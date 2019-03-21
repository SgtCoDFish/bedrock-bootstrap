OUTPUT_ARCH( "riscv" )
ENTRY(_start)
MEMORY
{
	ram : ORIGIN = 0x80000000, LENGTH = 0x4000
}

SECTIONS
{
	.text : { *(.text*) } > ram
	.rodata : { *(.rodata*) } > ram
	.bss : { *(.bss*) } > ram
}
