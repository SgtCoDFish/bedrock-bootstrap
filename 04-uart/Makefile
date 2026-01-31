include ../common.mk

.PHONY: all
all: BUILD/uart.bin BUILD/uart.elf

.PHONY: clean
clean:
	rm -rf BUILD

BUILD/uart.bin: uart.hex | BUILD
	sed "s/#.*$$//g" $< | xxd -r -p > $@

.PHONY: dump
dump: BUILD/uart.bin
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $<

BUILD/uart.elf: elfheader.hex uart.hex | BUILD
	cat elfheader.hex uart.hex | sed "s/#.*$$//g" | xxd -r -p > $@

.PHONY: qemu
qemu: BUILD/uart.elf
	qemu-system-riscv32 -nographic -serial pty -gdb tcp::1234 -S -machine sifive_e -kernel BUILD/uart.elf

BUILD:
	@mkdir -p $@
