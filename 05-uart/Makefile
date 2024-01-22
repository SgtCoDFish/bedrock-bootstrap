RISCV_PREFIX ?= riscv32-unknown-elf-

.PHONY: all
all: BUILD/uart.bin BUILD/uart.elf

.PHONY: clean
clean:
	rm -rf BUILD

BUILD/uart.bin: uart.hex
	@mkdir -p BUILD
	sed "s/#.*$$//g" $< | xxd -r -p > $@

.PHONY: dump
dump: BUILD/uart.bin
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $<

BUILD/uart.elf: elfheader.hex uart.hex
	@mkdir -p BUILD
	cat elfheader.hex uart.hex | sed "s/#.*$$//g" | xxd -r -p > $@

.PHONY: qemu
qemu:
	qemu-system-riscv64 -nographic -serial pty -gdb tcp::1234 -S -machine sifive_e -kernel BUILD/uart.elf

.PHONY: qemu-virt
qemu-virt:
	# TODO: this won't work; it jumps to 0x80000000
	qemu-system-riscv64 -nographic -serial pty -gdb tcp::1234 -S -machine virt -bios none -kernel BUILD/uart.elf
