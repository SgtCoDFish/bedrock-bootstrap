addi x0, x0, 0x00  # nop
lui a5, 0x10012
addi a5, a5, 0x3c
lw a0, 0(a5)
addi a1, x0, 0x3
slli a1, a1, 0x10
sub a1, x0, a1
and a0, a1, a0
sw a0, 0(a5)
addi a5, a5, -0x8
lw a0, 0(a5)
sub a1, x0, a1
or a0, a0, a1
lui a5, 0x10013
addi a5, a5, 0x18
addi a0, x0, 0x8a
sw a0, 0(a5)
addi a5, a5, -0x10
addi a0, x0, 0x1
sw a0, 0(a5)
addi a5, a5, 0x4
addi a0, x0, 0x1
sw a0, 0(a5)
addi x0, x0, 0x00
addi x0, x0, 0x00
addi a5, a5, -0xC
addi a1, x0, 0x35
lui a2, 0x80000
lw a0, 0(a5)
and a0, a0, a2
bne a0, x0, -8
sw a1, 0(a5)
addi x0, x0, 0x00
jal x0, 0x00
