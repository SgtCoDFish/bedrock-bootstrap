bindir := _bin

$(bindir)/bootloader.bin: bootloader.hex | $(bindir)
	sed "s/#.*$$//g" $< | xxd -r -p > $@

.PHONY: dump
dump: $(bindir)/bootloader.bin
	riscv64-elf-objdump -D -b binary -EL -m riscv:rv64 $<

$(bindir):
	mkdir -p $@
