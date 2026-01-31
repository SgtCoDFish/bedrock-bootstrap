include ../common.mk

BUILD/example: | BUILD
	@mkdir -p BUILD
	$(shell echo -n -e "\x37\x51\x34\x12" > $@)

.PHONY: dump
dump: BUILD/example
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $<

BUILD:
	@mkdir -p $@
