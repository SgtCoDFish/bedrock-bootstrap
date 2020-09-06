# Longan Nano GD32VF103

A small UART test app using platformio for the longan nano

Based on the [gd32v\_lcd example](https://github.com/sipeed/Longan_GD32VF_examples/tree/7fe21406e73dbe49e9aa3c5a9cd24a44af3b5f41/gd32v_lcd)

## Firmware Backup

Backing up original firmware:

```bash
sudo dfu-util -d 28e9:0189 -U firmware_backup
```

## Links

- Original firmware [dump](./firmware_backup)
- Overall details in a good article: [susa.net](https://www.susa.net/wordpress/2019/10/longan-nano-gd32vf103/)
- Patched [dfu-utils](https://github.com/riscv-mcu/gd32-dfu-utils)
