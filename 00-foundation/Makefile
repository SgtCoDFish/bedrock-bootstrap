SHELL = /bin/bash
RISCV_PREFIX ?=

.PHONY: all
all: clean nothing.bin nothing.elf nothing.dump

.PHONY: clean
clean:
	rm -rf BUILD nothing.elf nothing.bin nothing.dump

BUILD/start.o: start.s
	@mkdir -p BUILD
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-as -march=rv32ima $< -o $@

BUILD/nothing.o: nothing.c
	@mkdir -p BUILD
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-gcc -fPIC -O2 -g0 -static -fvisibility=hidden -nostdlib -nostartfiles -march=rv32ima -mabi=ilp32 $< -o $@

BUILD/bootloader1: 
	@mkdir -p BUILD
	$(shell echo -n -e "\xb7\x02\x40\x20\x67\x80\x02\x00" > BUILD/bootloader1)

nothing.elf: BUILD/nothing.o BUILD/start.o
	@mkdir -p BUILD
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-ld $^ -Tlinker.ld -o $@

nothing.bin: nothing.elf
	@mkdir -p BUILD
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-objcopy $< -O binary $@

nothing.dump: nothing.elf
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-objdump -d $< > $@
