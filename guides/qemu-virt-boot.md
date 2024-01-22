Virt board boot

`riscv_find_and_load_firmware` calls `riscv_find_firmware` to get the filename to load, which will be empty if `-bios none`.

If firmware isn't empty it calls `riscv_load_firmware` which returns the end address for the firmware.

The end address will be set to the load address if there was no firmware, which in turn defaults to `memmap[VIRT_DRAM].base` (0x80000000).
