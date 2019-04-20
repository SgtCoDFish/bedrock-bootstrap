# Bedrock Baremetal Tools

While we do want to minimise the amount of tooling we use so that we can minimise the amount of third-party code we need to rely on, it would be considerably harder to work without a few core tools.

## Building Qemu

If your system has qemu >= 3.1, then in theory RISC-V support was upstreamed and you can probably install a RISC-V-supporting qemu from your system package manager. In practise, at least on macOS, there are problems connecting to an emulated kernel with gdb. Better to build it yourself; it's not hard to do and it's actually quite quick to build.

You'll want to clone it yourself and build:

```bash
git clone --recursive --depth 1 https://github.com/qemu/QEMU
cd riscv-qemu
mkdir build && cd build
../configure --target-list=riscv32-softmmu
make -j4
```

## The GNU Toolchain

You can use a prebuilt toolchain from [SiFive](https://www.sifive.com/boards/) (search for "GCC Toolchain") or build your own. You'll need one in either case; the compilers won't be much use, but some of the other tools will be.

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

Note that we're using the arch `rv32ima` which means the base RISC-V 32-bit instruction set (`i`) plus extensions for `m`ultiplication and `a`tomic operations. You'll often see a `c` used in other toolchains; its omission is a conscious decision to simplify our efforts to write bare-metal machine code.

In any case, we can (and will) explicitly specify the arch later when building which will allow us to avoid using extensions even if the toolchain is built with support for various extensions.[1]

The executables will be in the `bin` dir.

## OpenOCD

OpenOCD is also available from SiFive as a prebuilt binary: [search "OpenOCD"](https://www.sifive.com/boards).

At the time of writing, RISC-V wasn't supported by any versions of OpenOCD which were installable from package managers like apt or homebrew. There is a [patched version](https://github.com/riscv/riscv-openocd) which is easy to build, however:

```bash
git clone https://github.com/riscv/riscv-openocd
cd riscv-openocd
./bootstrap
./configure
make -j4
```

The built binary will be in `src/openocd`.

### Notes

[1] In fact, it appears that even if one builds the toolchain without `c` support it still gets added in at the time of writing. That's why we'll explicitly not use it when building our applications.
