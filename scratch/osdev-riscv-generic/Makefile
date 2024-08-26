bindir=_bin

.PHONY: build
build: $(bindir)/kernel.elf

$(bindir)/kernel.o: kernel.c | $(bindir)
	riscv64-elf-gcc -Wall -Wextra -c -mcmodel=medany $< -o $@ -ffreestanding

$(bindir)/entry.o: entry.S | $(bindir)
	riscv64-elf-as -c $< -o $@

$(bindir)/kernel.elf: linker.ld $(bindir)/entry.o $(bindir)/kernel.o | $(bindir)
	riscv64-elf-ld -T linker.ld -nostdlib $(bindir)/kernel.o $(bindir)/entry.o -o $@

$(bindir):
	mkdir -p $@

.PHONY: qemu
qemu: $(bindir)/kernel.elf
	qemu-system-riscv64 -machine virt -bios none -kernel $(bindir)/kernel.elf -serial mon:stdio -nographic
