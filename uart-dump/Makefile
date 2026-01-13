include ../common.mk

.PHONY: all
all: clean BUILD/dump-uart.elf BUILD/qemu-dump.bin

.PHONY: clean
clean:
	rm -rf BUILD

BUILD/start.o: start.s | BUILD
	$(RISCV_PREFIX)as -march=rv32i $< -o $@

BUILD/dump-uart.elf: BUILD/start.o | BUILD
	$(RISCV_PREFIX)ld $^ -Tlinker.ld --architecture=rv32i -m elf32lriscv -o $@

.PHONY: dump-asm
dump-asm: BUILD/dump-uart.elf
	$(RISCV_PREFIX)objdump -d $<

.PHONY: dump-dumped
# Dumps the values dumped over UART
dump-dumped: BUILD/qemu-dump.bin
	riscv64-elf-objdump -D -b binary -EL -m riscv:rv32 $<

.PHONY: qemu
qemu:
	qemu-system-riscv32 -nographic -serial pty -gdb tcp::1234 -S -machine sifive_e -kernel $<

# Run dump-uart in qemu to dump the program itself.
# Note this will hang and need SIGINT to kill qemu.
.PRECIOUS: BUILD/qemu-dump.bin
BUILD/qemu-dump.bin: BUILD/dump-uart.elf | BUILD
	qemu-system-riscv32 -nographic -monitor none -serial stdio -machine sifive_e -kernel $< > $@

BUILD:
	@mkdir -p $@
