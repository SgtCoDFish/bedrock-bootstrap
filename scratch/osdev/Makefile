main.img: main.c linker.ld
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-gcc main.c -g -march=rv32imac -mabi=ilp32 -mcmodel=medany -o $@ -T linker.ld -nostartfiles
