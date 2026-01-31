# Foundations

Before we can do _anything_ we'll need to get some necessary tools, gather some required documentation, and find a way to debug the code we'll write.

- See [TOOLS.md](../guides/TOOLS.md) for a guide to tools you'll likely need to work in this repo. By design though, not much is needed!
- Check [RESOURCES.md](../guides/RESOURCES.md) for a list of references and other resources which might be handy.

Without the [tools](../guides/TOOLS.md) you won't be able to progress. This folder contains a test which confirms that all required tools are configured and working correctly.

## Aims

- Understand the very basics of compiling and inspecting RISC-V code.
- Run code on QEMU and use some basic commands to understand the state of the system.

## Testing our Toolchain

Dependencies for bedrock-bootstrap are minimal, but we do require a few tools. This folder contains a few minimal examples which will test that your local toolchain is set up correctly.

Put simply: if you can run `make all` in this directory without an error, you have everything you need to run through this whole repo!

First, if you haven't already you'll need to ensure that `RISCV_PREFIX` is set. If your RISC-V objdump command is called `riscv64-elf-objdump` then you don't need to do anything - that's the default (see [common.mk](../common.mk)).

If it's anything else (e.g. `riscv32-unknown-objdump`) you'll need to `export RISCV_PREFIX=riscv32-unknown-` (changed as appropriate).

### Running the Test

Simply run:

```bash
make all
```

This will delete all generated files from this directory and attempt to regenerate them. It tests both hex code (our ultimate aim!) and the assembler (which is helpful for illustrative purposes).

Several files will be dumped into this directory and intermediate files are placed in `BUILD/`.

For this example only, the artifacts are also committed into Git to make it possible to test QEMU without having the full toolchain set up and working.

The most interesting files are `toolchain-test-asm.elf` and `toolchain-test-hex.elf`. These are the kernels resulting from assembly and raw hex code respectively, and they can be run directly in QEMU.

There are also "dump" files showing the RISC-V instructions present in `toolchain-test.asm` and `toolchain-test.hex`.

There are also dumps of both the ASM and HEX examples, showing that they're equivalent.

```bash
$ cat toolchain-test-asm.dump toolchain-test-hex.dump

toolchain-test-asm.elf:     file format elf32-littleriscv


Disassembly of section .text:

80000000 <_start>:
80000000:	12345137          	lui	sp,0x12345

BUILD/toolchain-test-hex.bin:     file format binary


Disassembly of section .data:

00000000 <.data>:
   0:	12345137          	lui	sp,0x12345
```

Of course, the output isn't much use to us right now - we want to actually run it!

## Running in QEMU

We already know that the two kernels are functionally equivalent from the dumps, but still we provide two make targets for QEMU - one for ASM (`make qemu-asm`) and one for HEX (`make qemu-hex`).

Since they're functionally equivalent, you can choose either one to run. We'll show output from the hex version since that's more similar to what we'll use in the following steps of bootstrapping.

First, let's run QEMU:

```console
$ make qemu-hex
qemu-system-riscv32 -nographic -serial pty -gdb tcp::1234 -S -machine virt -bios none -kernel toolchain-test-hex.elf
QEMU 10.2.0 monitor - type 'help' for more information
char device redirected to /dev/ttys004 (label serial0)
(qemu)
```

Your output might differ slightly, but it should be close enough. First: to quit QEMU, type "quit" and then press return.

Now we'll attach GDB to ensure that works too.

In a separate terminal, run `make gdb`:

```console
$ make gdb
riscv64-elf-gdb -q -ex "set architecture riscv:rv32" -ex "target remote :1234"
The target architecture is set to "riscv:rv32".
Remote debugging using :1234
⚠️ warning: No executable has been specified and target does not support
determining executable automatically.  Try using the "file" command.
0x00001000 in ?? ()
(gdb)
```

Both the HEX and ASM programs are placed into memory at `0x80000000`. To confirm that they're loaded correctly, we'll dump the instruction at that location:

```text
(gdb) x/4i 0x80000000
   0x80000000:	lui	sp,0x12345
   0x80000004:	unimp
   0x80000006:	unimp
   0x80000008:	unimp
```

`x/4i` dumps 4 words of data at the given memory location, and disassembles them. We can see that the instruction at 0x8000000 matches the ASM and HEX code, and that the other instructions are empty (`unimp`).

If you're here without error, you've confirmed that your toolchain is fully working and you can move forwards!
