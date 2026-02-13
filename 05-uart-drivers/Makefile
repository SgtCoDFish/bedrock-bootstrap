include ../common.mk

.PHONY: all
all: BUILD/sifive-uart.bin BUILD/virt-uart.bin BUILD/sifive-example.bin BUILD/virt-example.bin BUILD/sifive-example.elf BUILD/virt-example.elf

.PHONY: clean
clean:
	rm -rf BUILD

BUILD/sifive-uart.bin: sifive-uart.hex | BUILD
	sed "s/#.*$$//g" $< | xxd -r -p > $@

BUILD/virt-uart.bin: virt-uart.hex | BUILD
	sed "s/#.*$$//g" $< | xxd -r -p > $@

BUILD/sifive-example.bin: sifive-uart.hex uart.hex | BUILD
	cat $^ | sed "s/#.*$$//g" | xxd -r -p > $@

BUILD/virt-example.bin: virt-uart.hex uart.hex | BUILD
	cat $^ | sed "s/#.*$$//g" | xxd -r -p > $@

BUILD/sifive-example.elf: elfheader.sifive.hex sifive-uart.hex uart.hex | BUILD
	cat $^ | sed "s/#.*$$//g" | xxd -r -p > $@

BUILD/virt-example.elf: elfheader.virt.hex virt-uart.hex uart.hex | BUILD
	cat $^ | sed "s/#.*$$//g" | xxd -r -p > $@

.PHONY: dump-sifive-driver
dump-sifive-driver: BUILD/sifive-uart.bin
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $<

.PHONY: dump-virt-driver
dump-virt-driver: BUILD/virt-uart.bin
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $<

.PHONY: dump-sifive-example
dump-sifive-example: BUILD/sifive-example.bin
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $<

.PHONY: dump-virt-example
dump-virt-example: BUILD/virt-example.bin
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $<

.PHONY: qemu-sifive
qemu-sifive: BUILD/sifive-example.elf
	qemu-system-riscv32 -nographic -serial pty -gdb tcp::1234 -S -machine sifive_e -kernel BUILD/uart.elf

.PHONY: qemu-virt
qemu-virt: BUILD/virt-example.elf
	qemu-system-riscv32 -nographic -serial pty -gdb tcp::1234 -S -machine virt -bios none -kernel $<
BUILD:
	@mkdir -p $@
