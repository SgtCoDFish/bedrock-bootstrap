    .section .text
    .globl _start

# --------------------------------------------------
# Constants
# --------------------------------------------------
.equ UART_BASE, 0x10000000

.equ UART_RBR, 0x00
.equ UART_THR, 0x00
.equ UART_DLL, 0x00
.equ UART_IER, 0x01
.equ UART_DLM, 0x01
.equ UART_FCR, 0x02
.equ UART_LCR, 0x03
.equ UART_LSR, 0x05

.equ LCR_DLAB, 0x80
.equ LSR_THRE, 0x20
.equ LSR_DR,   0x01

# --------------------------------------------------
# Entry
# --------------------------------------------------
_start:
    call uart_init

    # Send 'H'
    li a0, 'H'
    call uart_putc

    # Send 'i'
    li a0, 'i'
    call uart_putc

    # Newline
    li a0, '\n'
    call uart_putc

	# Write to test device to signal QEMU to exit	
	li t0, 0x100000
	li t1, 0x5555   # pass
	sw t1, 0(t0)

    # Echo loop
echo_loop:
    call uart_getc     # a0 = received char
    call uart_putc     # echo it back
    j echo_loop

# --------------------------------------------------
# uart_init
# --------------------------------------------------
uart_init:
    li t0, UART_BASE

    # Disable interrupts
    sb zero, UART_IER(t0)

    # Enable DLAB
    li t1, LCR_DLAB
    sb t1, UART_LCR(t0)

    # Set baud rate divisor = 1 (115200 @ 1.8432 MHz)
    li t1, 1
    sb t1, UART_DLL(t0)
    sb zero, UART_DLM(t0)

    # 8 data bits, no parity, 1 stop bit, clear DLAB
    li t1, 0x03
    sb t1, UART_LCR(t0)

    # Enable FIFO, clear RX/TX FIFO
    li t1, 0x07
    sb t1, UART_FCR(t0)

    ret

# --------------------------------------------------
# uart_putc
#   a0 = character to send
# --------------------------------------------------
uart_putc:
    li t0, UART_BASE
wait_tx:
    lb t1, UART_LSR(t0)
    andi t1, t1, LSR_THRE
    beqz t1, wait_tx

    sb a0, UART_THR(t0)
    ret

# --------------------------------------------------
# uart_getc
#   returns character in a0
# --------------------------------------------------
uart_getc:
    li t0, UART_BASE
wait_rx:
    lb t1, UART_LSR(t0)
    andi t1, t1, LSR_DR
    beqz t1, wait_rx

    lb a0, UART_RBR(t0)
    ret
