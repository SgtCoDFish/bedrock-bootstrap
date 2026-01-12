include ../common.mk

.PHONY: all
all: BUILD/uart-functions.bin BUILD/uart-functions.elf

.PHONY: clean
clean:
	rm -rf BUILD

BUILD/uart-functions.bin: uart-functions.hex | BUILD
	sed "s/#.*$$//g" $< | xxd -r -p > $@

.PHONY: dump
dump: BUILD/uart-functions.bin
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $<

BUILD/uart-functions.elf: elfheader.hex uart-functions.hex | BUILD
	cat $^ | sed "s/#.*$$//g" | xxd -r -p > $@

.PHONY: qemu
qemu: BUILD/uart-functions.elf
	qemu-system-riscv32 -nographic -serial pty -gdb tcp::1234 -S -machine sifive_e -kernel $<

BUILD:
	@mkdir -p $@
