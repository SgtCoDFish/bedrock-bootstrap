addi x0, x0, 0x0
addi x0, x0, 0x0
addi x0, x0, 0x0
addi x0, x0, 0x0

# Initialise UART
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

# Pad to 0x80
addi x0, x0, 0x0
addi x0, x0, 0x0


# Register allocation:
# t0 / x05: current word
# t1 / x06: shift for placement in word
# t2 / x07: byte shift, starts at 4
# s0 / x08: scratch register used for comparison values
# s1 / x09: scratch register used for comparison values
# a0 / x10: current received byte from UART
# a1 / x11: tmp scratch area for reading from UART
# a2 / x12: 0x8000_0000 (used for reading/writing UART)
# a4 / x14: RAM write pointer, starts at 0x8000_0000
# a5 / x15: UART0_TXDATA = 0x10013000 (write UART)
# a6 / x16: UART0_RXDATA = 0x10013004 (receive UART)
# a7 / x17: 0xa (used for comparing in comment mode)
# s2 / x18: 0x20 (used for comparing x06)
# s3 / x19: current little endian byte
# s4 / x20: flags, where nonzero means we're in comment mode
# s5 / x21: the location to which we jump when we're done
#            accepting input. same as the initial value of
#            x14

addi x5, x0, 0x0
addi x6, x0, 0x0
addi x7, x0, 0x4
addi x8, x0, 0x0
addi x9, x0, 0x0
addi x10, x0, 0x0
addi x11, x0, 0x0
lui x12, 0x80000
lui x14, 0x80000
lui x15, 0x10013
addi x16, x15, 0x4
addi x17, x0, 0xa
addi x18, x0, 0x20
addi x19, x0, 0x0
addi x20, x0, 0x0
addi x21, x14, 0x0

# b0: .READ_UART
# Read a byte from UART
.READ_UART:
lw x10, 0(x16)
and x11, x10, x12
bne x11, x0, -8
andi x10, x10, 0xFF

addi x0, x0, 0x0
addi x0, x0, 0x0
addi x0, x0, 0x0

# If we're in comment mode, ignore everything that isn't 0xa
# Since we only have 1 flag at the moment, we just skip comment mode
# if the flags register is zero

beq x20, x0, 0x24

# d0: Comment mode: skip anything which isn't 0xa
bne x10, x17, +0x8
addi x20, x0, 0x0
addi x0, x0, 0x0
beq x0, x0, .READ_UART

addi x0, x0, 0x0
addi x0, x0, 0x0
addi x0, x0, 0x0
addi x0, x0, 0x0

# f0: If we're here, we're not in comment mode
# We need to check if we got 0x23 (#) and if so,
# enter comment mode

addi x8, x0, 0x23
bne x10, x8, 0x14

# We got a #, so enter comment mode
addi x20, x0, 0x1
beq x0, x0, .READ_UART # (-0x4C)
addi x0, x0, 0x0
addi x0, x0, 0x0

# 108: If we didn't get 0x23, then we can ignore anything below 0x30 ('0')
addi x8, x0, 0x30
bge x10, x8, +0x10
addi x0, x0, 0x0
addi x0, x0, 0x0
beq x0, x0, .READ_UART

# The value is >= 0x30, so it could be a hex char
# If the value is less than 0x3a, we just subtract 0x30 and that's the value
addi x8, x0, 0x3a
bge x10, x8, +0x14

# We got a value in the range 0x30 - 0x39 inclusive. Subtract 0x30 and we have our raw value
addi x10, x10, -0x30
addi x0, x0, 0x0
addi x0, x0, 0x0
beq x0, x0, .ADDRAW  # +0x40

# Otherwise, we got a value >= 0x3A, so it could be a lower or uppercase letter.
# We OR the x10 register with 0x20 to force lowercase, and then check that the value is in
# the range 0x61 ('a') to 0x66 ('f'). If it is, we subtract 0x57 to get the raw value
# in the range 0xA - 0xF

ori  x10, x10, 0x20
addi x8, x0, 0x61
blt  x10, x8, .PANIC  # +0x90
addi x8, x0, 0x67
blt  x8, x10, +0xC  # (check for `j`)
addi x10, x10, -0x57
beq x0, x0, .ADDRAW  # +0x24

# Check if we got a "j" - if so, jump to the starting point where we began writing
# instructions we received, which is stored in x21

addi x8, x0, 0x6A
bne x10, x8, .PANIC  # + 0x78

# If we're here, we need to jump to the location in x21
# but first, we need to write an infinite loop instruction to the location stored in x14
# so that we don't jump to x21, execute a few instructions and then hit an invalid
# 0000_0000 instruction and trap

addi x10, x0, 0x63
sw x10, 0(x14)
jalr x0, x21, 0x0
addi x0, x0, 0x0
addi x0, x0, 0x0
addi x0, x0, 0x0

# 170: ADDRAW: we got a value to be added to our "running word" which will be written
# when it's complete.
# We shift left by some multiple of 4 bits and OR the resulting value with the
# running value in x05. If we shifted left by K, we reset the shift amount and write
# x5 to the RAM at x14 and increment x14 by 4.
# After all steps, we jump back and read another byte from UART.

# There's some nuance to handling little endian incoming ASCII,
# since the shifts aren't linear. To explain that, consider:
# We get the input "4a 8b 01 3c" which should translate as 0x3c018b4a
# The shifts are:
# 4 << 04 , a << 00
# 8 << 12 , b << 08
# 0 << 20 , 1 << 16
# 3 << 24 , c << 28
# So we have to build bytes (which we store in x19) by storing a shift
# in x07, which switches between 4 and 0, starting at 4.
# We get a nibble from uart, shift it by the number of places in x07
# and then OR it with x19.
# If x07 == 4, we set x07 to 0 and read another value from UART
# If x07 == 0, we:
# - add 4 to x07
# - shift x19 left by x06 places
# - x05 |= x19
# - x19 = 0
# - add 8 to x06
# - if x06 == 32:
# -- set x06 to 0
# -- write the byte in x05 to the address in x14
# -- x05 = 0
# -- add 4 to x14
# - finally, read another value from UART

.ADDRAW:

sll x10, x10, x7
or x19, x19, x10
beq x7, x0, +0x14
addi x7, x0, 0x0
beq x0, x0, .READ_UART # -0xb0
addi x0, x0, 0x0
addi x0, x0, 0x0

# 18c: ADDBYTE
addi x7, x0, 0x4
sll x19, x19, x6
or x5, x5, x19
addi x19, x0, 0x0
addi x6, x6, 0x8

beq x6, x18, 0x10
beq x0, x0, .READ_UART  # -0xf0
addi x0, x0, 0x0
addi x0, x0, 0x0

# 1b0: WRITEWORD

sw x5, 0(x14)
addi x14, x14, 0x4
addi x5, x0, 0x0
addi x6, x0, 0x0
addi x0, x0, 0x0
addi x0, x0, 0x0
beq x0, x0, .READ_UART  # -0x118

# 1cc: PANIC
.PANIC:
addi x0, x0, 0x0
addi x0, x0, 0x0
addi x0, x0, 0x0
addi x0, x0, 0x0

beq x0, x0, 0 # hang
