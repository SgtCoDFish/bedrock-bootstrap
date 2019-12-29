RISCV_PREFIX ?= riscv32-unknown-elf-

.PHONY: all
all: BUILD/uart-functions.bin BUILD/uart-functions.elf

.PHONY: clean
clean:
	rm -rf BUILD

BUILD/uart-functions.bin: uart-functions.hex
	@mkdir -p BUILD
	sed "s/#.*$$//g" $< | xxd -r -p > $@

.PHONY: dump
dump: BUILD/uart-functions.bin
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $<

BUILD/uart-functions.elf: elfheader.hex uart-functions.hex
	@mkdir -p BUILD
	cat $^ | sed "s/#.*$$//g" | xxd -r -p > $@
