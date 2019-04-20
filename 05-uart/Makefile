RISCV_PREFIX ?= riscv32-unknown-elf-

BUILD/uart.bin: uart.hex
	@mkdir -p BUILD
	xxd -r -p uart.hex > $@

.PHONY: dump
dump: BUILD/uart.bin
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $<
