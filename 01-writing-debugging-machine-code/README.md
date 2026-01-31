# Writing and Debugging Raw Machine Code

If we want to write in pure machine code, we'll need to be able to write raw binary into a file.

## Aims

- Learn how to create raw machine code files
- Understand objdump

## Using `echo` to Write Machine Code

Let's focus for now on how to convert ASCII to binary. We can focus on the instruction `lui sp,0x12345` from [00-foundation](../00-foundation/README.md), which we saw in the dumps corresponds to a binary value of 0x12345137.

`echo` is a super simple, widely available program but actually works fine to illustrate what we're trying to achieve.

We can use `echo` to convert bytes into binary using the `\x00` escape sequence. `echo -n -e "\x13"` will echo the binary value 0x13 to stdout. We pass `-n` to prevent echo adding a newline (0xA) to the output.

To write our instruction, we have to remember to write bytes as little-endian; GDB shows us 32-bit instructions whereas we're writing raw bytes one at a time.

You can run and dump the instruction yourself, or use the Makefile.

```bash
$ make dump
# or ...
$ mkdir -p BUILD
$ echo -n -e "\x37\x51\x34\x12" > BUILD/example
$ $(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 BUILD/example
riscv64-elf-objdump -D -b binary -EL -m riscv:rv32 BUILD/example
BUILD/example:     file format binary


Disassembly of section .data:

00000000 <.data>:
   0:	12345137          	lui	sp,0x12345
```

We can see that we correctly wrote the `lui` instruction into a binary file!

## Understanding objdump

Let's look at the objdump parameters we used:

- `-D` (not `-d`) disassembles "all" in the file, meaning every instruction
- `-b binary` indicates we're dealing with a raw binary file (versus say, an ELF)
- `-m riscv:rv32` hints that we're dealing with RISC-V 32-bit instructions. This is required; there's no context in a raw binary file which might allow objdump to infer the architecture

Basically, we need to tell objdump that the file is just a raw list of RISC-V instructions. From that, it can disassemble the contents for us.

## Practical Machine Code

`echo` is a great way to introduce the concept, but obviously we _can_ use `echo` to write binary files... but having to prefix everything with `\x` would be a pain, so we need an easier way to write hex and convert to binary. If you've not already, you'll probably want to read [guides/XXD\_COMPILER.md](../guides/XXD_COMPILER.md) which explains how we'll write "hex" files in a bit more detail.

In short: we'll actually use `xxd` to convert hex to binary. You'll see this in various Makefiles across this repo (often paired with a `sed` oneliner to strip comments).
