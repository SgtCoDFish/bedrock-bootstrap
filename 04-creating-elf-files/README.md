# Creating Bare Metal ELF Files

## Writing a Minimal ELF Header

We've already seen that we can use `echo` to write binary files; for example, `echo -ne "\x30"` writes the value `0x30` (i.e. decimal 48, which is an ASCII `0`) to stdout. However, having to prefix everything with `\x` would be a pain, so we need an easier way to write hex and convert to binary.

`xxd -r p` can do this. `-r` means "reverse", i.e. convert hex into binary. `-p` treats the input as a plain hexdump. For example, `echo -n "30" | xxd -r -p` gives us `0` as we'd expect.

The command will also handle whitespace, so we can type `echo -n "00 00 40 20" | xxd -r -p` which will write the value of `0x20400000`. Note that we have to handle the endianness of the data ourselves.

It's also useful to be able to comment our hex files; we can use `sed` to remove the comments as a processing step. We use `sed "s/#.*$//g"` to replace everything from a # to the end of a line with nothing, repeated globally inside the file.

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