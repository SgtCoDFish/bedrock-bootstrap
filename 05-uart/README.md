# UART

Next thing is enable a connection and communication over [UART](https://en.wikipedia.org/wiki/Universal_asynchronous_receiver-transmitter) which will let us upload to and download from our target device.

The UART on the `sifive_e` QEMU board matches the UART on the HiFive1 so we should be able to write once and run on both platforms.

References:
- [Freedom Metal UART Driver](https://github.com/sifive/freedom-metal/blob/6d69e6d48babe4472a6f4671b832cb7df941f274/src/drivers/sifive%2Cuart0.c)
- [dwelch UART Driver](https://github.com/dwelch67/sifive_samples/blob/e93a68e4dfed9f0cc5e3d23cc4ac7c4176f15b98/hifive1/uart02/notmain.c)
