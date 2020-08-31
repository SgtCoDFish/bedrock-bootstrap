# Using OpenOCD

OpenOCD (Open On-Chip Debugger) is a tool which allows us to debug and update a device which is connected to a computer. In the case of the HiFive1, we can connect it to a machine via a micro-USB to USB cable
and then use GDB to write new images to the on-device flash memory, which can be booted from later.

## Running OpenOCD Without `sudo`

We need to add our user to the relevant group and allowlist the specific device vendor and product IDs for the HiFive1.

Add the following to `/etc/udev/99-openocd-riscv.rules`:

```text
# Taken from https://forgge.github.io/theCore/guides/running-openocd-without-sudo.html

ACTION!="add|change", GOTO="openocd_rules_end"
SUBSYSTEM!="usb|tty|hidraw", GOTO="openocd_rules_end"

# idVendor and idProduct taken from https://github.com/sifive/freedom-e-sdk/blob/cebbb7cb6c7b16767b4ba04e7f231a784be9697a/bsp/sifive-hifive1/openocd.cfg#L5
# vid = vendor id
# pid = product id
ATTRS{idVendor}=="0403", ATTRS{idProduct}=="6010", MODE="664", GROUP="plugdev"

LABEL="openocd_rules_end"
```

```bash
sudo useradd -G plugdev $(whoami)
sudo udevadm trigger  # this might change on different distros, but works on Debian derivatives
```

This should be sufficient although you might need to log out and log back in first.

## Writing Using OpenOCD

To write an ELF file using OpenOCD you'll need a device definition file, which is provided by SiFive in the [freedom-e-sdk repository](https://github.com/sifive/freedom-e-sdk/blob/cebbb7cb6c7b16767b4ba04e7f231a784be9697a/bsp/sifive-hifive1/openocd.cfg).

To upload an ELF file, use the following command. Note that if you've not allowlisted your user to use OpenOCD, you'll need to call OpenOCD with `sudo`.

```bash
${OPENOCD} -f openocd.cfg \
  -c "flash protect 0 0 last off; \
      program main.img verify;\
      resume 0x20400000;\
      exit"
```

That's enough to write the file!

### Understanding the OpenOCD Commands

#### `flash protect 0 0 last off`

Disables flash protection on flash bank 0 (that's the first `0`) starting from the first block (the second `0`) and finishing on the `last` block. That ensure that the flash is writable.

We can use `${OPENOCD} -f openocd.cfg -c "flash info 0; exit;"` to show us information about flash bank 0:

```text
...
#0 : fespi at 0x20000000, size 0x01000000, buswidth 0, chipwidth 0
    #  0: 0x00000000 (0x10000 64kB) not protected
    #  1: 0x00010000 (0x10000 64kB) not protected
    #  2: 0x00020000 (0x10000 64kB) not protected
    #  3: 0x00030000 (0x10000 64kB) not protected
    #  4: 0x00040000 (0x10000 64kB) not protected
    #  5: 0x00050000 (0x10000 64kB) not protected
    #  6: 0x00060000 (0x10000 64kB) not protected
    #  7: 0x00070000 (0x10000 64kB) not protected
    #  8: 0x00080000 (0x10000 64kB) not protected
    #  9: 0x00090000 (0x10000 64kB) not protected
    # 10: 0x000a0000 (0x10000 64kB) not protected
...
    #250: 0x00fa0000 (0x10000 64kB) not protected
    #251: 0x00fb0000 (0x10000 64kB) not protected
    #252: 0x00fc0000 (0x10000 64kB) not protected
    #253: 0x00fd0000 (0x10000 64kB) not protected
    #254: 0x00fe0000 (0x10000 64kB) not protected
    #255: 0x00ff0000 (0x10000 64kB) not protected
...
```

We see that the flash bank is located at `0x2000_0000` and is `0x0100_0000` bytes long which is ~16MB. This 16MB is divided into 256 64kB blocks. That's what we'd expect to see, based on the information in the documentation.

#### `program main.img verify`

Write and verify the file called `main.img`. In the [OpenOCD documentation](http://www.openocd.org/doc/html/Flash-Programming.html#Flash-Programming) we see that `program` is a shortcut for a few other commands including `write_image` which actually does the write.

Since `main.img` is an ELF file, it contains the addresses at which code should be written as instructed by a linker script. It's possible to use `write_image` along with an explicit offset to write a raw binary file.

#### `resume 0x20400000`

Starts the core running again from the given address. This is where our `_start` function usually lives.
