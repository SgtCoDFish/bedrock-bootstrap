include ../common.mk

.PHONY: all
all: BUILD/uart-asm.elf

.PHONY: clean
clean:
	rm -rf BUILD

.PHONY: dump
dump: BUILD/uart-asm.elf
	$(RISCV_PREFIX)readelf -a $<

BUILD/start.o: uart.asm | BUILD
	$(RISCV_PREFIX)as -march=rv32i $< -o $@

BUILD/uart-asm.elf: BUILD/start.o | BUILD
	$(RISCV_PREFIX)ld $^ -Tlinker.ld --architecture=rv32i -m elf32lriscv -o $@

.PHONY: qemu-asm
qemu-asm: BUILD/uart-asm.elf
	qemu-system-riscv32 -nographic -serial pty -gdb tcp::1234 -S -machine virt -bios none -kernel $<

BUILD:
	@mkdir -p $@
