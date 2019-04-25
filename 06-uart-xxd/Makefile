RISCV_PREFIX ?= riscv32-unknown-elf-

.PHONY: all
all: BUILD/uart-rxxd.bin BUILD/uart-rxxd.elf

.PHONY: clean
clean:
	rm -rf BUILD

BUILD/uart-rxxd.bin: uart-rxxd.hex
	@mkdir -p BUILD
	sed "s/#.*$$//g" $< | xxd -r -p > $@

.PHONY: dump
dump: BUILD/uart-rxxd.bin
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $<

BUILD/uart-rxxd.elf: elfheader.hex uart-rxxd.hex
	@mkdir -p BUILD
	cat $^ | sed "s/#.*$$//g" | xxd -r -p > $@
