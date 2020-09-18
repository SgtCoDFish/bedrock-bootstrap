# Longan Nano GD32VF103

A small UART test app using platformio for the longan nano

Based on the [gd32v\_lcd example](https://github.com/sipeed/Longan_GD32VF_examples/tree/7fe21406e73dbe49e9aa3c5a9cd24a44af3b5f41/gd32v_lcd)

## Firmware Backup

Backing up original firmware:

```bash
sudo dfu-util -d 28e9:0189 -U firmware_backup
```

## JTAG

Using the [Sipeed USB-JTAG](https://www.seeedstudio.com/Sipeed-USB-JTAG-TTL-RISC-V-Debugger-p-2910.html) you can pause and debug the Longan Nano.

Upstream openocd doesn't seem to work, but there's one in platformio: `~/.platformio/packages/tool-openocd-gd32v/bin/openocd`

Config is provided in [openocd](./openocd); both files are needed. The config is taken from [this repo](https://github.com/riscv-rust/longan-nano).

## Links

- Original firmware [dump](./firmware_backup)
- Overall details in a good article: [susa.net](https://www.susa.net/wordpress/2019/10/longan-nano-gd32vf103/)
- Patched [dfu-utils](https://github.com/riscv-mcu/gd32-dfu-utils) (possible to use the one in platformio)
