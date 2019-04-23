# ELF Details

Analysing ELF files in detail is not the goal of this stage; nevertheless, if you want more in depth information on different ELF sections, that's provided here in a slightly free-flowing format.

## ELF File Header

An in-depth analysis of the ELF header in `ret1234.elf`, which looks as follows:

```text
od -Ax -tx1 -N52 ret1234.elf
000000 7f 45 4c 46 01 01 01 00 00 00 00 00 00 00 00 00
000010 02 00 f3 00 01 00 00 00 00 00 40 20 34 00 00 00
000020 ac 10 00 00 00 00 00 00 34 00 20 00 01 00 28 00
000030 06 00 05 00
```

There's a 4 byte magic number, then `01 01 01` denoting a 32-bit, little endian, version 1 ELF file. The 2 OS ABI hint bytes are unused, as are 7 bytes of padding. After the 16th byte, the endianness of the data swaps to whatever endianness is denoted in the header. So, `02 00` is `0x0002`, and so on.

At 0x10 we see `0x0002` to mark an executable ELF file and then `F3 00` to mark a RISC-V ELF. There's another ELF version - 4 bytes this time - and then our code's entry point: `0x20400000`.

`0x00000034` denotes the address of the ELF program header table, which follows the ELF header - `0x34` equals `52` which is the length of the ELF header. The ELF "section header table" pointer is `0x000010ac` - unsurprisingly, that's where the section header table lives.

A 4-byte flag field is unused, and followed by `0x0034` for the ELF header size again and `0x0020` (`32` in decimal) for the size of a program header table entry. `0x0001` is the number of entries in the program header table.

Next we have `0x0028` which is the size of an entry in the section header table, and `0x0006` which is the count of entries in that table. Finally, `0x0005` has the "index of the section header table entry that contains the section names" (from Wikipedia). That means the index of the `.shstrtab` section header.

## ELF Program Headers

The program header section contains details of different program sections, and looks like this:

```text
$ od -Ax -tx1 -N32 -j52 ret1234.elf
000034 01 00 00 00 00 10 00 00 00 00 40 20 00 00 40 20
000044 0c 00 00 00 0c 00 00 00 05 00 00 00 00 10 00 00
```

As we saw in the file header, there's only one program header in this file.

- `01 00 00 00` - `0x1` identifies a LOADable segment. We'll only encounter this type for our purposes.
- `00 10 00 00` - `0x1000` is the offset of the segment in this file image - this is where our actual compiled code is stored in this file.
- `00 00 40 20` - `0x20400000` is the virtual address of the segment, and is followed by the same value denoting the physical address of the segment. In our case, with no virtualised memory, these match.
- `0c 00 00 00` - `0xc` is the size of the segment in this file (12 bytes, or 3 RV32I instructions). It's only a small program!
- `0c 00 00 00` - `0xc` is the size of the segment in memory, which in our case is a repitition of the size of the segment in this file.
- `05 00 00 00` - `0x5` provides some segment-specific flags. `5 == 1 | 4` which indicates that the segment is readable (`1`) and executable (`4`).
- `00 10 00 00` - `0x1000` is an alignment helper; the virtual address should equal the offset modulo this value. In other words, `0x20400000 == 0x20401000 MOD 0x1000`.

## ELF Section Header Table

The section header table in this file contains 6 entries, which are all dumped here:

```text
# 0: NULL header inserted by the compiler, seems to be required
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

### `.text` Section

```text
$ od -v -Ax -tx1 -j4308 -N40 ret1234.elf
0010d4 1b 00 00 00 01 00 00 00 06 00 00 00 00 00 40 20
0010e4 00 10 00 00 0c 00 00 00 00 00 00 00 00 00 00 00
0010f4 04 00 00 00 00 00 00 00
```

The first 4 bytes are the index into the `.shstrtab` section which has the header's name - `.text`. The second 4 bytes give the section type, which is `0x1` in this case, denoting "program data", i.e. code.

Next `0x06` gives the flags for the section; this is equal to `0x2 | 0x4`, which means that the section "occupies memory during execution" and is "executable", respectively. `0x2040_0000` is the now-familiar address at which the section is placed.

The section's offset in the file is given next: `0x0000_1000` (4096), followed by the size of the section (`0xc` or 12 bytes). The section has already been dumped above, following the program header section.

The next 8 zero bytes consist of 4 bytes for a section link, and 4 for section info. Both are unused. Finally we have `0x0000_0004` which is the required alignment of the section and an unused 4 byte zero section.

#### `.riscv.attributes` Section

```bash
$ od -v -Ax -tx1 -j4340 -N40 ret1234.elf
0010fc 21 00 00 00 03 00 00 70 00 00 00 00 00 00 00 00
00110c 0c 10 00 00 24 00 00 00 00 00 00 00 00 00 00 00
00111c 01 00 00 00 00 00 00 00

# from readelf:
Attribute Section: riscv
File Attributes
  Tag_RISCV_arch: "rv32i2p0"
```

The first 4 bytes are the offset into the `.shstrtab`, giving the name `.riscv.attributes`. The next 4 bytes are the identifier for a RISCV_ATTRIBUTE type: `0x70000003`.

8 bytes for flags and addresses are unused, followed by the offset `0x100c` which is the location in the file and 0x24 which is the length of the section, which is dumped below.

The link and info sections make up the next 8 nul bytes, followed by an alignment of 0x1 and an "entity size" of 0.

```bash
$ od -v -Ax -tx1 -j4108 -N36 ret1234.elf
000100c    41 19 00 00 00 72 69 73 63 76 00 01 0f 00 00 00
000101c    05 72 76 33 32 69 32 70 30 00 00 00 00 00 00 00
000102c    00 00 00 00

$ od -Ax -v -ta -j4108 -N36 ret1234.elf
000100c    A  em nul nul nul   r   i   s   c   v nul soh  si nul nul nul
000101c  enq   r   v   3   2   i   2   p   0 nul nul nul nul nul nul nul
000102c  nul nul nul nul
0001030
```

Viewed in ASCII, the `.riscv.attributes` section clearly contains architectural information about the program; `rv32i2p0` is embedded in the section, and refers to 32-bit RISC-V, with only the base instruction set version 2.0 (`p` is the version separator).

From the binutils [source](https://fossies.org/linux/binutils/include/elf/riscv.h) we see that "Tag_RISCV_arch" has the value 5, which is what precedes the arch information "rv32i2p0". That's a useful thing to extract from the section, but the documentation regarding the rest seems almost completely non-existant. At the time of writing there's an open [PR](https://github.com/riscv/riscv-elf-psabi-doc/pull/71) which seems to add _some_ documentation, but it's nowhere near enough to actually be able to parse the format. We can work it out, however, from binutils which wrote the section in the first place!

The first byte is the letter 'A', which identifies the section as an "attributes" section according to [binutils](https://github.com/bminor/binutils-gdb/blob/4a4153dfc945701938b6f52795cf234fa0a5f5fe/binutils/readelf.c#L15533-L15534).

Next is a 4 byte section length (0x19), followed by a null-terminated section name "riscv\0". `01` tells us we're reading [file attributes](https://github.com/bminor/binutils-gdb/blob/4a4153dfc945701938b6f52795cf234fa0a5f5fe/binutils/readelf.c#L15648) and the 4-byte `0x0000000f` is the size of the section. `05` is the RISCV_arch tag as above, and then the null-terminated string follows. The rest of the null bytes are just waste.

#### `.strtab` Section

```text
$ od -v -Ax -tx1 -j4428 -N40 ret1234.elf
00114c 09 00 00 00 03 00 00 00 00 00 00 00 00 00 00 00
00115c 70 10 00 00 08 00 00 00 00 00 00 00 00 00 00 00
00116c 01 00 00 00 00 00 00 00
```

The first 4 bytes are the index into the `.shstrtab` section which has the header's name - `.strtab`. The second 4 bytes give the section type, which is `0x3` in this case for a `STRTAB` or string table.

We then have all zeroes for the section attributes and section virtual address which aren't relevant for a string table.

The section's offset in the file is given next: `0x0000_1070` (4208), followed by the size of the section (`0x8` or 8 bytes). That's enough for us to dump the section, which looks like this:

(We use `-ta` to dump ASCII chars for the section)

```text
od -Ax -ta -j4208 -N8 ret1234.elf
001070 nul   _   s   t   a   r   t nul
```

There are 8 unused bytes, followed by an alignment section of `0x1` and then 4 more unused bytes.

#### `.symtab` Section

```text
$ od -v -Ax -tx1 -j4388 -N40 ret1234.elf
001124 01 00 00 00 02 00 00 00 00 00 00 00 00 00 00 00
001134 30 10 00 00 40 00 00 00 04 00 00 00 03 00 00 00
001144 04 00 00 00 10 00 00 00
```

The first 4 bytes are the index into the `.shstrtab` section which has the header's name - `.symtab`. The second 4 bytes give the section type, which is `0x2` in this case for a `SYMTAB` or symbol table.

We then have all zeroes for the section attributes and section virtual address which aren't relevant for a symbol table.

The section's offset is `0x0000_1030` with a length of `0x40`. Its contents follow:

```text
$ od -Ax -tx1 -j4144 -N64 ret1234.elf
001030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00
001040 00 00 00 00 00 00 40 20 00 00 00 00 03 00 01 00
001050 00 00 00 00 00 00 00 00 00 00 00 00 03 00 02 00
001060 01 00 00 00 00 00 40 20 00 00 00 00 10 00 01 00
```

The index of an associated section is `0x0000_0004`, which references the `.strtab` section. That's associated because it's where the symbol table gets the names of the symbols.

The info section is `0x0000_0003`, which is section-specific. In this case, it's the index in the symbol table which holds the first non-local symbol.

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

The section's offset in the file is given next: `0x0000_1078` (4216), followed by the size of the section (`0x33` or 51 bytes). That's enough for us to dump the section, which looks like this:

(We use `-ta` to dump ASCII chars for the section)

```text
$ od -Ax -ta -j4216 -N51 ret1234.elf
001078 nul   .   s   y   m   t   a   b nul   .   s   t   r   t   a   b
001088 nul   .   s   h   s   t   r   t   a   b nul   .   t   e   x   t
001098 nul   .   r   i   s   c   v   .   a   t   t   r   i   b   u   t
0010a8   e   s nul
```

The next 8 zero bytes consist of 4 bytes for a section link, and 4 for section info. Both are unused. Towards the end, we have `0x0000_0001` which is the required alignment of the section and an unused 4 byte zero section.

## The Symbol Table

The symbol table looks a little arcane, but mostly it can be ignored. The null first entry is required, and then the only important entry is the last, which describes a global symbol for our program.

We'll get rid of the other two when we roll our own header.

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
01 00  # section header table index
```

```text
# 0x13D: Third Entry; an empty LOCAL section which can be ignored
00 00 00 00 00 00 00 00 00 00 00 00 03 00 02 00

# 0x14D: Fourth Entry
01 00 00 00  # string table point; 1 points to "_start"
00 00 40 20  # value of the symbol, 0x2040_0000
00 00 00 00  # size of the symbol
10  # NOTYPE type, GLOBAL bind type
00  # default visibility
01 00  # section header table index
```

## `readelf` Dump

`readelf` can be used to give fairly verbose output describing pretty much everything about an ELF file. A dump follows:

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
