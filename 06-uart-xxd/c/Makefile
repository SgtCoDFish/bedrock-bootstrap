SHELL = /bin/bash
RISCV_PREFIX ?= riscv32-unknown-elf-

.PHONY: all
all: clean BUILD/uart-rxxd.bin BUILD/uart-rxxd.elf BUILD/uart-rxxd.dump

.PHONY: clean
clean:
	rm -rf BUILD

BUILD/start.o: start.s
	@mkdir -p BUILD
	$(RISCV_PREFIX)as -march=rv32i $< -o $@

BUILD/uart-rxxd.o: uart-rxxd.c
	@mkdir -p BUILD
	$(RISCV_PREFIX)gcc -fPIC -O0 -g0 -static -fvisibility=hidden -nostdlib -nostartfiles -march=rv32i -mabi=ilp32 $< -o $@

BUILD/uart-rxxd.elf: BUILD/start.o BUILD/uart-rxxd.o
	@mkdir -p BUILD
	$(RISCV_PREFIX)ld $^ -Tlinker.ld -o $@

BUILD/uart-rxxd.bin: BUILD/uart-rxxd.elf
	@mkdir -p BUILD
	$(RISCV_PREFIX)objcopy $< -O binary $@

BUILD/uart-rxxd.dump: BUILD/uart-rxxd.elf
	@mkdir -p BUILD
	$(RISCV_PREFIX)objdump -d $< > $@
