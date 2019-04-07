# RISC-V / HiFive / Bedrock Bootstrap Resources

Part of what makes RISC-V and the HiFive1 a good choice for this project is the wide
availability of resources which can be used to learn more about chips and the board.

## HiFive1

### References

- The [HiFive1 Getting Started Guide](https://sifive.cdn.prismic.io/sifive%2F9c57065b-6d28-465b-b67d-f416894123a9_hifive1-getting-started-v1.0.2.pdf) is a good read for getting acquainted with the spectrum of components contained on the board.
- There's a [more detailed datasheet](https://sifive.cdn.prismic.io/sifive%2F4d063bf8-3ae6-4db6-9843-ee9076ebadf7_fe310-g000.pdf) which will be a useful reference.

### Libraries

- [Freedom Metal](https://sifive.github.io/freedom-metal-docs/introduction.html#what-is-freedom-metal) (code [on GitHub](https://github.com/sifive/freedom-metal/tree/master)) is a bare-metal library for programming various SiFive products including the HiFive1. It's far too heavyweight for our use case, but can inform certain choices we might make.
- The [freedom-e-sdk](https://github.com/sifive/freedom-e-sdk) is the main source of code and config files we might need during development.
- Most imporantly from the SDK is the "board support package" (bsp) for the [HiFive1](https://github.com/sifive/freedom-e-sdk/tree/30c143eb5445f47edb351ba54c84ff8285dc27a9/bsp/sifive-hifive1) which contains a device map, a (slightly complicated) linker script and an OpenOCD config file which we'll use to upload code to the device.

We're not going to be _using_ the "Freedom Metal" library or much of the freedom-e-sdk directly, but since they're open source we can learn from them and use any relevant bits.

## RISC-V

- [RISC-V ISA Manual 2.2](https://github.com/riscv/riscv-isa-manual/blob/3f98f6087b75e52ec4f61681769b5f6931df2f06/release/riscv-spec-v2.2.pdf) is the canonical reference for RISC-V assembly. There's plenty of detailed information in there, but Chapter 20, the "RISC-V Assembly Programmer's Handbook" is well worth a read since it contains details about the calling convention we'll use along with a mapping of pseudoinstructions to base instructions.
- [rv8.io](https://rv8.io/asm.html) is a good reference / cheat sheet for RISC-V assembly programming

## Other Links

- dwelch67 always has good guides for bare metal programming; for example, check out his [uart01 sample](https://github.com/dwelch67/sifive_samples/tree/master/hifive1/uart01) for a bare-metal RISC-V UART example
- Running RISC-V on qemu bare metal [google groups thread](https://groups.google.com/a/groups.riscv.org/forum/#!topic/sw-dev/IET9LBFJohU). This thread was useful for figuring out how Qemu could be used productively with bare metal development.
- A very minimal bare metal example, similar to dwelch but based on riscv64 and SPIKE: [schoeberl](https://github.com/schoeberl/cae-examples).
- Embedded bare-metal Rust on the Raspberry Pi 3: [rust-raspi3-OS-tutorials](https://github.com/rust-embedded/rust-raspi3-OS-tutorials). Quite different in scope, but a brilliantly written resource in general.
- [OpenOCD SPI Flash Driver](https://github.com/riscv/riscv-openocd/blob/riscv/src/flash/nor/fespi.c) - a neat insight into how the driver works which lets us upload our code.
- [OSDev HiFive1 Bare Bones](https://wiki.osdev.org/HiFive-1_Bare_Bones#The_Bare_Bones) - a brilliant C tutorial which outlines working on bare metal on the HiFive1. The OSDev Wiki is a fantastic resource itself.
- [bcompiler (mirror)](https://github.com/smtlaissezfaire/bcompiler): the page which originally inspired this project was "Bootstrapping a simple compiler from nothing" by Edmund Grimley Evans, which does a similar project using x86. The original site has disappeared but this GitHub repository has a mirror of the code.