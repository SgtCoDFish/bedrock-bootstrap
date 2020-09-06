#include "gd32vf103.h"
#include "lcd/lcd.h"
#include "fatfs/tf_card.h"
#include <stdio.h>

int main(void)
{
	/* enable GPIO clock */
	rcu_periph_clock_enable(RCU_GPIOA);
	rcu_periph_clock_enable(RCU_GPIOC);

	/* enable USART clock */
	rcu_periph_clock_enable(RCU_USART0);

	// enable red LED control
	gpio_init(GPIOC, GPIO_MODE_OUT_PP, GPIO_OSPEED_50MHZ, GPIO_PIN_13);

	// enable blue and green LED control
	gpio_init(GPIOA, GPIO_MODE_OUT_PP, GPIO_OSPEED_50MHZ, GPIO_PIN_1|GPIO_PIN_2);

	/* connect port to USARTx_Tx */
	gpio_init(GPIOA, GPIO_MODE_AF_PP, GPIO_OSPEED_50MHZ, GPIO_PIN_9);

	/* connect port to USARTx_Rx */
	gpio_init(GPIOA, GPIO_MODE_IN_FLOATING, GPIO_OSPEED_50MHZ, GPIO_PIN_10);

	/* USART configure */
	usart_deinit(USART0);
	usart_baudrate_set(USART0, 115200U);
	usart_word_length_set(USART0, USART_WL_8BIT);
	usart_stop_bit_set(USART0, USART_STB_1BIT);
	usart_parity_config(USART0, USART_PM_NONE);
	usart_hardware_flow_rts_config(USART0, USART_RTS_DISABLE);
	usart_hardware_flow_cts_config(USART0, USART_CTS_DISABLE);
	usart_receive_config(USART0, USART_RECEIVE_ENABLE);
	usart_transmit_config(USART0, USART_TRANSMIT_ENABLE);
	usart_enable(USART0);

	// 1 means "off" for LEDs
	LEDR(1);
	LEDG(1);
	LEDB(1);

	while (1) {
		uint32_t flag = 0u;
		while (flag = usart_flag_get(USART0, USART_FLAG_RBNE)) {
			uint16_t val = usart_data_receive(USART0);
			printf("got: %08x\n", val);
			if (val == 0x31) {
				LEDR_TOG;
			} else if (val == 0x32) {
				LEDG_TOG;
			} else if (val == 0x33) {
				LEDB_TOG;
			}
		}

		delay_1ms(200);
	}

}

/* retarget the C library printf function to the USART */
int _put_char(int ch)
{
	usart_data_transmit(USART0, (uint8_t) ch );
	while ( usart_flag_get(USART0, USART_FLAG_TBE)== RESET){
	}

	return ch;
}
