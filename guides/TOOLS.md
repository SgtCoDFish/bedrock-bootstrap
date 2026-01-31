# Bedrock Baremetal Tools

While we do want to minimise the amount of tooling we use so that we can minimise the amount of third-party code we need to rely on, it would be considerably harder to work without a few core tools.

The main tools required if you want to follow along and write your own bedrock bootstrapping code are:

1. A text editor to write hex
2. `make` - should be available on almost any system
2. Some program to convert hex into binary (reverse hex dumping)

In this repo, the reverse hex dump tool of choice is `xxd`, which is commonly available by default or else can be found in Linux package managers or bundled with `vim` (which is the case in Homebrew).

The initial bootstrapping also uses `sed` to strip whitespace and comments from `.hex` files, which should be available on any POSIX system.

## QEMU

Unless you want to run everything on hardware (or if you're running through this repo on a RISC-V machine!) you'll need QEMU to be able to run the RISC-V code we generate.

If running on Linux or macOS it should be trivial to get a version of QEMU which supports RISC-V.

On macOS the `qemu` package in Homebrew includes `qemu-system-riscv32` which is all you need.

On Linux, as long as you're running a distro which was released in the last few years QEMU is almost certain to be available in your package manager with RISC-V support. For example, the `qemu-system-riscv` package is available in Debian Trixie and in Debian Bullseye backports.

## Binutils

There are a few tools commonly used in this repo which make it much easier to work with the low-level bootstrapping code.

The key tool is `objdump`. This is available in the `riscv64-elf-binutils` package in Homebrew, and is available in many package manager. Don't worry about specifically finding a `riscv32` objdump - the 64 bit version will work fine for us.

To be clear: you don't _need_ objdump - it just helps to check that the hex you're writing is correct.

## GDB

Debuggers are an incredibly useful set of tools, and using GDB to debug QEMU will make bootstrapping much simpler.

GDB should be available in package manager, and `riscv64-elf-gdb` is in Homebrew.

## Compilers

Binutils are much more important for our purposes than GCC. You can bootstrap a system without needing a compiler at all.

This repository does contain some assembly code for instructional or illustrative purposes and compiling that will require an assembler (`as`) and a linker (`ld`). Again, importantly, you can just choose to not bother with those - the point of bootstrapping is to avoid them.

These tools are also available in `riscv64-elf-binutils` in Homebrew or in package managers.
