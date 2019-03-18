# Foundations
Before we can do _anything_ we'll need to be able to debug what we're doing and gather some required documentation.

The [HiFive1 getting started guide](https://sifive.cdn.prismic.io/sifive%2F9c57065b-6d28-465b-b67d-f416894123a9_hifive1-getting-started-v1.0.2.pdf) is a good place for some details on the HiFive1, but it mainly uses
a slightly clunky-feeling IDE based workflow. Ideally we want to do everything on the command line where it'll be fast
and easily scriptable. There's also a [more detailed datasheet](https://sifive.cdn.prismic.io/sifive%2F4d063bf8-3ae6-4db6-9843-ee9076ebadf7_fe310-g000.pdf).

[Freedom Metal](https://sifive.github.io/freedom-metal-docs/introduction.html#what-is-freedom-metal) (on [GitHub](https://github.com/sifive/freedom-metal/tree/master)) looks more useful, as does the [freedom-e-sdk](https://github.com/sifive/freedom-e-sdk). We see they provide
a [Board Support Package](https://github.com/sifive/freedom-e-sdk/tree/30c143eb5445f47edb351ba54c84ff8285dc27a9/bsp/sifive-hifive1) for the HiFive1, with details we can use. We also see example programs
such as [hello](https://github.com/sifive/example-hello/tree/d1397bec64187efb8b791fe1eb307aa3c760c694) which we can use as a sanity check. The makefile is a big beast designed to support multiple different boards, and we can definitely trim it down.

We're not actually going to be _using_ this library directly, but since it's open source we can learn from it. It's
also in C, which isn't any use to us in a bedrock bare metal world - but again, will point us in the right direction.
In any case, we'll be needing a [GNU compiler toolchain](https://github.com/riscv/riscv-gnu-toolchain) which is going to take a _long_ time to clone if you don't have fast internet.

While you're at it you'll want to clone the `freedom-e-sdk` and `riscv-qemu`.

```bash
git clone --recursive https://github.com/riscv/riscv-gnu-toolchain
git clone --recursive https://github.com/sifive/riscv-qemu
git clone --recursive https://github.com/sifive/freedom-e-sdk
```

Once that's done we'll need to build.

# TODO: Building "hello world" and running on qemu / hardware
