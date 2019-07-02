# Investigations

This file is mostly just a mind dump for debugging things.

## System QEMU / Built gdb

On a reasonably fresh arch install, with the `qemu-headless-arch-extra 4.0.0-2` package providing `qemu-riscv32` a bare metal program can be loaded, but cannot be debugged using a user-built version of gdb.

 gives:

```bash
$ riscv-unknown-elf-gdb --version
GNU gdb (GDB) 8.3.0.20190516-git
...
```

The error message from gdb is:

```text
bfd requires xlen 8, but target has xlen 4
```

which implies a register-size mismatch - gdb wants/expects 8 byte registers (rv64) and qemu has 4-byte registers (rv32).

The relevant file seems to be in [riscv-tdep.c](https://sourceware.org/git/gitweb.cgi?p=binutils-gdb.git;a=blob;f=gdb/riscv-tdep.c;h=bae987cf66ba4b052550a8f822ba515ba85797b4;hb=HEAD#l3156) which references the [riscv\_gdbarch\_features struct](https://sourceware.org/git/gitweb.cgi?p=binutils-gdb.git;a=blob;f=gdb/arch/riscv.h;h=05c19054dce00da97f44bb6cce00a3ee800707cd;hb=HEAD#l36).

This seems fishy; they're both meant to be 32-bit.

```bash
$ gdb --configure
This GDB was configured as follows:
   configure --host=x86_64-pc-linux-gnu --target=riscv32-unknown-elf
             --with-auto-load-dir=$debugdir:$datadir/auto-load
             --with-auto-load-safe-path=$debugdir:$datadir/auto-load
             --with-expat
             --with-gdb-datadir=.../riscv-gnu-toolchain/build/share/gdb (relocatable)
             --with-jit-reader-dir=.../riscv-gnu-toolchain/build/lib/gdb (relocatable)
             --without-libunwind-ia64
             --with-lzma
             --without-babeltrace
             --without-intel-pt
             --disable-libmcheck
             --with-mpfr
             --with-python=/usr
             --with-guile
             --disable-source-highlight
             --with-separate-debug-dir=.../riscv-gnu-toolchain/build/lib/debug (relocatable)

("Relocatable" means the directory can be moved with the GDB installation
tree, and GDB will still find it.)
```

### Fix

On further investigation, it turns out the architecture auto-discovery wasn't working. This can be fixed with "set architecture riscv:rv32" in gdb.
