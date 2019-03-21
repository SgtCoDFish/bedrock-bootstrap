# The Qemu/Hardware Shim

We know that our code will start running at `0x2040_0000` on Qemu and `0x2000_0000` on the HiFive1. So, we can add a couple of instructions to force compatibility between on _both_ systems (so that we only have to build one output binary that can run on both). We'll write code at `0x2000_0000` which always jumps to `0x2040_0000` and then write our actual code at `0x2040_0000`.

The code to do the jump is similar to what we've seen already:

```assembly
.section .text.shim
    li t0, 0x20400000
    jr t0
```

Ignore the `.section`, that's used later.

All we need to do is ensure that those two instructions are placed at `0x2000_0000` and then our code is placed at `0x2040_0000`. Great; how do we do that? The answer (at least for now) is a linker script, which lets us specify where code ends up. [This blog](https://shobhitsharda.wordpress.com/2011/03/11/linker-scripts-inner-concepts/) is a great introduction to the basic concepts, but in fact, we already used a linker script earlier, in the form of `linker.ld` from the foundational section which happened to put our program into RAM at `0x8000_0000`.

There are some files placed in this directory which make up a "shim test" for us to check that our shim does what we expect. The linker script is especially important:

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

In practise, we place `start.s` at the start of our main `.text` section. In this example, we load the value `0x1234` into `x1`. Run the program in qemu and confirm that that value ends up in `x1`; now we can actually work out how to get the code onto hardware to test, since we have some confidence about the boot process.
