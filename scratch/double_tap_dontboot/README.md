# "Double Tap Dontboot" and 0x20000000 Dump

The directory has the `double_tap_dontboot.c` bootloader provided by Sifive alongside a dump of the `double_tap_dontboot.c` bootloader which is distributed on the HiFive1 at address `0x2000_0000`.

The necessary includes are provided, although are unlikely to be sufficient to recreate the full bootloader without a linker script.

This is handy as a reference of what exists on the board and as an audit of what exists.

To create the dump, use OpenOCD:

```bash
openocd -f openocd.cfg -c "dump_image dump.img 0x20000000 0x3f0; exit;"
$RISCV_PREFIX/bin/riscv32-unknown-elf-objdump -D 0x20000000_dump.img -b binary -m riscv:rv32 > 0x20000000_dump.list
```
