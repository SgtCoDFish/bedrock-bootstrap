# OSDev

This folder contains code from the wiki.osdev.org [bare bones](https://wiki.osdev.org/HiFive-1_Bare_Bones) tutorial, which is a handy test bed uploads of real code to hardware.

You can upload main.img to a HiFive1 using the following command:

```bash
${OPENOCD} -f openocd.cfg \
  -c "flash protect 0 0 last off; \
      program main.img verify;\
      resume 0x20400000;\
      exit"
```

See the OpenOCD upload guide for more details on the command.
