include ../common.mk

.PHONY: all
all: BUILD/uart-rxxd.bin BUILD/uart-rxxd.elf

.PHONY: clean
clean:
	rm -rf BUILD

BUILD/uart-rxxd.bin: uart-rxxd.hex | BUILD
	@mkdir -p BUILD
	sed "s/#.*$$//g" $< | xxd -r -p > $@

BUILD/uart-rxxd.elf: elfheader.hex uart-rxxd.hex | BUILD
	cat $^ | sed "s/#.*$$//g" | xxd -r -p > $@

.PHONY: dump
dump: BUILD/uart-rxxd.bin
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $<

.PHONY: qemu
qemu: BUILD/uart-rxxd.elf
	qemu-system-riscv32 -nographic -serial pty -gdb tcp::1234 -S -machine sifive_e -kernel $<

BUILD:
	@mkdir -p $@
