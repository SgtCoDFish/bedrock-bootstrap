include ../common.mk

.PHONY: all
all: clean toolchain-test-asm.elf toolchain-test-asm.dump toolchain-test-hex.elf toolchain-test-hex.dump

.PHONY: clean
clean:
	rm -rf BUILD toolchain-test-asm.elf toolchain-test-asm.dump toolchain-test-hex.elf toolchain-test-hex.dump

.PHONY: gdb
gdb:
	$(RISCV_PREFIX)gdb -q -ex "set architecture riscv:rv32" -ex "target remote :1234" -ex "set confirm off"

# HEX example

BUILD/elfheader.bin: elfheader.hex
	sed "s/#.*$$//g" $< | xxd -r -p > $@

BUILD/toolchain-test-hex.bin: toolchain-test.hex
	sed "s/#.*$$//g" $< | xxd -r -p > $@

toolchain-test-hex.elf: BUILD/elfheader.bin BUILD/toolchain-test-hex.bin
	cat $^ > $@

toolchain-test-hex.dump: BUILD/toolchain-test-hex.bin
	$(RISCV_PREFIX)objdump -D -b binary -EL -m riscv:rv32 $< > $@

.PHONY: qemu-hex
qemu-hex: toolchain-test-hex.elf
	qemu-system-riscv32 -nographic -serial pty -gdb tcp::1234 -S -machine virt -bios none -kernel $<

# ASM Example

BUILD/toolchain-test.o: toolchain-test.asm | BUILD
	$(RISCV_PREFIX)as -march=rv32i $< -o $@

toolchain-test-asm.elf: BUILD/toolchain-test.o linker.ld
	$(RISCV_PREFIX)ld $< -Tlinker.ld -m elf32lriscv -o $@

toolchain-test-asm.dump: toolchain-test-asm.elf
	$(RISCV_PREFIX)objdump -d $< > $@

.PHONY: qemu-asm
qemu-asm: toolchain-test-asm.elf
	qemu-system-riscv32 -nographic -serial pty -gdb tcp::1234 -S -machine virt -bios none -kernel $<

BUILD:
	@mkdir -p $@
