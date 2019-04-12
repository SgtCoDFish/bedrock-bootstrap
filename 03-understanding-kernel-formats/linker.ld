OUTPUT_ARCH( "riscv" )
ENTRY(_start)
MEMORY
{
	spiflash : ORIGIN = 0x20400000, LENGTH = 0x1FC00000
	ram : ORIGIN = 0x80000000, LENGTH = 0x4000
}

SECTIONS
{
	.text : { *(.text*) } > spiflash
	.rodata : { *(.rodata*) } > spiflash
	.bss : { *(.bss*) } > spiflash
}
