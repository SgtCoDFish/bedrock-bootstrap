addi x0, x0, 0x00  # nop

# IOF0_UART0_MASK = 0x00030000
# GPIO_IOF_SEL = 0x1001203C
# GPIO_IOF_SEL &= ~(IOF0_UART0_MASK) to enable UART

lui a5, 0x10012
addi a5, a5, 0x3c
lw a0, 0(a5)
addi a1, x0, 0x3
slli a1, a1, 0x10
sub a1, x0, a1
and a0, a1, a0
sw a0, 0(a5)

# GPIO_IOF_EN = 0x10012038
# GPIO_IOF_EN |= IOF0_UART0_MASK to enable UART
addi a5, a5, -0x8
lw a0, 0(a5)
sub a1, x0, a1
or a0, a0, a1

# UART_REG_DIV = 0x10013018
# UART_REG_DIV := 8A == 138
lui a5, 0x10013
addi a5, a5, 0x18
addi a0, x0, 0x8a
sw a0, 0(a5)

# UART_REG_TXCTRL = 0x10013008
# UART_REG_TXCTRL := 0x01 | 0x03
# TODO: Should this be set to 0x03 or 0x01?
addi a5, a5, -0x10
addi a0, x0, 0x1
sw a0, 0(a5)

# UART_REG_RXCTRL = 0x1001300C
# UART_REG_RXCTRL := 0x01
addi a5, a5, 0x4
addi a0, x0, 0x1
sw a0, 0(a5)

# UART init completed!
# Pad to 0x60; we should be at 0x58 bytes into the file

addi x0, x0, 0x00
addi x0, x0, 0x00

# We need to read from UART0_TXDATA to confirm that UART is ready.
# UART0_TXDATA = 0x10013000, which is UART_REG_RXCTRL (already in register a5) minus 0xC (0x1001300C - 0xC = 0x10013000)
# So we can subtract 0xC from a5 and save an instruction
addi a5, a5, -0xC

# Read UART0_TXDATA into a0 until (a0 & 0x80000000) == 0
# Then put our data from a1, which we'll make 0x35 (an ASCII '5', for RISC-V!)
addi a1, x0, 0x35
lui a2, 0x80000
lw a0, 0(a5)
and a0, a0, a2
bne a0, x0, -8


sw a1, 0(a5)

addi x0, x0, 0x00
jal x0, 0x00
