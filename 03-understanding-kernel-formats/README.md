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

We'll dump the first 52 bytes of ret1234.elf, which is the ELF header that can be compared against the detail given in the Wikipedia article on the ELF format. This header is further analysed in [3] and a full dump of the file using `readelf` is given in [4]. Despite all that, we can ignore most of the details.

```text
od -Ax -tx1 -N52 ret1234.elf
000000 7f 45 4c 46 01 01 01 00 00 00 00 00 00 00 00 00
000010 02 00 f3 00 01 00 00 00 00 00 40 20 34 00 00 00
000020 ac 10 00 00 00 00 00 00 34 00 20 00 01 00 28 00
000030 06 00 05 00
```

`.riscv.attributes` is a special section which contains details about the specific RISC-V architecture the object was compiled for - `rv32ima` in our case - among other things. It's not important for us; all we need to get started is a `.text` section which gets put at the address `0x2040_0000`.

## Running and Debugging Code

Now we're aware of the file formats needed to get runnable code on both platforms, let's actually run some code! For now we'll use `ret1234.elf` as created by the GNU toolchain's linker, but we'll soon create a bare-metal ELF file.

### Running on Qemu

We've already seen how to run something in qemu but we'll do the same again here:

```bash
$ qemu-system-riscv32 -machine sifive_e -nographic -s -S -kernel ret1234.elf
# in another terminal...

$ $RISCV_PREFIX/bin/riscv32-unknown-elf-gdb
(gdb) target remote :1234
Remote debugging using :1234
0x00001000 in ?? ()
(gdb) x 0x20400000
0x20400000:    0x012340b7
(gdb) x 0x20400004
0x20400004:    0x00100073
(gdb) x 0x20400008
0x20400008:    0x0000006f
(gdb) x 0x2040000c
0x2040000c:    0x00000000
(gdb) i r x1
x1    0x0    0x0
(gdb) nexti
0x00001004 in ?? ()
(gdb) nexti
0x20400004 in ?? ()
(gdb) nexti
^C
Program received signal SIGINT, Interrupt.
0x00000000 in ?? ()
(gdb) i r x1
x1             0x1234000           0x1234000
```

That's the value we expected, so all is well for Qemu. Let's try real hardware!

### Running on HiFive1

TODO: Run with openocd/gdb and inspect `x1` to check output

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

[4] Running `readelf` on `ret1234.elf` gives the following:

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
