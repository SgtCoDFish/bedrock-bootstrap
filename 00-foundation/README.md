# Foundations

Before we can do _anything_ we'll need to get some necessary tools, gather some required documentation, and find a way to debug the code we'll write.

- See [TOOLS.md](../guides/TOOLS.md) for a guide to tools
- Check [RESOURCES.md](../guides/RESOURCES.md) for a list of references and other resources which might be handy.
- Read [OPENOCD_WRITING.md](../guides/OPENOCD_WRITING.md) for a guide to uploading code onto a HiFive1 using OpenOCD

Without the [tools](../guides/TOOLS.md) you'll struggle to do much at all.

## Compiling Our Toolchain Test

We can test our toolchain using the files in this repo. Be warned however that the code built for these examples is intentionally placed into a "weird" location and that running `nothing.elf` on a HiFive1 might give strange results. The idea here is to use Qemu to check that our toolchain is doing what we expect.

The command to run for (re-)building the binaries is simple:

```bash
make RISCV_PREFIX=/path/to/riscv all
```

Several files will be dumped into this directory and intermediate files are placed in `BUILD/`. The end-product files are also committed into Git to make it possible to test the tools without having the full GNU toolchain.

The most interesting files from our perspective are `nothing.bin` and `nothing.dump`:

```bash
$ cat nothing.dump
nothing.elf:     file format elf32-littleriscv


Disassembly of section .text:

80000000 <main>:
80000000:   00b00513           li a0,11
80000004:   00008067           ret

80000008 <_start>:
80000008:   80004137           lui sp,0x80004
8000000c:   ff5ff0ef           jal ra,80000000 <main>
80000010:   00100073           ebreak
80000014:   0000006f           j 80000014 <_start+0xc>

$ od -Ax -tx1 nothing.bin
000000 13 05 b0 00 67 80 00 00 37 41 00 80 ef f0 5f ff
000010 73 00 10 00 6f 00 00 00
000018
```

Note that the raw binary dump (`nothing.bin`) is very small, and that the bytes match up with the disassembly in `nothing.dump`. For example, we see that the first instruction, at address `0x80000000` has the hex machine-code representation `00b00513` and that this matches the first 4 bytes of `nothing.bin` if you account for the disassembly showing big-endian and the hexdump showing raw bytes.

The dump also has assembly on the right to make it easier to read, so we can see that the sole command in `main` is `li a0,11` which loads the value `11` into the register named `a0` (which is designated by the RISC-V [Calling Convention](https://riscv.org/wp-content/uploads/2015/01/riscv-calling.pdf) as being for return values).

Of course, the compiled output isn't much use to us right now. We want to run it -on qemu or on hardware - to make sure it actually runs.

## Running in Qemu

First, note that when running Qemu headless, you exit by pressing Ctrl+A, releasing and then pressing `x` (use Ctrl+A and then `h` for help on other such commands).

We pass a few arguments to the following command which look initially confusing:

- `-s` starts a GDB debugger on port `1234` which lets us dump memory
- `-S` pauses the CPU before running anything, giving us time to debug
- `-machine sifive_e` tells qemu we're running a sifive-e machine (which the HiFive 1 is!)
- `-nographic` disables graphics - we're not going to need them.

Finally we point at our "kernel" - which is our ELF file - and kick off Qemu. macOS users take note: you may run into problems with using the homebrew qemu with gdb. Best to compile your own qemu for this.

```bash
$ qemu-system-riscv32 -machine sifive_e -nographic -s -S -kernel nothing.elf
# Open a separate terminal
$ $RISCV_PREFIX/bin/riscv32-unknown-elf-gdb
(gdb) target remote localhost:1234
(gdb) x 0x80000000
0x80000000:    0x00b00513
(gdb) info register pc # can also write "i r pc"
                       # or just "i r" to dump all registers
pc 0x1000     0x1000
(gdb) x 0x1000
0x1000:    0x204002b7
(gdb) x 0x1004
0x1004:    0x00028067
```

`x 0x80000000` dumps the memory at that address, which we can see is our `main` function above. We know we've loaded the file correctly into qemu.

So we see our program is loaded correctly, and `info register pc` shows us that the board has defaulted `pc` to `0x1000`, and at that address there's actually an instruction already there which we didn't write! We _could_ cheat and look up in the HiFive manual for clues as to what that instruction does, but this is actually a great opportunity to write some raw machine code and disassemble it.

## Writing Raw Machine Code

If we want to write in pure machine code, we'll need to be able to write raw bytes into a file. We might need better tooling in the future, but for now we can just use `echo`.

We must remember to write the bytes as little-endian; remember that gdb shows us 32-bit instructions in hex, whereas we're writing raw bytes. We also need to pass `-n` so that echo doesn't append a newline (which would show up as `0x0a`). You can run the commands yourself or run `make BUILD/bootloader1`:

```bash
$ mkdir -p BUILD
$ echo -n -e "\xb7\x02\x40\x20\x67\x80\x02\x00" > BUILD/bootloader1
$ od -Ax -tx1 BUILD/bootloader1
000000 b7 02 40 20 67 80 02 00
000008
```

Now we've written the instructions into a file, we can use `riscv32-unknown-elf-objdump` to help us work out what they are, as long as we give the disassembler a few tips:

- `-D` (not `-d` which looks similar) dissassembles "all" in the file, meaning every instruction
- `-b binary` indicates we're dealing with a raw binary file (as oppossed to, say, an ELF)
- `-m riscv:rv32` hints that we're dealing with RISC-V 32-bit instructions, since there's no context in a binary file which might allow objdump to work out the architecture

```bash
$ $RISCV_PREFIX/bin/riscv32-unknown-elf-objdump -D -b binary -m riscv:rv32 BUILD/bootloader1
Disassembly of section .data:

00000000 <.data>:
   0:    b7024020    lui t0,0x20400
   4:    67800200    jr  t0 # 0x20400000
```

And the first bootloader is revealed!

`lui t0,0x20400` loads the unsigned value `0x20400000` into `t0` which is a "temporary" register. You'll note that the immediate value in the instruction, `0x20400` doesn't _exactly_ match the value that ends up in `t0`. To clear that up, the [manual](https://content.riscv.org/wp-content/uploads/2016/06/riscv-spec-v2.1.pdf) says:

> LUI (load upper immediate) is used to build 32-bit constants ... LUI places the immediate value in the top 20 bits of the destination register, filling in the lowest 12 bits with zeros.

Since one hex "value" is 4 bits, that means the bottom 12 bits of the loaded value are `000`, which explains the difference.

`jr t0` unconditionally jumps to the address in `t0`, which has the effect of setting `pc` to the value in `t0`.

So we conclude that the bootloader just immediately jumps to the address 0x20400000! Let's try it out:

```bash
(gdb) nexti
0x00001004 in ?? ()
(gdb) i r pc t0
pc    0x1004
t0    0x20400000
(gdb) nexti
<execution hangs, so we press Ctrl-C>
^C
Program received signal SIGINT, Interrupt.
0x00000000 in ?? ()
(gdb) i r pc t0
pc    0x0
t0    0x20400000
```

That pretty much confirms what we expected to happen. We set `t0` and then jump to the address it contains, where the debugger mysteriously hangs. Even after the hang, we still see the value that we set in `t0`.

## Boot Part 2

So we know that after the first bootloader, control will pass to 0x20400000 immediately and then everything hangs. What's at that location?

```bash
(gdb) x 0x20400000
0x20400000:    0x00000000
```

The answer: absolutely nothing! `0x00000000` is an illegal instruction, which causes the process to trap and thereby sets the PC to `0x00000000`... which in RISC-V always contains `0x0` by definition and so causes an infinite loop! We'll get our code running later, once we've figured out how to get to this point on actual hardware.

## Booting on the HiFive1

If you already read the HiFive1 documentation regarding the boot process, you'll have noticed that the process for the HiFive1 is slightly different.

The HiFive1 comes with slightly more code with the intention of making it easier to develop for, taking into account the fact that it's harder to develop on hardware than it is under an emulator.

We can get the more detail in the following lightly edited description of the boot process from [a datasheet](https://sifive.cdn.prismic.io/sifive%2Ffeb6f967-ff96-418f-9af4-a7f3b7fd1dfc_fe310-g000-ds.pdf)[1], where the sections in brackets are added.

> The FE310-G000 \[starts at address `0x0001_0000`\] and boots by jumping to the beginning of OTP memory \[at `0x0002_0000`\] and executing code found there. As shipped, OTP memory at the boot location is preprogrammed to jump immediately to the end of the OTP memory \[`0x0002_1FF4`\], which contains the following code to jump to the beginning of the SPI-Flash at `0x2000_0000`:

```assembly
0x0002_1FF4:
   0:    0000000f    fence 0,0
   4:    200002b7    li t0, 0x20000000
   8:    00028067    jr t0
```

(Note that "OTP" means "one time programmable" memory - that is, once you "burn" a program there, it's there permanently. "SPI Flash" means flash memory connected over [SPI](https://en.wikipedia.org/wiki/Serial_Peripheral_Interface))

What matters most from the above text is the code; `fence 0,0` is described in [this StackOverflow answer](https://stackoverflow.com/a/26374650) and is provided essentially to enable a neat trick[2]. We can ignore it as a no-op.

`li t0,0x20000000` and `jr t0` look very similar to the qemu code we saw above and there aren't any surprises: we load `0x20000000` into `t0` and then jump there.

What happens at in SPI-Flash at `0x20000000`? As shipped, there's [a program](https://github.com/sifive/freedom-e-sdk/tree/f9271b91257e0a8a989faf3eff0757ee46694fe0/software/double_tap_dontboot) written there, whose source is reproduced in the `scratch` directory of this project along with a dump of the bootloader from real hardware. We might choose to overwrite that bootloader code later, but for now we can basically just take it as a given that it jumps to `0x2040_0000`.

To summarise, the HiFive1 boots like this:

- Start at `0x0001_0000` (this is the same as Qemu)
- Jump to  `0x0002_0000` (start of OTP)
- Jump to  `0x0002_1FF4` (end of OTP)
- Jump to  `0x2000_0000` (start of SPI-Flash)
- As shipped: Do some init and magic stuff, then jump to `0x2040_0000`

That means our code needs to live at `0x2040_0000`. In the next section, we'll get something of our own there.

## Notes

[1] (Note also that that datasheet is actually old, but the description of the boot process is in some ways easier to follow. More "up to date" details are available [here](https://sifive.cdn.prismic.io/sifive%2F4d063bf8-3ae6-4db6-9843-ee9076ebadf7_fe310-g000.pdf)).

[2] From the same datasheet's boot description:
> `fence 0,0` is encoded as `0x0000000F`, and the instruction may be modified by burning additional bits to transform it into a `JAL` instruction (opcode `0x6F`) to execute arbitrary code rather than jumping directly to the beginning of the SPI-Flash.

This means that it's easy to change `0x0000_000F` into a different instruction by burning bits. We won't be doing anything so permanent any time soon!