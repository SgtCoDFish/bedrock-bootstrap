# Legacy Tooling

This file attempts to preserve information from a while ago to avoid losing it.

The instructions in TOOLS.md should be simpler and appropriate for almost all modern systems.

## Building QEMU

If your system has QEMU >= 3.1 RISC-V support is likely to be included and you can probably install a RISC-V-supporting QEMU from your system package manager. If you run into problems using the version from your package manage, it's very easy to build QEMU from source and use that.

Clone it yourself and build:

```bash
git clone --recursive --depth 1 https://github.com/qemu/QEMU riscv-qemu
cd riscv-qemu
mkdir build && cd build
../configure --target-list=riscv32-softmmu
make -j4
```

## The GNU Toolchain

You can use a pre-built toolchain or build your own.

If you're building from source, you can see in the freedom-e-sdk HiFive1 [BSP](https://github.com/sifive/freedom-e-sdk/blob/30c143eb5445f47edb351ba54c84ff8285dc27a9/bsp/sifive-hifive1/settings.mk) that we need to target a different arch and ABI, since the toolchain defaults to 64 bit.

You can choose where you want to run the tools from, but the build commands used in this project assume you've set `$RISCV_PREFIX` to point at the correct executables on the path. For example, if you install to `/opt/riscv` then you should set `$RISCV_PREFIX` to `/opt/riscv/bin/riscv32-unknown-elf-`.

Also, note that if you're installing somewhere where your user doesn't have write permissions (e.g. `/opt/riscv`) you'll probably need to build as root (`sudo make -j4`).

Before building, you're likely to need to install some additional requirements. See [the repo](https://github.com/riscv/riscv-gnu-toolchain) for requirements on various platforms including popular Linus distros and macOS.

```bash
git clone --recursive https://github.com/riscv/riscv-gnu-toolchain
cd riscv-gnu-toolchain
mkdir build && cd build
../configure --with-arch=rv32ima --with-abi=ilp32 --with-cmodel=medlow --prefix=SOME_PREFIX  # change SOME_PREFIX to whatever you like
make -j4  # might need to be done as root depending on where you're installing
```

Note that we're using the arch `rv32ima` which means the base RISC-V 32-bit instruction set (`i`) plus extensions for multiplication and atomic operations. You'll often see a `c` used in other toolchains; its omission is a conscious decision to simplify our efforts to write bare-metal machine code.

In any case, we can (and will) explicitly specify the arch later when building which will allow us to avoid using extensions even if the toolchain is built with support for various extensions.

The executables will be in the `bin` directory.

## OpenOCD

OpenOCD is used to transfer binaries to the HiFive1 system, and for debugging.

At the time that bedrock-boostrap was started RISC-V wasn't supported by any versions of OpenOCD which were installable from package managers like apt or homebrew. There is a [patched version](https://github.com/riscv/riscv-openocd) which is easy to build, however:

```bash
git clone https://github.com/riscv/riscv-openocd
cd riscv-openocd
./bootstrap
./configure
make -j4
```

The built binary will be in `src/openocd`.

