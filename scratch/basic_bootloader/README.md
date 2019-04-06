# HiFive1 Basic Bootloader

We know that our code will start running directly at `0x2040_0000` on Qemu whereas the HiFive1 has an intermediate bootloader at `0x2000_0000` which then jumps to `0x2040_0000`. If we accept that bootloader, we can always rely on our code starting at `0x2040_0000` on both systems and not worry about it.

However, if we want to have more complete control of the whole boot process, we can write a small bit of code at `0x2000_0000` which always jumps to `0x2040_0000`.

The code to do the jump is simple:

```assembly
.section .text.shim
    li t0, 0x20400000
    jr t0
```

All we need to do is ensure that those two instructions are placed at `0x2000_0000` and then our code is placed at `0x2040_0000` using a linker script, which lets us specify where code ends up in the resultant binary.

[This blog](https://shobhitsharda.wordpress.com/2011/03/11/linker-scripts-inner-concepts/) is a great introduction to the basic concepts.

There are some files placed in this directory which show a working basic bootloader.

```text
OUTPUT_ARCH( "riscv" )
ENTRY(_start)
MEMORY
{
    shimloc : ORIGIN = 0x20000000, LENGTH = 0x00400000
    spiflash : ORIGIN = 0x20400000, LENGTH = 0x1FC00000
    ram : ORIGIN = 0x80000000, LENGTH = 0x4000
}

SECTIONS
{
    .text.shim : { *(.text.shim*) } > shimloc
    .text : { *(.text*) } > spiflash
    .rodata : { *(.rodata*) } > ram
    .bss : { *(.bss*) } > ram
}
```

We define the "shimloc" to be an approximately 4MB (`0x0040_0000`) block in SPI-Flash, starting at `0x2000_0000`. Then under "SECTIONS" we say that the section `.text.shim` must be placed at the beginning of shimloc. Likewise, we place `.text` - regular code - at `0x2040_0000` where both Qemu and the HiFive1 will end up jumping to.
