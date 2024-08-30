# Pico 2 Research

## Background

The Raspberry Pi Pico 2 is a really interesting option for a bedrock bootstrap for one main reason: it's a Raspberry Pi!

Part of the risk (no pun intended) of many RISC-V boards is that even in 2024, it can feel like you're an early adopter taking a punt on a new technology provided by startups.

That's not to say that all the companies in the space will fail, but there's certainly a threat of that. Even if the companies don't fail, many of the boards suffer from the same problem that a lot of Raspberry Pi clones present - poor documentation.

It's very possible to get powerful hardware from a lot of vendors, but many of them fail at the software side of things and those boards can turn into a nightmare to find any details for.

Raspberry Pi doesn't have that problem. I can be confident that, now that the RP2350 chip is available with RISC-V support, there will be [excellent documentation](https://datasheets.raspberrypi.com/rp2350/rp2350-datasheet.pdf) and [easily accessible open source SDKs and examples](https://github.com/raspberrypi/pico-examples/tree/master?tab=readme-ov-file).

## Building Examples

As an easy way to get examples of how UART works on the Pico 2, we can use the Pico examples repo which produces
 all the files we need - a binary, an ELF file, the UF2 format which can easily be uploaded to a Pico and a disassembled dump to make code easier to browse.

First, we clone the repos we need and download a compiler. I'm checking out specific tags because that's how I first built locally; this helps to make these instructions easier to follow later.

As for the compiler, I'll provide a link to the one I downloaded [from here](https://www.embecosm.com/resources/tool-chain-downloads/#riscv-stable) which worked for me, but you might need a different one if that download goes away.

```console
git clone --branch sdk-2.0.0 git@github.com:raspberrypi/pico-examples.git && cd pico-examples

# warning: this download is over 1GB!
curl -LO https://buildbot.embecosm.com/job/riscv32-gcc-ubuntu2004-release/19/artifact/riscv32-embecosm-ubuntu2004-gcc13.2.0.tar.gz

sha256sum riscv32-embecosm-ubuntu2004-gcc13.2.0.tar.gz
# 508048f05a3ac8fd7058123b63a458bf8256a157c734d4cdff96defa768dd053  riscv32-embecosm-ubuntu2004-gcc13.2.0.tar.gz

tar xf riscv32-embecosm-ubuntu2004-gcc13.2.0.tar.gz

git clone --branch 2.0.0 git@github.com:raspberrypi/pico-sdk.git

mkdir build && cd build

cmake -DPICO_SDK_PATH=../pico-sdk -DPICO_PLATFORM=rp2350-riscv -DPICO_TOOLCHAIN_PATH=`pwd`/../riscv32-embecosm-ubuntu2004-gcc13.2.0 ..
make -j
```

The build will take a while and will "fail" because we haven't set up enough libraries to build all of the examples. We don't really care about that though - we just want the UART example and that should build without a problem.

```console
$ sha256sum uart/hello_uart/hello_uart.*
8ed837e3df7868800e6478f6b45635e370409a77911ba4400b1500fd8aa43dc3  uart/hello_uart/hello_uart.bin
f14f534354d0062873cddf1c60c8707d30a054be4f1e8506ebf29029364ecd41  uart/hello_uart/hello_uart.dis
d1949c882b8536d1568edf810d111e292748db213f7ff4098d853af8fdf56402  uart/hello_uart/hello_uart.elf
02c9a58674d054d2a65e6362c9b377abed2aab93d8fb242da41a09ac65df1fae  uart/hello_uart/hello_uart.elf.map
b67f6f8d636d01066dbce8758031da53a35f32037d4c6b62a11744867c59a902  uart/hello_uart/hello_uart.uf2
```

I've included a gzipped tarball of these 5 artifacts in this folder, since they compress very well.
