# Understanding Kernel Formats

Now we understand the boot process and how control passes to our programs, we need to get some code into the correct location to run it, which will require some more research and analysis. We'll learn about the ELF format. If you know about the ELF format's inner structure already, this section won't provide much of use!

## Aims

- Understand what format we need to use to get code running on both platforms
- Know which parts of an ELF file are important for our purposes and which are not needed.

## ret1234

In this directory we have a simple assembly language program which loads the value `0x01234` into the upper part of `x1`. The code is placed (via the linker script, `linker.ld`) into memory address `0x2040_0000` so it can be booted on both QEMU and the HiFive1.

After loading the value, it calls `ebreak` which breaks to a debugger, and then loops in the same position infinitely.

The binaries are once again committed to git in the form of `ret1234.elf` - an ELF file - and `ret1234.bin` which is a raw binary file with no relocation information.

## File Formats

So what binary file format - ELF or raw binary - do we actually need? We'll need to analyse both QEMU and our HiFive1 hardware.

### QEMU

So far we've used an ELF file to provide the kernel for QEMU, and in fact that's the only choice we have as evidenced by the riscv-qemu source code for [load_kernel](https://github.com/riscv/riscv-qemu/blob/32a1a94dd324d33578dca1dc96d7896a0244d768/hw/riscv/sifive_e.c#L77-L88)[1].

That means that to run under QEMU, we always have to provide an [ELF](https://en.wikipedia.org/wiki/Executable_and_Linkable_Format) file[2].

That introduces a decent amount of complexity in our build process that it would be desirable to avoid since it'll be harder to generate an ELF than a raw binary file - but there's no reason it should stop us.

### HiFive1

When developing for the HiFive1 we can choose to produce ELF files (which can be uploaded natively by OpenOCD) but we also have the option to use OpenOCD's `write_image` directive to write a raw binary file at a memory offset we specify. This is friendlier to raw machine code since we could avoid the cost of having to create an ELF header, but QEMU forces our hand

## ELF Files

In this directory we see `ret1234.elf`. From [Wikipedia](https://en.wikipedia.org/wiki/Executable_and_Linkable_Format#File_header) we learn that a 32-bit ELF binary has a 52-byte long header, followed by program headers and section headers. The Wikipedia article is an excellent reference for the format and you'll probably want to have the article open while analysing the file.

The following is a brief overview of the different sections. More in-depth analysis is provided in [ELF\_DETAIL](./ELF_DETAIL.md)

### ELF Header

The first 52 bytes of ret1234.elf, when dumped, show the ELF file header. This mostly doesn't change between different RV32I ELF files in our case, since we'll be aiming for a very simple file with very few program headers and section headers.

### Program Header

```text
$ od -Ax -tx1 -N32 -j52 ret1234.elf
000034 01 00 00 00 00 10 00 00 00 00 40 20 00 00 40 20
000044 0c 00 00 00 0c 00 00 00 05 00 00 00 00 10 00 00
```

In short, this header points at a `0xc` byte long executable section which is located at offset `0x1000` in this file, and indicates that the section should be loaded at `0x2040_0000`.

We can dump that section and confirm that the expected code is present:

```text
$ od -v -Ax -tx1 -j4096 -N12 ret1234.elf
001000 b7 40 23 01 73 00 10 00 6f 00 00 00
```

The program header is one of the parts of the file that must change for each program, since it contains the size of the program itself.

### Section Headers

The 6 section headers in this file are located at `0x000010ac` and, according to the ELF file header, each has a size of 40 bytes.

The types are summarised here, with more detail available [in this doc](https://docs.oracle.com/cd/E23824_01/html/819-0690/chapter7-1.html#scrolltoc):

- a null, empty section, which appears to be required
- `.text` is the most important section type since it contains our code! The others are less important, but some are required.
- `.riscv.attributes` is a special RISC-V section which contains details about the specific RISC-V architecture the object was compiled for - `rv32i` in this case - among other things. We can ignore it for now, and it's not well documented in any case.
- `.symtab` is the [symbol table](https://docs.oracle.com/cd/E23824_01/html/819-0690/chapter6-79797.html#scrolltoc), which is needed to perform relocations on symbolic definitions and references.
- `.strtab` is a table containing null-terminated strings relating to symbol table entries.
- `.shstrtab` is a table also containting null-terminated strings which are the names of section headers

Three of these sections are vital (aside from NULL which is required):

- `.text` contains our code
- `.shstrtab` which has header names for everything. We can get rid of the separate `.strtab` and place everything in here if needed.
- `.symtab` which points to our code and gives details about how function calls work. This will be very simple for our use cases.

### Section Bodies

`.shstrtab` and `.strtab` both contain lists of null-terminated strings, and `.text` just contains our compiled code.

`.symtab` has a "global" section which points at our code, but the rest of the sections can be ignored.

The program section pointed to by the program header is the same as the `.text` section in our simple case.

Mostly, there's a large amount of padding up to `0x1000` bytes which is all zeroes. We can get rid of almost all of it.

## Next

We've seen what an ELF file looks like, and we know that we need an ELF file for QEMU. The next step is to write an ELF header in hex.

## Notes

[1] We see a call to `load_elf` which ultimately ends up [load\_elf\_ram\_sym](https://github.com/riscv/riscv-qemu/blob/32a1a94dd324d33578dca1dc96d7896a0244d768/hw/core/loader.c#L461) - the ELF parsing logic might be useful to us later so it's noted here, but in any case it's neat to drill down into these things.

[2] The alternative to generating ELF files for both platforms is to patch QEMU's HiFive-compatible device to support raw binary files, perhaps with a specified offset.

We could for example patch the `load_kernel` function to parse `-kernel "0x20400000:ret1234.bin"` as "place ret1234.bin at `0x2040_0000`". It's a judgement call, but it feels like patching QEMU is further away from what we're trying to do than just banging out an ELF header, which is what we'll be doing later.

At the end of the day, one of the tasks involves writing C and the other involves writing hex!
