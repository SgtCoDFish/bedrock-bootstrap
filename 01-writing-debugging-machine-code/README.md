# Writing and Debugging Raw Machine Code

If we want to write in pure machine code, we'll need to be able to write raw bytes into a file. We'll use the bootloader instructions we discovered at `0x1000` and `0x1004` in the previous section, write them to a binary file and then inspect that file to learn what the instructions are.

## Aims

- Learn how to create raw machine code files
- Dump and inspect raw machine code

## Using `echo` to Write Machine Code

`echo` is a super simple, widely available program but actually works fine to get the idea of what we're trying to achieve here.

We must remember to write the bytes as little-endian; remember that gdb shows us 32-bit instructions in hex, whereas we're writing raw bytes. We also need to pass `-n` so that echo doesn't append a newline (which would show up as `0x0a`). You can run the commands yourself or run `make BUILD/bootloader1`:

```bash
$ mkdir -p BUILD
$ echo -n -e "\xb7\x02\x40\x20\x67\x80\x02\x00" > BUILD/bootloader1
$ od -Ax -tx1 BUILD/bootloader1
000000 b7 02 40 20 67 80 02 00
000008
```

Now we've written the instructions into a file, we can use `riscv32-unknown-elf-objdump` to help us work out what they are, as long as we give the disassembler a few tips:

- `-D` (not `-d`) disassembles "all" in the file, meaning every instruction
- `-b binary` indicates we're dealing with a raw binary file (versus say, an ELF)
- `-m riscv:rv32` hints that we're dealing with RISC-V 32-bit instructions, since there's no context in a raw binary file which might allow objdump to infer the architecture

```bash
$ $RISCV_PREFIX/bin/riscv32-unknown-elf-objdump -D -b binary -m riscv:rv32 BUILD/bootloader1
Disassembly of section .data:

00000000 <.data>:
   0:    b7024020    lui t0,0x20400
   4:    67800200    jr  t0 # 0x20400000
```

And the first bootloader is revealed!

## Understanding the First Bootloader

`lui t0,0x20400` loads the unsigned value `0x20400000` into `t0` which is a "temporary" register. You'll note that the immediate value in the instruction, `0x20400` doesn't _exactly_ match the value that ends up in `t0`. To clear that up, the [manual](https://content.riscv.org/wp-content/uploads/2016/06/riscv-spec-v2.1.pdf) says:

> LUI (load upper immediate) is used to build 32-bit constants ... LUI places the immediate value in the top 20 bits of the destination register, filling in the lowest 12 bits with zeros.

Since one hex "value" is 4 bits, that means the bottom 12 bits of the loaded value are `000`, which explains the difference.

`jr t0` unconditionally jumps to the address in `t0`, which has the effect of setting `pc` to the value in `t0`.

So we conclude that the bootloader just immediately jumps to the address `0x2040_0000`! Let's try it out in a GDB session connected to QEMU with `nothing.elf` from the previous section:

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

That pretty much confirms what we expected to happen. We set `t0` and then jump to the address it contains, where the debugger hangs after jumping to `0x0`. Even after the hang, we still see the value that we set in `t0`. This was to be expected; we didn't write any code at `0x2040_0000` and there's nothing placed there at design time.

In the next section, we'll investigate the boot process further, on both QEMU and HiFive1 hardware.
