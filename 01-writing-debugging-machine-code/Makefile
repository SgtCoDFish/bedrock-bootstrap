SHELL = /bin/bash
RISCV_PREFIX ?=

BUILD/bootloader1:
	@mkdir -p BUILD
	$(shell echo -n -e "\xb7\x02\x40\x20\x67\x80\x02\x00" > BUILD/bootloader1)
