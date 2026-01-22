    .section .text
    .globl _start

/* --------------------------------------------------------------------
 * Constants
 * ------------------------------------------------------------------*/
UART_BASE   = 0x10013000
UART_TXDATA = 0x00

SRC_START  = 0x20400000
SRC_END    = 0x20400100

/* --------------------------------------------------------------------
 * Program entry
 * ------------------------------------------------------------------*/
_start:
/* UART setup, taken from 05-uart */
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

    /* Load source start address */
    lui     t0, %hi(SRC_START)
    addi    t0, t0, %lo(SRC_START)   /* t0 = current pointer */

    /* Load source end address */
    lui     t1, %hi(SRC_END)
    addi    t1, t1, %lo(SRC_END)     /* t1 = end pointer */

    /* Load UART base address */
    lui     t2, %hi(UART_BASE)
    addi    t2, t2, %lo(UART_BASE)   /* t2 = UART base */

send_loop:
    /* Check if we've reached the end */
    beq     t0, t1, done

    /* Load one byte from memory */
    lbu     t3, 0(t0)

uart_wait:
    /* Read UART TXDATA */
    lw      t4, UART_TXDATA(t2)

    /* Check txfull bit (bit 31) */
    bltz    t4, uart_wait

    /* Write byte to UART TXDATA */
    sw      t3, UART_TXDATA(t2)

    /* Increment source pointer */
    addi    t0, t0, 1

    /* Loop */
    j       send_loop

done:
    /* Spin forever */
    j       done
