RISCV_PREFIX ?=

.PHONY: all
all: clean BUILD/shimtest.bin BUILD/shimtest.elf BUILD/shimtest.dump

.PHONY: clean
clean:
	rm -rf BUILD

BUILD/start.o: start.s
	@mkdir -p BUILD
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-as -march=rv32imac $< -o $@

BUILD/shim.o: shim.s
	@mkdir -p BUILD
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-as -D -march=rv32imac $< -o $@

BUILD/shimtest.elf: BUILD/start.o BUILD/shim.o
	@mkdir -p BUILD
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-ld $^ -Tlinker.ld -o $@

BUILD/shimtest.bin: BUILD/shimtest.elf
	@mkdir -p BUILD
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-objcopy $< -O binary $@

BUILD/shimtest.dump: BUILD/shimtest.elf
	$(RISCV_PREFIX)/bin/riscv32-unknown-elf-objdump -d $< > $@
