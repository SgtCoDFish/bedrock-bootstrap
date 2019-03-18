# Foundations

Before we can do _anything_ we'll need to gather some required documentation, and find a way to debug the code we're going to write.

The [HiFive1 getting started guide](https://sifive.cdn.prismic.io/sifive%2F9c57065b-6d28-465b-b67d-f416894123a9_hifive1-getting-started-v1.0.2.pdf) is actually not the best place for us to get started. It uses a fairly heavyweight IDE and isn't targetted at true bare metal. There is however a [more detailed datasheet](https://sifive.cdn.prismic.io/sifive%2F4d063bf8-3ae6-4db6-9843-ee9076ebadf7_fe310-g000.pdf) which will be a useful reference for us.

[Freedom Metal](https://sifive.github.io/freedom-metal-docs/introduction.html#what-is-freedom-metal) (on [GitHub](https://github.com/sifive/freedom-metal/tree/master)) looks more useful, as does the [freedom-e-sdk](https://github.com/sifive/freedom-e-sdk). We see they provide
a [Board Support Package](https://github.com/sifive/freedom-e-sdk/tree/30c143eb5445f47edb351ba54c84ff8285dc27a9/bsp/sifive-hifive1) for the HiFive1, with details we need such as the architecture we're compiling for. We also see example programs
such as the vanishingly simple [hello](https://github.com/sifive/example-hello/tree/d1397bec64187efb8b791fe1eb307aa3c760c694) which we can use as a sanity check. The Makefile is a big beast designed to support multiple different boards in testing, and we don't need most of it.

We're not actually going to be _using_ Freedom Metal library directly, but since it's open source we can learn from it. Note also that it's in C, which isn't any use to us in a bedrock bare metal world.

While you're at it you'll probably want to clone the `freedom-e-sdk` just to have it around.

```bash
git clone --recursive https://github.com/riscv/riscv-gnu-toolchain
git clone --recursive https://github.com/sifive/freedom-e-sdk
```

Once that's done we'll need to build.

## Building Qemu

If your system has qemu >= 3.1, then RISC-V support was upstreamed and you can probably install a RISC-V-supporting qemu from your system package manager.

Otherwise you'll want to clone it yourself and build:

```bash
git clone --recursive https://github.com/sifive/riscv-qemu
cd riscv-qemu
mkdir build && cd build
../configure --target-list=riscv32-softmmu # TODO: Check this is right
make -j4
```

## The GNU Toolchain

You can use a prebuilt toolchain from [Sifive](https://www.sifive.com/boards/) (search for "GCC Toolchain") or build your own. You'll need one in either case; the compilers won't be much use, but some of the other tools will be.

If you're building from source, you can see in the freedom-e-sdk HiFive1 [BSP](https://github.com/sifive/freedom-e-sdk/blob/30c143eb5445f47edb351ba54c84ff8285dc27a9/bsp/sifive-hifive1/settings.mk) that we need to target a different arch and ABI, since the toolchain defaults to 64 bit.

You can choose a different prefix, but we'll assume you've set `$RISCV_PREFIX` to something.

```bash
cd riscv-gnu-toolchain
mkdir build && cd build
../configure --with-arch=rv32imac --with-abi=ilp32 --with-cmodel=medlow --prefix=$RISCV_PREFIX 
make -j4  # might need to be done as root depending on where you're installing
```

The build files will be in the `$RISCV_PREFIX/bin` dir.

## Compiling Our Sample

Now we have a toolchain, we can test it using the files in this repo. Just run:

```bash
make RISCV_PREFIX=/path/to/riscv all
```

and several files will be dumped into BUILD in this directory. The most interesting ones from our perspective are `nothing.bin` and `nothing.dump`:

```bash
$ cat BUILD/nothing.elf
BUILD/nothing.elf:     file format elf32-littleriscv

Disassembly of section .text:

80000000 <main>:
80000000: 452d        li a0,11
80000002: 8082        ret

80000004 <_start>:
80000004: 80004137    lui sp,0x80004
80000008: ff9ff0ef    jal ra,80000000 <main>
8000000c: 9002        ebreak
8000000e: a001        j 8000000e <_start+0xa>

$ hexdump BUILD/nothing.bin
0000000 4501 8082 4137 8000 f0ef ff9f 9002 a001
0000010
```

Note that the raw binary dump (`nothing.bin`) is very small, and that the bytes match up with the disassembly in `nothing.dump`. For example, we see that the first instruction, at address `0x80000000` has the hex machine-code representation `452d` and that this matches the first 2 bytes of `nothing.bin`. The dump also has assembly on the right to make it easier to read, so we can see that the sole command is `li a0,11` which loads the value `11` into the register named `a0` (which is designated by the RISC-V ABI as being for return values).

Of course, the compiled output isn't much use to us right now. We want to run it on qemu to make sure it actually runs.

## Running in Qemu

Note that when running Qemu headless, you exit by pressing Ctrl+A and then `x` (use Ctrl+A and then `h` for help).

We pass a few arguments to the following command which look initially confusing:

- `-s` starts a GDB debugger on port `1234` which lets us dump memory
- `-S` pauses the CPU before running anything, giving us time to debug
- `-machine sifive_e` tells qemu we're running a sifive-e machine (which the HiFive 1 is!)
- `-nographic` disables graphics - we're not going to need them.

Finally we point at our kernel (which is our ELF file)

```bash
$ qemu-system-riscv32 -machine sifive_e -nographic -s -S -kernel BUILD/nothing.elf
# Open a separate terminal
$ $RISCV_PREFIX/bin/riscv32-unknown-elf-gdb
(gdb) target remote localhost:1234
(gdb) x 0x80000000
0x80000000:    0x8082452d
(gdb) quit
```

`x 0x80000000` dumps the memory at that address, which we can see is our `main` function above. We know we've loaded the file correctly into qemu.

# TODO: Actually run the program

## Other Links

- dwelch67 always has good guides and you can check his [uart01 sample](https://github.com/dwelch67/sifive_samples/tree/master/hifive1/uart01) for a baremetal UART example, which is getting a little ahead of ourselves if we're focusing on bedrock.
- Running RISC-V on qemu bare metal [google groups thread](https://groups.google.com/a/groups.riscv.org/forum/#!topic/sw-dev/IET9LBFJohU)
- A very minimal bare metal example, similar to dwelch but based on riscv64 and SPIKE: [schoeberl](https://github.com/schoeberl/cae-examples)
