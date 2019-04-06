#include <stdint.h>
#include <stddef.h>
 
/* GPIO */
#define GPIO_CTRL_ADDR 0x10012000UL
#define GPIO_IOF_EN    0x38
#define GPIO_IOF_SEL   0x3C
#define IOF0_UART0_MASK         0x00030000UL
 
/* UART */
#define UART0_CTRL_ADDR 0x10013000UL
#define UART_REG_TXFIFO         0x00
#define UART_REG_RXFIFO         0x04
#define UART_REG_TXCTRL         0x08
#define UART_REG_RXCTRL         0x0c
#define UART_REG_IE             0x10
#define UART_REG_IP             0x14
#define UART_REG_DIV            0x18
#define UART_TXEN               0x1
 
/* PRCI */
#define PRCI_CTRL_ADDR 0x10008000UL
#define PRCI_HFROSCCFG (0x0000)
#define PRCI_PLLCFG (0x0008)
#define ROSC_EN(x) (((x) & 0x1) << 30)
#define PLL_REFSEL(x) (((x) & 0x1) << 17)
#define PLL_BYPASS(x) (((x) & 0x1) << 18)
#define PLL_SEL(x) (((x) & 0x1) << 16)
 
 
/* This function will read a 32-bit value from an MMIO register */
static inline uint32_t
mmio_read_u32(unsigned long reg, unsigned int offset)
{
	return (*(volatile uint32_t *) ((reg) + (offset)));
}
 
/* This function will write a byte to an MMIO register */
static inline void
mmio_write_u8(unsigned long reg, unsigned int offset, uint8_t val)
{
	(*(volatile uint32_t *) ((reg) + (offset))) = val;
}
 
/*This function will write a 32-bit value to an MMIO register */
static inline void
mmio_write_u32(unsigned long reg, unsigned int offset, uint32_t val)
{
	(*(volatile uint32_t *) ((reg) + (offset))) = val;
}
 
/* Initialize the UART */
static void
uart_init()
{
	/* These two writes enable the UART via the GPIO */
	mmio_write_u32(GPIO_CTRL_ADDR,
		GPIO_IOF_SEL,
		mmio_read_u32(GPIO_CTRL_ADDR, GPIO_IOF_SEL)
			& ~IOF0_UART0_MASK);
 
	mmio_write_u32(GPIO_CTRL_ADDR,
		GPIO_IOF_EN,
		mmio_read_u32(GPIO_CTRL_ADDR, GPIO_IOF_EN)
			| IOF0_UART0_MASK);
 
	/*
	 * Assuming a 16Mhz Crystal (which is Y1 on the HiFive1), the divisor
	 * for a 115200 baud rate is 138
	 */
	mmio_write_u32(UART0_CTRL_ADDR, UART_REG_DIV, 138);
	mmio_write_u32(UART0_CTRL_ADDR,
		UART_REG_TXCTRL,
		mmio_read_u32(UART0_CTRL_ADDR, UART_REG_TXCTRL)
			| UART_TXEN);
 
	/* busy loop until the line is asserted... */
	volatile int i = 0;
	while(i++ < 1000000);
}
 
/* Transmit a single byte over the UART */
static void
__uart_write(uint8_t byte)
{
	/* wait for the UART to become ready */
	while (mmio_read_u32(UART0_CTRL_ADDR, UART_REG_TXFIFO) & 0x80000000)
		;
 
	/* write to the UART transmit FIFO */
	mmio_write_u8(UART0_CTRL_ADDR, UART_REG_TXFIFO, byte);
}
 
 
/* Transmit a buffer of length "len" over the UART */
static void
uart_write(uint8_t *buf, size_t len)
{
	int i;
	for (i = 0; i < len; i ++) {
		__uart_write(buf[i]);
		/* If an LF was written, also write a CR */
		if (buf[i] == '\n') {
			__uart_write('\r');
		}
	}
}
 
/* People, the simplest ever strlen function */
static size_t
strlen(char *str)
{
    int len = 0;
    int i;
 
    for (i = 0; str[i] != 0; i ++)
        len ++;
 
    return len;
}
 
/* Write a null-terminated string to the UART, transmitting it */
static void
uart_write_string(uint8_t *buf) 
{
	uart_write(buf, strlen((char *) buf));
}
 
/* Initialize the clock source for the UART, in this case the 16MHz crystal */
static void
prci_init(void)
{
	/* Make sure the HFROSC is on */
	mmio_write_u32(PRCI_CTRL_ADDR, PRCI_HFROSCCFG,
			mmio_read_u32(PRCI_CTRL_ADDR, PRCI_HFROSCCFG)
			 | ROSC_EN(1));
 
	/* Run off 16 MHz Crystal for accuracy */
	mmio_write_u32(PRCI_CTRL_ADDR, PRCI_PLLCFG,
			mmio_read_u32(PRCI_CTRL_ADDR, PRCI_PLLCFG)
			 | (PLL_REFSEL(1) | PLL_BYPASS(1)));
	mmio_write_u32(PRCI_CTRL_ADDR, PRCI_PLLCFG,
			mmio_read_u32(PRCI_CTRL_ADDR, PRCI_PLLCFG)
			 | (PLL_SEL(1)));
 
	/* Turn off HFROSC to save power */
	mmio_write_u32(PRCI_CTRL_ADDR, PRCI_HFROSCCFG,
			mmio_read_u32(PRCI_CTRL_ADDR, PRCI_HFROSCCFG)
			 & ~(ROSC_EN(1)));
}
 
/* The entry point */
void
main(void)
{
	prci_init();
	uart_init();
 
	uart_write_string("Hello, world!\nThis is myOS on the HiFive-1 Board!\n");
 
	/* For now, just halt */
	for (;;);
}
 
/* The _actual_ entry point, this is then fixated to 0x20400000 via the linker script */
 __attribute__((section(".init")))
void
_start(void)
{
	main();
}
