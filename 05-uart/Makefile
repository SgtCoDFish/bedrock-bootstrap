RISCV_PREFIX ?= riscv32-unknown-elf-

.PHONY: all
all: BUILD/uart.bin BUILD/uart.elf BUILD/uart-virt.bin BUILD/uart-virt.elf

.PHONY: clean
clean:
	rm -rf BUILD

BUILD/uart.bin: uart.hex | BUILD
	sed "s/#.*$$//g" $< | xxd -r -p > $@

BUILD/uart-virt.bin: uart-virt.hex | BUILD
	sed "s/#.*$$//g" $< | xxd -r -p > $@

.PHONY: dump
dump: BUILD/uart.bin
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $<

.PHONY: dump-virt
dump-virt: BUILD/uart-virt.bin
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $<

BUILD/uart.elf: elfheader.hex uart.hex | BUILD
	@mkdir -p BUILD
	cat elfheader.hex uart.hex | sed "s/#.*$$//g" | xxd -r -p > $@

BUILD/uart-virt.elf: elfheader-virt.hex uart-virt.hex | BUILD
	@mkdir -p BUILD
	cat elfheader-virt.hex uart-virt.hex | sed "s/#.*$$//g" | xxd -r -p > $@

.PHONY: qemu
qemu: BUILD/uart.elf
	qemu-system-riscv64 -nographic -serial pty -gdb tcp::1234 -S -machine sifive_e -kernel BUILD/uart.elf

.PHONY: qemu-virt
qemu-virt: BUILD/uart-virt.elf
	qemu-system-riscv64 -nographic -serial mon:stdio -gdb tcp::1234 -S -machine virt -bios none -kernel BUILD/uart-virt.elf

BUILD:
	@mkdir -p $@
