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
