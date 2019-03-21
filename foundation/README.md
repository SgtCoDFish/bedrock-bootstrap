# Foundations

Before we can do _anything_ we'll need to gather some required documentation, and find a way to debug the code we're going to write.

The [HiFive1 getting started guide](https://sifive.cdn.prismic.io/sifive%2F9c57065b-6d28-465b-b67d-f416894123a9_hifive1-getting-started-v1.0.2.pdf) is actually not the best place for us to get started. It uses a fairly heavyweight IDE and isn't targetted at true bare metal. There is however a [more detailed datasheet](https://sifive.cdn.prismic.io/sifive%2F4d063bf8-3ae6-4db6-9843-ee9076ebadf7_fe310-g000.pdf) which will be a useful reference for us.

[Freedom Metal](https://sifive.github.io/freedom-metal-docs/introduction.html#what-is-freedom-metal) (on [GitHub](https://github.com/sifive/freedom-metal/tree/master)) looks more useful, as does the [freedom-e-sdk](https://github.com/sifive/freedom-e-sdk). We see they provide
a [Board Support Package](https://github.com/sifive/freedom-e-sdk/tree/30c143eb5445f47edb351ba54c84ff8285dc27a9/bsp/sifive-hifive1) for the HiFive1, with details we need such as the architecture we're compiling for. We also see example programs
such as the vanishingly simple [hello](https://github.com/sifive/example-hello/tree/d1397bec64187efb8b791fe1eb307aa3c760c694) which we can use as a sanity check. The Makefile is a big beast designed to support multiple different boards in testing, and we don't need most of it.

We're not actually going to be _using_ Freedom Metal library directly, but since it's open source we can learn from it. Note also that it's in C, which isn't any use to us in a bedrock bare metal world.

While you're at it you'll probably want to clone the `freedom-e-sdk` just to have it around.

```bash
git clone --recursive https://github.com/riscv/riscv-gnu-toolchain
git clone --recursive https://github.com/sifive/freedom-e-sdk
```

Once that's done we'll need to build.

## Building Qemu

If your system has qemu >= 3.1, then in theory RISC-V support was upstreamed and you can probably install a RISC-V-supporting qemu from your system package manager. In practise, at least on macOS, there are problems connecting to an emulated kernel with gdb. Better to build it yourself; it's quick to do.

Otherwise you'll want to clone it yourself and build:

```bash
git clone --recursive https://github.com/sifive/riscv-qemu
cd riscv-qemu
mkdir build && cd build
../configure --target-list=riscv32-softmmu
make -j4
```

## The GNU Toolchain

You can use a prebuilt toolchain from [Sifive](https://www.sifive.com/boards/) (search for "GCC Toolchain") or build your own. You'll need one in either case; the compilers won't be much use, but some of the other tools will be.

If you're building from source, you can see in the freedom-e-sdk HiFive1 [BSP](https://github.com/sifive/freedom-e-sdk/blob/30c143eb5445f47edb351ba54c84ff8285dc27a9/bsp/sifive-hifive1/settings.mk) that we need to target a different arch and ABI, since the toolchain defaults to 64 bit.

You can choose a different prefix, but we'll assume you've set `$RISCV_PREFIX` to something. You'll need that evironment variable set to do anything.

Before building, you're likely to need to install some additional requirements. See [the repo](https://github.com/riscv/riscv-gnu-toolchain) for requirements on various platforms including Ubuntu and macOS.

```bash
cd riscv-gnu-toolchain
mkdir build && cd build
../configure --with-arch=rv32imac --with-abi=ilp32 --with-cmodel=medlow --prefix=$RISCV_PREFIX
make -j4  # might need to be done as root depending on where you're installing
```

The executables will be in the `$RISCV_PREFIX/bin` dir, and you can sanity check by running `$RISCV_PREFIX/bin/riscv32-unknown-elf-gcc -v`.

## Compiling Our Sample

Now we have a toolchain, we can test it using the files in this repo. Just run:

```bash
make RISCV_PREFIX=/path/to/riscv all
```

and several files will be dumped into BUILD in this directory. The most interesting ones from our perspective are `nothing.bin` and `nothing.dump`:

```bash
$ cat BUILD/nothing.elf
BUILD/nothing.elf:     file format elf32-littleriscv

Disassembly of section .text:

80000000 <main>:
80000000: 452d        li a0,11
80000002: 8082        ret

80000004 <_start>:
80000004: 80004137    lui sp,0x80004
80000008: ff9ff0ef    jal ra,80000000 <main>
8000000c: 9002        ebreak
8000000e: a001        j 8000000e <_start+0xa>

$ hexdump BUILD/nothing.bin
0000000 4501 8082 4137 8000 f0ef ff9f 9002 a001
0000010
```

Note that the raw binary dump (`nothing.bin`) is very small, and that the bytes match up with the disassembly in `nothing.dump`. For example, we see that the first instruction, at address `0x80000000` has the hex machine-code representation `452d` and that this matches the first 2 bytes of `nothing.bin`. The dump also has assembly on the right to make it easier to read, so we can see that the sole command is `li a0,11` which loads the value `11` into the register named `a0` (which is designated by the RISC-V ABI as being for return values).

Of course, the compiled output isn't much use to us right now. We want to run it on qemu to make sure it actually runs.

## Running in Qemu

First, note that when running Qemu headless, you exit by pressing Ctrl+A, releasing and then pressing `x` (use Ctrl+A and then `h` for help on other such commands).

We pass a few arguments to the following command which look initially confusing:

- `-s` starts a GDB debugger on port `1234` which lets us dump memory
- `-S` pauses the CPU before running anything, giving us time to debug
- `-machine sifive_e` tells qemu we're running a sifive-e machine (which the HiFive 1 is!)
- `-nographic` disables graphics - we're not going to need them.

Finally we point at our kernel (which is our ELF file) and kick off Qemu. macOS users take note: there are some problems with using the homebrew qemu with gdb. Best to compile your own qemu for this.

```bash
$ qemu-system-riscv32 -machine sifive_e -nographic -s -S -kernel BUILD/nothing.elf
# Open a separate terminal
$ $RISCV_PREFIX/bin/riscv32-unknown-elf-gdb
(gdb) target remote localhost:1234
(gdb) x 0x80000000
0x80000000:    0x8082452d
(gdb) info register pc # can also write "i r pc"
                       # or just "i r" to dump all registers
pc 0x1000     0x1000
(gdb) x 0x1000
0x1000:    0x204002b7
(gdb) x 0x1004
0x1004:    0x00028067
```

`x 0x80000000` dumps the memory at that address, which we can see is our `main` function above. We know we've loaded the file correctly into qemu.

So we see our program is loaded correctly, and `info register pc` shows us that the board has defaulted `pc` to `0x1000`, and at that address there's actually an instruction already there which we didn't write! We can cheat and look up in the HiFive1 manual for clues as to what that instruction does, but this is actually a great opportunity to write our first machine code!

## Writing Raw Machine Code

If we want to write in pure machine code, we'll need to be able to write raw bytes into a file. We'll come onto tooling for that later, but for now we can make do with a very quick solution.

In `bootloader` there's a Python script which will dump raw binary into a file called "bootloader" in the same directory. The order looks "reversed" compared to the instructions we see above to account for endianness when we're writing to the file - gdb is showing us full hex values, which we need to write little-endian.

```bash
$ python3 write_bootloader.py && hexdump bootloader
0000000 b7 02 40 20 67 80 02 00
0000008
```

Now we've written the instructions into a file, we can use `riscv32-unknown-elf-objdump` to help us work out what they are, as long as we give the disassembler a few tips about what exactly it's disassembling:

- `-D` (not `-d`!) dissassembles "all" in the file, so every instruction
- `-b binary` indicates we're dealing with a raw binary file (as oppossed to, say, an ELF)
- `-m riscv:rv32` hints that we're dealing with RISC-V 32-bit instructions, since there's no context in our binary file to indicate that this is what it contains

```bash
$ $RISCV_PREFIX/bin/riscv32-unknown-elf-objdump -D -b binary -m riscv:rv32 bootloader

bootloader: file format binary

Disassembly of section .data:

00000000 <.data>:
   0:    b7024020    lui t0,0x20400
   4:    67800200    jr  t0 # 0x20400000
```

And the "bootloader" is revealed!

`lui t0,0x20400` basically loads the unsigned value `0x20400000` into `t0` which is a "temporary" register (also known as `x5`). You'll note that the immediate value in the instruction, `0x20400` doesn't match the value that ends up in `t0`. To clear that up, the [manual](https://content.riscv.org/wp-content/uploads/2016/06/riscv-spec-v2.1.pdf) says:

> LUI (load upper immediate) is used to build 32-bit constants ... LUI places the immediate value in the top 20 bits of the destination register, filling in the lowest 12 bits with zeros.

Since one hex "value" is 4 bits, that means the bottom 12 bits of the loaded value are `000`, which explains the difference!

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

So we know that after booting, control will pass to 0x20400000 immediately and then everything hangs. What's at that location?

```bash
(gdb) x 0x20400000
0x20400000:    0x00000000
```

The answer: absolutely nothing! `0x00000000` is an illegal instruction, which causes the process to trap and thereby sets the PC to `0x00000000`... which in RISC-V always contains `0x0` by definition and so causes an infinite loop!

## Booting on the HiFive1

If you "cheated" like was mentioned earlier and read the HiFive1 documentation regarding the boot process, you'll have noticed that the process for the HiFive1 is different.

The HiFive1 comes with slightly more code with the intention of making it easier to develop for, taking into account the fact that it's harder to develop on hardware than it is under an emulator.

We can get the more detail in the following lightly edited description of the boot process from [a datasheet](https://sifive.cdn.prismic.io/sifive%2Ffeb6f967-ff96-418f-9af4-a7f3b7fd1dfc_fe310-g000-ds.pdf), where the sections in brackets are added.

(Note also that that datasheet is actually old, but the description is in some ways easier to follow. More "up to date" details are available [here](https://sifive.cdn.prismic.io/sifive%2F4d063bf8-3ae6-4db6-9843-ee9076ebadf7_fe310-g000.pdf).).

> The FE310-G000 \[starts at address `0x0001_0000`\] and boots by jumping to the beginning of OTP memory \[at `0x0002_0000`\] and executing code found there. As shipped, OTP memory at the boot location is preprogrammed to jump immediately to the end of the OTP memory \[around `0x0002_1FFF`\], which contains the following code to jump to the beginning of the SPI-Flash at `0x2000_0000`:

```assembly
fence 0,0
li t0, 0x20000000
jr t0
```

> `fence 0,0` is encoded as `0x0000000F`, and the instruction may be modified by burning additional bits to transform it into a `JAL` instruction (opcode `0x6F`) to execute arbitrary code rather than jumping directly to the beginning of the SPI-Flash.

(Note that "OTP" means "one time programmable" memory - that is, once you "burn" a program there, it's there permanently. [SPI Flash](https://en.wikipedia.org/wiki/Serial_Peripheral_Interface) means flash memory connected over SPI)

What matters most from the above text is the code; `fence 0,0` is best described in [this StackOverflow answer](https://stackoverflow.com/a/26374650) and is effectively a no-op here to enable the neat trick in the second paragraph. We're not going to be burning OTP any time soon, so we'll ignore it.

`li t0,0x20000000` and `jr t0` look very similar to the qemu code we saw above and there aren't any surprises: we load `0x20000000` into `t0` and then jump there.

What happens at in SPI-Flash at `0x20000000`? As shipped, there's [a program](https://github.com/sifive/freedom-e-sdk/tree/f9271b91257e0a8a989faf3eff0757ee46694fe0/software/double_tap_dontboot) written there, whose source is reproduced in this directoy in `double_tap_dontboot.c`. That doesn't mean much to us since we'll be overwriting it, but it's neat (and very cool of SiFive!) to have the insight and be able to restore the program if we choose.

## NEXT

// TODO: Explain how we'll handle the difference between qemu and hardware and actually run some code from boot

## Other Links

- dwelch67 always has good guides and you can check his [uart01 sample](https://github.com/dwelch67/sifive_samples/tree/master/hifive1/uart01) for a baremetal UART example, which is getting a little ahead of ourselves if we're focusing on bedrock.
- Running RISC-V on qemu bare metal [google groups thread](https://groups.google.com/a/groups.riscv.org/forum/#!topic/sw-dev/IET9LBFJohU)
- A very minimal bare metal example, similar to dwelch but based on riscv64 and SPIKE: [schoeberl](https://github.com/schoeberl/cae-examples)
