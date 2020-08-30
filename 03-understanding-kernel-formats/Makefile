SHELL = /bin/bash
RISCV_PREFIX ?= riscv32-unknown-elf-

.PHONY: all
all: clean ret1234.bin ret1234.elf ret1234.dump

.PHONY: clean
clean:
	rm -rf BUILD ret1234.elf ret1234.bin ret1234.dump

BUILD/start.o: start.s
	@mkdir -p BUILD
	$(RISCV_PREFIX)as -march=rv32i $< -o $@

ret1234.elf: BUILD/start.o
	@mkdir -p BUILD
	$(RISCV_PREFIX)ld $^ -Tlinker.ld --architecture=rv32i -m elf32lriscv -o $@

ret1234.bin: ret1234.elf
	@mkdir -p BUILD
	$(RISCV_PREFIX)objcopy $< -O binary $@

ret1234.dump: ret1234.elf
	$(RISCV_PREFIX)objdump -d $< > $@
