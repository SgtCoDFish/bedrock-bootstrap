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

BUILD/uart-functions.o: uart-functions.asm | BUILD
	$(RISCV_PREFIX)as -march=rv32i $< -o $@

BUILD:
	@mkdir -p $@
