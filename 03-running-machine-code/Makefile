SHELL = /bin/bash
RISCV_PREFIX ?=

.PHONY: all
all: clean ret1234.bin ret1234.elf ret1234.dump

.PHONY: clean
clean:
	rm -rf BUILD ret1234.elf ret1234.bin ret1234.dump

BUILD/start.o: start.s
	@mkdir -p BUILD
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-as -march=rv32ima $< -o $@

ret1234.elf: BUILD/start.o
	@mkdir -p BUILD
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-ld $^ -Tlinker.ld -o $@

ret1234.bin: ret1234.elf
	@mkdir -p BUILD
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-objcopy $< -O binary $@

ret1234.dump: ret1234.elf
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-objdump -d $< > $@
