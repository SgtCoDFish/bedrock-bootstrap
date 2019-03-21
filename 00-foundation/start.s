.globl _start
_start:
	lui x2,0x80004
	jal main
	sbreak
	j .
