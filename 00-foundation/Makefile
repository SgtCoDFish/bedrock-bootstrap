SHELL = /bin/bash

.PHONY: all
all: clean nothing.bin nothing.elf nothing.dump

.PHONY: clean
clean:
	rm -rf BUILD nothing.elf nothing.bin nothing.dump

BUILD/start.o: start.s
	@mkdir -p BUILD
	$(RISCV_PREFIX)as -march=rv32i $< -o $@

BUILD/nothing.o: nothing.c
	@mkdir -p BUILD
	$(RISCV_PREFIX)gcc -fPIC -O2 -g0 -static -fvisibility=hidden -nostdlib -nostartfiles -march=rv32i -mabi=ilp32 -c $< -o $@

nothing.elf: BUILD/nothing.o BUILD/start.o
	@mkdir -p BUILD
	$(RISCV_PREFIX)ld $^ -Tlinker.ld -m elf32lriscv -o $@

nothing.bin: nothing.elf
	@mkdir -p BUILD
	$(RISCV_PREFIX)objcopy $< -O binary $@

nothing.dump: nothing.elf
	$(RISCV_PREFIX)objdump -d $< > $@

.PHONY: qemu
qemu:
	qemu-system-riscv32 -machine sifive_e -nographic -gdb tcp::1234 -S -kernel nothing.elf
