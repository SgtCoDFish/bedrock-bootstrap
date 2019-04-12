# Creating Bare Metal ELF Files

## Writing a Minimal ELF Header

We've already seen that we can use `echo` to write binary files; for example, `echo -ne "\x30"` writes the value `0x30` (i.e. decimal 48, which is an ASCII `0`) to stdout. However, having to prefix everything with `\x` would be a pain, so we need an easier way to write hex and convert to binary.

`xxd -r p` can do this. `-r` means "reverse", i.e. convert hex into binary. `-p` treats the input as a plain hexdump. For example, `echo -n "30" | xxd -r -p` gives us `0` as we'd expect.

The command will also handle whitespace, so we can type `echo -n "00 00 40 20" | xxd -r -p` which will write the value of `0x20400000`. Note that we have to handle the endianness of the data ourselves.

The resulting hex header is in elfheader.hex, and we can "compile" it to a binary file with:

```bash
xxd -r -p elfheader.hex
```

TODO: FINISH ELF HEADER

### Program Headers

```text
$ od -Ax -tx1 -N32 -j52 ret1234.elf
0000034  01  00  00  00  00  10  00  00  00  00  40  20  00  00  40  20
0000044  0c  00  00  00  0c  00  00  00  05  00  00  00  00  10  00  00
```

`01 00 00 00` - `0x00000001` identifies a LOADable segment

`00 10 00 00` - `0x00001000` is the offset of the segment in this file image

`00 00 40 20` - `0x20400000` is the virtual address of the segment

`00 00 40 20` - `0x20400000` is the physical address (the same as the virtual address in our case here)

`0c 00 00 00` - `0x0000000c` is the size of the segment in this file (12 bytes, or 3 instructions here)

`0c 00 00 00` - `0x0000000c` is the size of the segment in memory (same as the size in the file here)

`05 00 00 00` - `0x00000005` are some segment-specific flags. `5 == 1 | 4` which indicates that the segment is readable (`1`) and executable (`4`)

`00 10 00 00` - `0x00001000` is an alignment helper; the virtual address should equal the offset modulo this value.

### Next

Then there's a whole lot of nothing until `0x1000` - this isn't unexpected given the pointer in the program header. We'll crush this empty space down later because we have no need for the extra buffer space for program headers which is left here.

We know that the section header table has entries with size `0x28` from the ELF header. Let's dump out the first section header in the table:

```text
$ od -Ax -tx1 -j4096 -N40 ret1234.elf
0001000  b7  40  23  01  73  00  10  00  6f  00  00  00  41  23  00  00
0001010  00  72  69  73  63  76  00  01  19  00  00  00  05  72  76  33
0001020  32  69  32  70  30  5f  6d  32
```

`b7 40 23 01` - `0x012340b7` is the first instruction in our program, `lui x1,0x1234`

`73 00 10 00` - `0x00100073` is the second, `ebreak`

`6f 00 00 00` - `0x0000006f` is the third, `j .`