addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00

# Initialise UART; see comments in 05-uart/uart.hex for a detailed explanation
# of what these instructions are doing.
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

# Register allocation:
# t0 / x05: current word
# t1 / x06: shift for placement in word
# t2 / x07: byte shift, starts at 4
# s0 / x08: scratch register used for comparison values
# s1 / x09: scratch register used for comparison values
# a0 / x10: current received byte from UART
# a1 / x11: tmp scratch area for reading from UART
# a2 / x12: 0x8000_0000 (used for reading/writing UART)
# a3 / x13: 0x8000_0ea0 (location of function "m" in memory)
# a4 / x14: RAM write pointer, starts at the value in x13 but is
#           changed whenever a new function is defined
# a5 / x15: UART0_TXDATA = 0x10013000 (write UART)
# a6 / x16: UART0_RXDATA = 0x10013004 (receive UART)
# a7 / x17: max size in bytes of each function, set to 312 (0x138)
#           this gives us room for 13 functions each with a maximum
#           of 78 instructions. The main function is "m", and can have
#           an extra 40 bytes / 10 instructions at the end if we're to
#           fit into 4096 (0x1000) bytes of RAM
# s2 / x18: 0x20 (used for comparing x06)
# s3 / x19: current little endian byte
# s4 / x20: flags, where nonzero means we're in comment mode

addi x5, x0, 0x0
addi x6, x0, 0x0
addi x7, x0, 0x4
addi x8, x0, 0x0
addi x9, x0, 0x0
addi x10, x0, 0x0
addi x11, x0, 0x0
lui x12, 0x80000
lui x13, 0x80001
addi x13, x13, -0x160
addi x14, x13, 0x00
lui x15, 0x10013
addi x16, x15, 0x4
addi x17, x0, 0x138 # 312
addi x18, x0, 0x20
addi x19, x0, 0x00
addi x20, x0, 0x00

# b0: .READ_UART
# Read a byte from UART
.READ_UART:
lw x10, 0(x16)
and x11, x10, x12
bne x10, x0, -8
andi x10, x10, 0xFF

addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00

# If the flags register is zero, we continue as normal, otherwise
# we need to check which flags are enabled and handle them as appropriate
beq x20, x0, 0x8

# d0: If the flags register is nonzero we need to check which flag is enabled
beq x0, x0, .HANDLE_FLAGS
addi x0, x0, 0x00

# Check for "." (0x2E) to start function definition mode
addi x8, x0, 0x2e
bne x10, x8, 0x10

# If we got a ".", set bit 2 of flags and jump back to READ_UART
ori x20, x20, 0x2
beq x0, x0, .READ_UART
addi x0, x0, 0x00
addi x0, x0, 0x00

# f0: If we're here, we're not in any kind of flag mode
# We need to check if we got 0x23 (#) and if so, enter comment mode
addi x8, x0, 0x23
bne x10, x8, 0x14

# We got a #, so enter comment mode
ori x20, x20, 0x1
beq x0, x0, .READ_UART
addi x0, x0, 0x00
addi x0, x0, 0x00

# 108: If we didn't get 0x23, then we can ignore anything below 0x30 ('0')
addi x8, x0, 0x30
bge x10, x8, 0x10
addi x0, x0, 0x00
addi x0, x0, 0x00
beq x0, x0, .READ_UART

# The value is >= 0x30, so it could be a hex char
# If the value is less than 0x3a, we just subtract 0x30 and that's the value
addi x8, x0, 0x3a
bge x10, x8, 0x14

# We got a value in the range 0x30 - 0x39 inclusive. Subtract 0x30 and we have our raw value
addi x10, x10, -0x30
addi x0, x0, 0x00
addi x0, x0, 0x00
beq x0, x0, .ADDRAW

# Otherwise, we got a value >= 0x3A, so it could be a lower or uppercase letter.
# We OR the x10 register with 0x20 to force lowercase, and then check that the value is in
# the range 0x61 ('a') to 0x66 ('f'). If it is, we subtract 0x57 to get the raw value
# in the range 0xA - 0xF

ori  x10, x10, 0x20
addi x8, x0, 0x61
blt  x10, x8, .PANIC
addi x8, x0, 0x67
blt  x8, x10, 0xC
addi x10, x10, 0x57
beq x0, x0, .ADDRAW

# Check if we got a "j" or "x" and handle them if needed in .HANDLE_SPECIAL
addi x0, x0, 0x00
beq x0, x0, .HANDLE_SPECIAL

addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00

# ADDRAW: we got a value to be added to our "running word" which will be written
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
beq x7, x0, 0x14
addi x7, x0, 0x00
beq x0, x0, .READ_UART
addi x0, x0, 0x00
addi x0, x0, 0x00

.ADDBYTE:
addi x7, x0, 0x4
sll x19, x19, x6
or x5, x5, x19
addi x19, x0, 0x00
addi x6, x6, 0x8

beq x6, x18, 0x10
beq x0, x0, .READ_UART
addi x0, x0, 0x00
addi x0, x0, 0x00

# 1b0: WRITEWORD
.WRITEWORD:
sw x5, 0(x14)
addi x14, x14, 0x4
addi x5, x0, 0x0
addi x6, x0, 0x0
addi x0, x0, 0x00
addi x0, x0, 0x00
beq x0, x0, .READ_UART

# 1cc: PANIC
.PANIC:
addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00
beq x0, x0, 0 # hang

# 1e0: HANDLE_FLAGS
.HANDLE_FLAGS:
addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00

# 1 is comment mode
# 2 is function definition mode; expect 1 ASCII char from a-m which should be
# lowercased, and then turned into an integer from 0-12. The value is then multiplied
# by the value in x17 to give the location at which the function is written

# check for comment mode

andi x8, x20, 0x1
beq x8, x0, 0x24

# handle comment mode

# if the byte we read is not 0xa (\n) just go back to read more
addi x8, x0, 0xa
bne x10, x8, 0x18

# otherwise, if we got 0xa we need to clear bit 1 and then go back
addi x8, x0, 0x1
xori x8, x8, -1
and x20, x20, x8

addi x0, x0, 0x00
addi x0, x0, 0x00

# 0x214: jump back after handling comment node
beq x0, x0, .READ_UART

# check function definition mode

andi x8, x20, 0x2
beq x8, x0, 0x4c

# handle function definition mode

# force the function name to lowercase, and then convert it into a number from 0-13 inclusive
ori  x10, x10, 0x20
addi x10, x10, -0x61

# if x10 is lower than 0 or higher than 13 (0xD) just skip and continue to read the next char
blt x10, x0, +0x3C
addi x8, x0, 0xd
bge x10, x8, +0x34

# otherwise, it's a valid function name so we need to set the current function register
# and then clear bit 2 of the flags register
# also needs to reset "shift" registers, etc: x19, x07, x06, x05
lui x14, 0x80000

# subtract from x10 until it's 0, adding the "function size" to x14 each time
beq x10, x0, 0x10
add x14, x14, x17
addi x10, x10, -1
beq x0, x0, -0x10

addi x5, x0, 0x0
addi x6, x0, 0x0
addi x7, x0, 0x4
addi x19, x0, 0x00

addi x8, x0, 0x2
xori x8, x8, -1
and x20, x20, x8

# 0x264: jump back after handling function definition mode
beq x0, x0, .READ_UART

# 0x268: check for "x" mode, where we'll create and write a jump instruction
andi x8, x20, 0x4
beq x8, x0, 0x44

# we're in "x" mode, so we need to assemble a jump instruction to the named function
# which is stored in x10

# force the function name to lowercase, and then convert it into a number from 0-13 inclusive
ori  x10, x10, 0x20
addi x10, x10, -0x61

# if x10 is lower than 0 or higher than 13 (0xD) just skip and continue to read the next char
blt x10, x0, +0x38
addi x8, x0, 0xd
bge x10, x8, +0x34

# we can't use branch instructions to make the calls, so we need to write instructions which:
# - load the correct address into x08 (this could take multiple instructions)
# - jump to x08 (jr x08)
# in other words, we need to write the following instructions; as such, we'll copy them across
# so we jump over them first
beq x0, x0, 0x1c

# as above, these aren't actually executed; they're for copying
lui x8, 0x80000
beq x10, x0, +0x10
add x8, x8, x17
addi x10, x10, -1
beq x0, x0, -0x10
jr x8

# we need to:
# write a li x10, *x10 instruction (i.e. an instruction which loads into x10 the current value of x10)
# load address of lui instruction above into x08
# while x08 <= address of jr instruction above:
# - load data from x08 into x09
# - write data from x09 into x14
# - increase x14 by 4

# TODO: write `li x10, *x10` instruction and loop to copy above instructions
# TODO: gonna need more space here
addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00

# 0x2b0: jump back after handling "x" mode
beq x0, x0, .READ_UART

# 0x2b4: .HANDLE_SPECIAL
# if we got a "j", jump to the starting point where we began writing
# instructions we received, which is stored in x21
.HANDLE_SPECIAL:
addi x8, x0, 0x6A # ("j")
bne x8, x10, 0x20

addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00

jalr x0, x13, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00

# check for "x" and if found, enter "insert jump" mode
addi x8, x0, 0x78 # ("x")
bne x8, x10, 0x20

addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00

ori x20, x20, 0x4
addi x0, x0, 0x00
addi x0, x0, 0x00
addi x0, x0, 0x00

beq x0, x0, .READ_UART
