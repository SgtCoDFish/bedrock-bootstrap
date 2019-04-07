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

- [rv8.io](https://rv8.io/asm.html) is a good references / cheat sheet for RISC-V assembly programming
- [RISC-V Calling Convention](https://riscv.org/wp-content/uploads/2015/01/riscv-calling.pdf) details a calling convention for RISC-V programs. It's very useful to have such a calling convention to work to.