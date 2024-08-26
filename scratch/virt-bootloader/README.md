# QEMU Virt Bootloader (No Firmware)

This folder contains the bootloader for the QEMU riscv64 "virt" device without firmware (`-bios none`), extracted by running
the following QEMU version locally:

```console
$ qemu-system-riscv64 -version
QEMU emulator version 9.0.2
Copyright (c) 2003-2024 Fabrice Bellard and the QEMU Project developers
```

Most of the interesting content is in `bootloader.hex` which implements the bootloader (in 64-bit RISC-V) as ripped from QEMU, using
hex. You can run `make dump` to see the bootloader commands disassembled.

The end result of the boot loader is that execution jumps to `0x8000_0000` with the registers in the following state:

```text
ra   0x0              0x0
sp   0x0              0x0
gp   0x0              0x0
tp   0x0              0x0
t0   0x80000000       -2147483648
t1   0x0              0
t2   0x0              0
fp   0x0              0x0
s1   0x0              0
a0   0x0              0
a1   0x87e00000       -2015363072
a2   0x1028           4136
a3   0x0              0
a4   0x0              0
a5   0x0              0
a6   0x0              0
a7   0x0              0
s2   0x0              0
s3   0x0              0
s4   0x0              0
s5   0x0              0
s6   0x0              0
s7   0x0              0
s8   0x0              0
s9   0x0              0
s10  0x0              0
s11  0x0              0
t3   0x0              0
t4   0x0              0
t5   0x0              0
t6   0x0              0
pc   0x80000000       0x80000000
```

Note that the bootloader sets a few registers, but in practice the only important action seems to be loading the
value at 0x1018 (which is `0x8000_0000`) into t0 and then jumping there.
