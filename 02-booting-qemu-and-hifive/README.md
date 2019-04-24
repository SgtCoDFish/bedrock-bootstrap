# Booting on QEMU and the HiFive1

So we know that after the first bootloader, control will pass to 0x20400000 immediately and then everything hangs. Why does it hang, and does the HiFive1 do the same thing?

## Aims

- Learn how the boot process works on both QEMU and the HiFive1
- Work out how we can run the same code on both platforms

## Why QEMU Hangs

```bash
(gdb) x 0x20400000
0x20400000:    0x00000000
```

What's at `0x2040_0000`? Absolutely nothing! `0x00000000` is an illegal instruction, which causes the process to trap and thereby sets the PC to `0x00000000`... which in RISC-V always contains `0x0` by definition and so causes an infinite loop! We'll get our code running later, once we've figured out how to get to this point on actual hardware. For now we can summarise the QEMU boot process:

- Start at `0x0001_0000`
- Jump to  `0x2040_0000`

The process is very simple, and the upshot is we need to write our code at `0x2040_0000`.

## Booting on the HiFive1

If you already read the HiFive1 documentation regarding the boot process, you'll have noticed that the process for the HiFive1 is slightly different. The HiFive1 comes with slightly more ceremony with the intention of making it easier to develop for and customise in the future.

If you've run a HiFive1 from the factory, you'll also notice that it _does_ do things out of the box, unlike our QEMU image.

We can get the more detail in the following lightly edited description of the boot process from [a datasheet](https://sifive.cdn.prismic.io/sifive%2Ffeb6f967-ff96-418f-9af4-a7f3b7fd1dfc_fe310-g000-ds.pdf)[1]. The sections in brackets have been added.

> The FE310-G000 \[starts at address `0x0001_0000`\] and boots by jumping to the beginning of OTP memory \[at `0x0002_0000`\] and executing code found there. As shipped, OTP memory at the boot location is preprogrammed to jump immediately to the end of the OTP memory \[`0x0002_1FF4`\], which contains the following code to jump to the beginning of the SPI-Flash at `0x2000_0000`:

```assembly
0x0002_1FF4:
   0:    0000000f    fence 0,0
   4:    200002b7    li t0, 0x20000000
   8:    00028067    jr t0
```

(Note that "OTP" means "one time programmable" memory - that is, once you "burn" a program there, it's there permanently. "SPI Flash" means flash memory connected over [SPI](https://en.wikipedia.org/wiki/Serial_Peripheral_Interface))

What matters most from the above text is the code; `fence 0,0` is described in [this StackOverflow answer](https://stackoverflow.com/a/26374650) and is provided essentially to enable a neat trick[2]. We can ignore it as a no-op.

`li t0,0x20000000` and `jr t0` look very similar to the QEMU code we saw above and there aren't any surprises: we load `0x20000000` into `t0` and then jump there.

What happens at in SPI-Flash at `0x20000000`? As shipped, there's [a program](https://github.com/sifive/freedom-e-sdk/tree/f9271b91257e0a8a989faf3eff0757ee46694fe0/software/double_tap_dontboot) written there, whose source is reproduced in the `scratch` directory of this project along with a dump of the bootloader from real hardware. We might choose to overwrite that bootloader code later, but for now we can basically just take it as a given that it jumps to `0x2040_0000`.

To summarise, the HiFive1 boots like this:

- Start at `0x0001_0000` (this is the same as QEMU)
- Jump to  `0x0002_0000` (start of OTP)
- Jump to  `0x0002_1FF4` (end of OTP)
- Jump to  `0x2000_0000` (start of SPI-Flash)
- As shipped: Do some init and magic stuff, then jump to `0x2040_0000`

## Running Code on Both Platforms

We've established that on both platforms, our code needs to live at `0x2040_0000`, even if there's a little magic required to get there on the HiFive1. In the next section, we'll start to get something of our own running there.

## Notes

[1] (Note also that that datasheet is actually old, but the description of the boot process is in some ways easier to follow. More "up to date" details are available [here](https://sifive.cdn.prismic.io/sifive%2F4d063bf8-3ae6-4db6-9843-ee9076ebadf7_fe310-g000.pdf)).

[2] From the same datasheet's boot description:
> `fence 0,0` is encoded as `0x0000000F`, and the instruction may be modified by burning additional bits to transform it into a `JAL` instruction (opcode `0x6F`) to execute arbitrary code rather than jumping directly to the beginning of the SPI-Flash.

This means that it's easy to change `0x0000_000F` into a different instruction by burning bits. We won't be doing anything so permanent any time soon!
