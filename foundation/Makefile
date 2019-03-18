RISCV_PREFIX ?=

.PHONY: all
all: clean BUILD/nothing.bin BUILD/nothing.elf BUILD/nothing.dump

.PHONY: clean
clean:
	rm -rf BUILD

BUILD/start.o: start.s
	@mkdir -p BUILD
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-as -march=rv32imac $< -o $@

BUILD/nothing.o: nothing.c
	@mkdir -p BUILD
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-gcc -fPIC -O2 -g0 -static -fvisibility=hidden -nostdlib -nostartfiles -march=rv32imac -mabi=ilp32 $< -o $@

BUILD/nothing.elf: BUILD/nothing.o BUILD/start.o
	@mkdir -p BUILD
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-ld $^ -Tlinker.ld -o $@

BUILD/nothing.bin: BUILD/nothing.elf
	@mkdir -p BUILD
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-objcopy $< -O binary $@

BUILD/nothing.dump: BUILD/nothing.elf
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-objdump -d $< > $@
