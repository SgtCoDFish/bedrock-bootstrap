include ../common.mk

.PHONY: all
all: BUILD/sifive-uart.bin BUILD/virt-uart.bin

.PHONY: clean
clean:
	rm -rf BUILD

BUILD/sifive-uart.bin: sifive-uart.hex | BUILD
	sed "s/#.*$$//g" $< | xxd -r -p > $@

BUILD/virt-uart.bin: virt-uart.hex | BUILD
	sed "s/#.*$$//g" $< | xxd -r -p > $@

.PHONY: dump-sifive
dump-sifive: BUILD/sifive-uart.bin
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $<

.PHONY: dump-virt
dump-virt: BUILD/virt-uart.bin
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $<

BUILD:
	@mkdir -p $@
