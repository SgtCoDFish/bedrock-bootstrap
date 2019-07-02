#include <stdint.h>

#define GPIOBASE        0x10012000
#define GPIO_IOF_EN     (GPIOBASE+0x38)
#define GPIO_IOF_SEL    (GPIOBASE+0x3C)

#define UART0BASE       0x10013000
#define UART0_TXDATA    (UART0BASE+0x00)
#define UART0_RXDATA    (UART0BASE+0x04)
#define UART0_TXCTRL    (UART0BASE+0x08)
#define UART0_RXCTRL    (UART0BASE+0x0C)
#define UART0_IE        (UART0BASE+0x10)
#define UART0_IP        (UART0BASE+0x14)
#define UART0_DIV       (UART0BASE+0x18)

inline static uint32_t __attribute__((always_inline)) GET32(uint32_t addr) {
	return *((volatile uint32_t *volatile) addr);
}

inline static void __attribute__((always_inline)) PUT32(uint32_t addr, uint32_t val) {
	*((volatile uint32_t *) addr) = val;
}

inline static void __attribute__((always_inline)) uart_init (void) {
	unsigned int ra = GET32(GPIO_IOF_SEL);

	ra &= ~(1<<16);
	ra &= ~(1<<17);

	PUT32(GPIO_IOF_SEL,ra);

	ra = GET32(GPIO_IOF_EN);
	ra |= 1<<16;
	ra |= 1<<17;

	PUT32(GPIO_IOF_EN, ra);

	PUT32(UART0_DIV, 138);
	PUT32(UART0_TXCTRL, 0x00000003);
	PUT32(UART0_RXCTRL, 0x00000001);
}

void __attribute__((noreturn)) xmain(void) {
	//uart_init();
	uint32_t x05;
	uint32_t x06;
	uint32_t x07;
	uint32_t x08;
	uint32_t x09;
	uint32_t x10;
	uint32_t x11;
	uint32_t x12;
	uint32_t x14;
	uint32_t x15;
	uint32_t x16;
	uint32_t x17;
	uint32_t x18;

SETUP:
	x05 = 0x00;
	x06 = 0x00;
	x07 = 0x04;
	x08 = 0x00;
	x09 = 0x00;
	x10 = 0x00;
	x11 = 0x00;
	x12 = 0x80000000;
	x14 = 0x80001000;
	x15 = UART0_TXDATA;
	x16 = UART0_RXDATA;
	x17 = 0xA;
	x18 = 0x20;

READ_UART:
	//x10 = GET32(x16);
	x10 = *((volatile uint32_t *) x16);
	x11 = x10 & x12;
	if (x11 != 0) { goto READ_UART; }

	x10 &= 0xFF;

	if (x08 == 0) { goto NOT_COMMENT_MODE; }

	if (x10 != x17) { goto READ_UART; }

	x08 = 0;
	goto READ_UART;

NOT_COMMENT_MODE:
	x09 = 0x23;
	if (x10 != x09) { goto NO_START_COMMENT_MODE; }
	x08 = 1;
	goto READ_UART;

NO_START_COMMENT_MODE:
	x09 = 0x30;
	if (x10 > x09) { goto NO_SKIP_LOW_CHAR; }
	goto READ_UART;

NO_SKIP_LOW_CHAR:
	x08 = 0x3A;
	if (x10 >= x08) { goto PARSE_CHAR; }
	x10 += -0x30;
	goto ADDRAW;

PARSE_CHAR:
	x10 |= 0x20;
	x08 = 0x61;
	if (x10 < x08) { goto PANIC; }

	x08 = 0x67;
	if (x10 > x08) { goto PANIC; }

	x10 += -0x57;
	goto ADDRAW;

ADDRAW:
	x10 <<= x07;
	x09 |= x10;
	if (x07 == 0) { goto ADDBYTE; }
	x07 = 0x0;
	goto READ_UART;

ADDBYTE:
	x07 = 0x4;
	x09 <<= x06;
	x05 |= x09;
	x06 += 0x8;
	if (x06 == x18) { goto WRITEWORD; }
	goto READ_UART;

WRITEWORD:
	PUT32(x14, x05);
	x14 += 4;
	goto READ_UART;

PANIC:
	for(;;){}
}
