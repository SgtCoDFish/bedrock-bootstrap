# Foundations

Before we can do _anything_ we'll need to get some necessary tools, gather some required documentation, and find a way to debug the code we'll write.

- See [TOOLS.md](../guides/TOOLS.md) for a guide to tools
- Check [RESOURCES.md](../guides/RESOURCES.md) for a list of references and other resources which might be handy.
- Read [OPENOCD\_WRITING.md](../guides/OPENOCD_WRITING.md) for a guide to uploading code onto a HiFive1 using OpenOCD

Without the [tools](../guides/TOOLS.md) you'll struggle to do much at all.

## Aims

- Understand the very basics of compiling and inspecting RISC-V code.
- Run a binary on QEMU and use some basic commands to inspect the state of the system.

## Compiling Our Toolchain Test

**NOTE**: The linker-script used for these examples is intentionally "weird". Running `nothing.elf` on hardware won't work. The idea here is to use QEMU to check that our toolchain is doing what we expect; actually running stuff on hardware comes later.

We can test our toolchain using the files in this repo. As long as you've set `$RISCV_PREFIX` correctly and you've got a riscv compiler installed, you should be good to go.

Building the binaries is simple; you can check the underlying `Makefile` for details of the specific commands.

```bash
make all
```

Several files will be dumped into this directory and intermediate files are placed in `BUILD/`.

For this example only, the artifacts are also committed into Git to make it possible to test QEMU without having the full toolchain set up and working.

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
# dumps nothing.bin in ASCII hex, one byte at a time
000000 13 05 b0 00 67 80 00 00 37 41 00 80 ef f0 5f ff
000010 73 00 10 00 6f 00 00 00
000018

$ od -Ax -tx4 nothing.bin
# dumps nothing.bin in ASCII hex, four bytes at a time
000000 00b00513 00008067 80004137 ff5ff0ef
000010 00100073 0000006f
000018
```

Note that the raw binary dump (`nothing.bin`) is very small, and that the bytes match up with the disassembly in `nothing.dump`. For example, we see that the first instruction, at address `0x80000000` has the hex machine-code representation `00b00513` and that this matches the first 4 bytes of `nothing.bin` in the hex dumps using `od`.

The dump also has assembly on the right to make it easier to read, so we can see that the sole command in `main` is `li a0,11` which loads the value `11` into the register named `a0` (which is designated by the RISC-V [Calling Convention](https://riscv.org/wp-content/uploads/2015/01/riscv-calling.pdf) as being for return values).

Of course, the compiled output isn't much use to us right now - we want to actually run it!

## Running in QEMU

First, we should note how to exit QEMU when running headless, which is done by pressing Ctrl+A, releasing and then pressing `x` (use Ctrl+A and then `h` for help on other such commands).

We can run qemu using `make qemu` but first we'll take a look at some of the arguments we pass to QEMU:

- `-gdb tcp::1234` starts a GDB debugger on port `1234` which lets us dump memory
- `-S` pauses the CPU before running anything, giving us time to attach using GDB
- `-machine sifive_e` tells QEMU we're running a `sifive-e` machine (which the HiFive1 is!)
- `-nographic` disables graphics - we're not going to need them.

Finally we point at our "kernel" - which is our ELF file - and kick off QEMU. If you encounter any issues using QEMU from a package manager, it's very easy to build your own.

```bash
$ make riscv
# Open a separate terminal
$ "$RISCV_PREFIX"gdb
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

`x 0x80000000` dumps the memory at that address, which we can see is our `main` function above. We know we've loaded the file correctly into QEMU.

So we see our program is loaded correctly, and `info register pc` shows us that the board has defaulted `pc` to `0x1000`, and at that address there's actually an instruction already there which we didn't write!

We _could_ cheat and look up in the HiFive1 manual for clues as to what that instruction does, but this is actually a great opportunity to write some raw machine code and disassemble it, which we'll try in the next section.
