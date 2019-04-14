# Running Machine Code on Qemu and HiFive 1

Now we understand the boot process and how control passes to our programs, we need to get some code into the correct location to run it, which will require some more understanding.

## Aims

- Understand what format we need to use to get code running on both platforms
- Be able to create machine-code files which can be run on Qemu or the HiFive1

## ret1234

In this directory we have a simple assembly language program which loads the value `0x01234` into the upper part of `x1`. The code is placed (via the linker script, `linker.ld`) into memory address `0x2040_0000` so it can be booted on both Qemu and the HiFive1.

After loading the value, it calls `ebreak` which breaks to a debugger, and then loops in the same position infinitely.

The binaries are once again committed to git in the form of `ret1234.elf` - an ELF file - and `ret1234.bin` which is a raw binary file.

## File Formats

### Qemu

So far we've used an ELF file to provide the kernel for Qemu, and in fact that's the only choice we have as evidenced by the riscv-qemu source code for [load_kernel](https://github.com/riscv/riscv-qemu/blob/32a1a94dd324d33578dca1dc96d7896a0244d768/hw/riscv/sifive_e.c#L77-L88)[1].

That means that to run under Qemu, we always have to provide an [ELF](https://en.wikipedia.org/wiki/Executable_and_Linkable_Format) file[2]. That introduces a decent amount of complexity in our build process that it would be desirable to avoid since it'll be harder to generate an ELF than a raw binary file if we're writing in machine code - but there's no reason it should stop us. It's worth it to be able to debug on Qemu.

### HiFive1

When developing for the HiFive1 we could choose to produce ELF files (which can be uploaded natively by OpenOCD). We also have the option to use OpenOCD's `write_image` directive to write a raw binary file at a memory offset we specify. This is friendlier to raw machine code since we could avoid the price of having to create an ELF header, but Qemu forces our hand.

## ELF Files

In this directory we see `ret1234.elf`. From [Wikipedia](https://en.wikipedia.org/wiki/Executable_and_Linkable_Format#File_header) we learn that a 32-bit ELF binary has a 52-byte long header, followed by program headers and section headers. Let's examine the contents.

### ELF Header

We'll dump the first 52 bytes of ret1234.elf, which is the ELF header that can be compared against the detail given in the Wikipedia article on the ELF format. This header is further analysed in [3] and a full dump of the file using `readelf` is given in [7]. Despite all that, we can ignore most of the details.

```text
od -Ax -tx1 -N52 ret1234.elf
000000 7f 45 4c 46 01 01 01 00 00 00 00 00 00 00 00 00
000010 02 00 f3 00 01 00 00 00 00 00 40 20 34 00 00 00
000020 ac 10 00 00 00 00 00 00 34 00 20 00 01 00 28 00
000030 06 00 05 00
```

### Program Header

```text
$ od -Ax -tx1 -N32 -j52 ret1234.elf
000034 01 00 00 00 00 10 00 00 00 00 40 20 00 00 40 20
000044 0c 00 00 00 0c 00 00 00 05 00 00 00 00 10 00 00
```

Cross referencing with Wikipedia again, there's further analysis in [4]. In short, this header points at a `0xc` byte long executable section which is located at offset `0x1000` in this file, and indicates that the section should be loaded at `0x2040_0000`.

We can dump that section and confirm that the expected code is present:

```text
$ od -v -Ax -tx1 -j4096 -N12 ret1234.elf
001000 b7 40 23 01 73 00 10 00 6f 00 00 00
```

### Section Headers

The 6 section headers in this file are located at `0x000010ac` and have a size of 40 bytes, according to the ELF header. We can dump them all, which is shown (annotated) here:

```text
# 0: Pointless NULL header inserted by the compiler
$ od -v -Ax -tx1 -j4268 -N40 ret1234.elf
0010ac 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
0010bc 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
0010cc 00 00 00 00 00 00 00 00

# 1: .text section header
$ od -v -Ax -tx1 -j4308 -N40 ret1234.elf
0010d4 1b 00 00 00 01 00 00 00 06 00 00 00 00 00 40 20
0010e4 00 10 00 00 0c 00 00 00 00 00 00 00 00 00 00 00
0010f4 04 00 00 00 00 00 00 00

# 2: .riscv.attributes section header
$ od -v -Ax -tx1 -j4348 -N40 ret1234.elf
0010fc 21 00 00 00 03 00 00 70 00 00 00 00 00 00 00 00
00110c 0c 10 00 00 24 00 00 00 00 00 00 00 00 00 00 00
00111c 01 00 00 00 00 00 00 00

# 3: .symtab section header
$ od -v -Ax -tx1 -j4388 -N40 ret1234.elf
001124 01 00 00 00 02 00 00 00 00 00 00 00 00 00 00 00
001134 30 10 00 00 40 00 00 00 04 00 00 00 03 00 00 00
001144 04 00 00 00 10 00 00 00

# 4: .strtab section header
$ od -v -Ax -tx1 -j4428 -N40 ret1234.elf
00114c 09 00 00 00 03 00 00 00 00 00 00 00 00 00 00 00
00115c 70 10 00 00 08 00 00 00 00 00 00 00 00 00 00 00
00116c 01 00 00 00 00 00 00 00

# 5: .shstrtab section header
$ od -v -Ax -tx1 -j4468 -N40 ret1234.elf
001174 11 00 00 00 03 00 00 00 00 00 00 00 00 00 00 00
001184 78 10 00 00 33 00 00 00 00 00 00 00 00 00 00 00
001194 01 00 00 00 00 00 00 00
```

The types are summarised here, with more detail available [in this doc](https://docs.oracle.com/cd/E23824_01/html/819-0690/chapter7-1.html#scrolltoc):

- `.text` is the most important section type since it contains our code! The others are less important, but some are required.
- `.riscv.attributes` is a special RISC-V section which contains details about the specific RISC-V architecture the object was compiled for - `rv32ima` in our case - among other things. We can ignore it for now.
- `.symtab` is the [symbol table](https://docs.oracle.com/cd/E23824_01/html/819-0690/chapter6-79797.html#scrolltoc), which is needed to perform relocations on symbolic definitions and references.
- `.strtab` is a table containing null-terminated strings relating to symbol table entries.
- `.shstrtab` is a table also containting null-terminated strings which are the names of section headers

#### `.text` Section

```text
$ od -v -Ax -tx1 -j4308 -N40 ret1234.elf
0010d4 1b 00 00 00 01 00 00 00 06 00 00 00 00 00 40 20
0010e4 00 10 00 00 0c 00 00 00 00 00 00 00 00 00 00 00
0010f4 04 00 00 00 00 00 00 00
```

The first 4 bytes are the index into the `.shstrtab` section which has the header's name - `.text`. The second 4 bytes give the section type, which is `0x1` in this case, denoting program data.

Next `0x06` gives the flags for the section; this is equal to `0x2 | 0x4`, which means that the section "occupies memory during execution" and is "executable", respectively. `0x2040_0000` is the now-familiar address at which the section is placed.

The section's offset in the file is given next: `0x0000_1000` (4096), followed by the size of the section (`0xc` or 12 bytes). The section has already been dumped above, following the program header section.

The next 8 zero bytes consist of 4 bytes for a section link, and 4 for section info. Both are unused. Finally we have `0x0000_0004` which is the required alignment of the section and an unused 4 byte zero section.

#### `.strtab` Section

```text
$ od -v -Ax -tx1 -j4428 -N40 ret1234.elf
00114c 09 00 00 00 03 00 00 00 00 00 00 00 00 00 00 00
00115c 70 10 00 00 08 00 00 00 00 00 00 00 00 00 00 00
00116c 01 00 00 00 00 00 00 00
```

The first 4 bytes are the index into the `.shstrtab` section which has the header's name - `.strtab`. The second 4 bytes give the section type, which is `0x3` in this case for a `STRTAB` or string table.

We then have all zeroes for the section attributes and section virtual address which aren't relevant for a string table.

The section's offset in the file is given next: `0x0000_1070` (4208), followed by the size of the section (`0x8` or 8 bytes). That's enough for us to dump the section, which is in [7].

There are 8 unused bytes, followed by an alignment section of `0x1` and then 4 more unsude bytes.

#### `.symtab` Section

```text
$ od -v -Ax -tx1 -j4388 -N40 ret1234.elf
001124 01 00 00 00 02 00 00 00 00 00 00 00 00 00 00 00
001134 30 10 00 00 40 00 00 00 04 00 00 00 03 00 00 00
001144 04 00 00 00 10 00 00 00
```

The first 4 bytes are the index into the `.shstrtab` section which has the header's name - `.symtab`. The second 4 bytes give the section type, which is `0x2` in this case for a `SYMTAB` or symbol table.

We then have all zeroes for the section attributes and section virtual address which aren't relevant for a symbol table.

The section's offset is `0x0000_1030` with a length of `0x40`. It's dumped in [8].

The index of an associated section is `0x0000_0004`, which references the `.strtab` section. The info section is `0x0000_0003`, which is section-specific. In this case, it's the index in the symbol table which holds the first non-local symbol.

Finally the alignment is `0x0000_0004` and `0x0000_0010` is the size of each symbol table entry.

#### `.shstrtab` Section

```text
$ od -v -Ax -tx1 -j4468 -N40 ret1234.elf
001174 11 00 00 00 03 00 00 00 00 00 00 00 00 00 00 00
001184 78 10 00 00 33 00 00 00 00 00 00 00 00 00 00 00
001194 01 00 00 00 00 00 00 00
```

The first 4 bytes are the index into the `.shstrtab` section which give the header's name - which is unsurprisingly `.shstrtab`. The second 4 bytes give the section type, which is `0x3` in this case, for a `STRTAB` or string table.

We then have all zeroes for the section attributes and section virtual address which aren't relevant for a text section.

The section's offset in the file is given next: `0x0000_1078` (4216), followed by the size of the section (`0x33` or 51 bytes). That's enough for us to dump the section, which is in [6].

The next 8 zero bytes consist of 4 bytes for a section link, and 4 for section info. Both are unused. Towards the end, we have `0x0000_0001` which is the required alignment of the section and an unused 4 byte zero section.

### Sections

#### .symtab Section

```text
# First entry is always null and ignored
00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
```

```text
# Second entry - not super important
00 00 00 00  # string table name pointer, 0 means no name
00 00 40 20  # value of the symbol, 0x2040_0000
00 00 00 00  # size of the symbol
03  # info: SECTION symbol type, LOCAL bind type
00  # default visibility
01 00  # section header table index`
```

```text
# 0x13D: Third Entry
00 00 00 00 00 00 00 00 00 00 00 00 03 00 02 00

# 0x14D: Fourth Entry
01 00 00 00 00 00 40 20 00 00 00 00 10 00 01 00
```

## Notes

[1] We see a call to `load_elf` which ultimately ends up [load\_elf\_ram\_sym](https://github.com/riscv/riscv-qemu/blob/32a1a94dd324d33578dca1dc96d7896a0244d768/hw/core/loader.c#L461) - the ELF parsing logic might be useful to us later so it's noted here, but in any case it's neat to drill down into these things.

[2] The alternative to generating ELF files for both platforms is to patch Qemu's HiFive-compatible device to support raw binary files, perhaps with a specified offset. We could for example patch the `load_kernel` function to parse `-kernel "0x20400000:ret1234.bin"` as "place ret1234.bin at `0x2040_0000`". It's a judgement call, but it feels like patching Qemu is further away from what we're trying to do than just banging out an ELF header, which is what we'll be doing later.

[3] An analysis of the ELF header in `ret1234.elf`:

```text
od -Ax -tx1 -N52 ret1234.elf
000000 7f 45 4c 46 01 01 01 00 00 00 00 00 00 00 00 00
000010 02 00 f3 00 01 00 00 00 00 00 40 20 34 00 00 00
000020 ac 10 00 00 00 00 00 00 34 00 20 00 01 00 28 00
000030 06 00 05 00
```

 There's a 4 byte magic number, then `01 01 01` denoting a 32-bit, little endian, version 1 ELF file. The 2 OS ABI hint bytes are unused, as are 7 bytes of padding. After the 16th byte, the endianness of the data swaps to whatever endianness is denoted in the header. So, `02 00` is `0x0002`, and so on.

At 0x10 we see `0x0002` to mark an executable ELF file and then `F3 00` to mark a RISC-V ELF. There's another ELF version - 4 bytes this time - and then our code's entry point: `0x20400000`.

`0x00000034` denotes the address of the ELF program header table, which follows the ELF header - `0x34` equals `52` which is the length of the ELF header. The ELF "section header table" pointer is `0x000010ac`.

A 4-byte flag field is unused, and followed by `0x0034` for the ELF header size again and `0x0020` (`32` in decimal) for the size of a program header table entry. `0x0001` is the number of entries in the program header table.

Next we have `0x0028` which is the size of an entry in the section header table, and `0x0006` which is the count of entries in that table. Finally, `0x0005` has the "index of the section header table entry that contains the section names" (from Wikipedia).

[4] An analysis of the program header section:

```text
$ od -Ax -tx1 -N32 -j52 ret1234.elf
000034 01 00 00 00 00 10 00 00 00 00 40 20 00 00 40 20
000044 0c 00 00 00 0c 00 00 00 05 00 00 00 00 10 00 00
```

`01 00 00 00` - `0x00000001` identifies a LOADable segment

`00 10 00 00` - `0x00001000` is the offset of the segment in this file image

`00 00 40 20` - `0x20400000` is the virtual address of the segment

`00 00 40 20` - `0x20400000` is the physical address (the same as the virtual address in our case here)

`0c 00 00 00` - `0x0000000c` is the size of the segment in this file (12 bytes, or 3 instructions here)

`0c 00 00 00` - `0x0000000c` is the size of the segment in memory (same as the size in the file here)

`05 00 00 00` - `0x00000005` are some segment-specific flags. `5 == 1 | 4` which indicates that the segment is readable (`1`) and executable (`4`)

`00 10 00 00` - `0x00001000` is an alignment helper; the virtual address should equal the offset modulo this value.

[5] The `.text` section header

[6] The `.shstrtab` section:
We use `-ta` to dump ASCII chars for the section itself at `0x1078`, rather than hex bytes:

```text
$ od -Ax -ta -j4216 -N51 ret1234.elf
001078 nul   .   s   y   m   t   a   b nul   .   s   t   r   t   a   b
001088 nul   .   s   h   s   t   r   t   a   b nul   .   t   e   x   t
001098 nul   .   r   i   s   c   v   .   a   t   t   r   i   b   u   t
0010a8   e   s nul
```

This illustrates that `.shstrtab` is a list of null-terminated strings which can be indexed into by section headers.

[7] The `.strtab` section:

```text
od -Ax -ta -j4208 -N8 ret1234.elf
001070 nul   _   s   t   a   r   t nul
```

[8] The `.symtab` section:

```text
$ od -Ax -tx1 -j4144 -N64 ret1234.elf
001030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
001040 00 00 00 00 00 00 40 20 00 00 00 00 03 00 01 00
001050 00 00 00 00 00 00 00 00 00 00 00 00 03 00 02 00
001060 01 00 00 00 00 00 40 20 00 00 00 00 10 00 01 00
```

[7] Running `readelf` on `ret1234.elf` gives the following:

```text
$ $RISCV_PREFIX/bin/riscv32-unknown-elf-readelf -a ret1234.elf
ELF Header:
  Magic:   7f 45 4c 46 01 01 01 00 00 00 00 00 00 00 00 00
  Class:                             ELF32
  Data:                              2's complement, little endian
  Version:                           1 (current)
  OS/ABI:                            UNIX - System V
  ABI Version:                       0
  Type:                              EXEC (Executable file)
  Machine:                           RISC-V
  Version:                           0x1
  Entry point address:               0x20400000
  Start of program headers:          52 (bytes into file)
  Start of section headers:          4268 (bytes into file)
  Flags:                             0x0
  Size of this header:               52 (bytes)
  Size of program headers:           32 (bytes)
  Number of program headers:         1
  Size of section headers:           40 (bytes)
  Number of section headers:         6
  Section header string table index: 5

Section Headers:
  [Nr] Name              Type            Addr     Off    Size   ES Flg Lk Inf Al
  [ 0]                   NULL            00000000 000000 000000 00      0   0  0
  [ 1] .text             PROGBITS        20400000 001000 00000c 00  AX  0   0  4
  [ 2] .riscv.attributes RISCV_ATTRIBUTE 00000000 00100c 000024 00      0   0  1
  [ 3] .symtab           SYMTAB          00000000 001030 000040 10      4   3  4
  [ 4] .strtab           STRTAB          00000000 001070 000008 00      0   0  1
  [ 5] .shstrtab         STRTAB          00000000 001078 000033 00      0   0  1
Key to Flags:
  W (write), A (alloc), X (execute), M (merge), S (strings), I (info),
  L (link order), O (extra OS processing required), G (group), T (TLS),
  C (compressed), x (unknown), o (OS specific), E (exclude),
  p (processor specific)

There are no section groups in this file.

Program Headers:
  Type           Offset   VirtAddr   PhysAddr   FileSiz MemSiz  Flg Align
  LOAD           0x001000 0x20400000 0x20400000 0x0000c 0x0000c R E 0x1000

 Section to Segment mapping:
  Segment Sections...
   00     .text

There is no dynamic section in this file.

There are no relocations in this file.

The decoding of unwind sections for machine type RISC-V is not currently supported.

Symbol table '.symtab' contains 4 entries:
   Num:    Value  Size Type    Bind   Vis      Ndx Name
     0: 00000000     0 NOTYPE  LOCAL  DEFAULT  UND
     1: 20400000     0 SECTION LOCAL  DEFAULT    1
     2: 00000000     0 SECTION LOCAL  DEFAULT    2
     3: 20400000     0 NOTYPE  GLOBAL DEFAULT    1 _start

No version information found in this file.
Attribute Section: riscv
File Attributes
  Tag_RISCV_arch: "rv32i2p0_m2p0_a2p0"
```
