# UART Drivers

Creating a UART "hello world" program is great - we did something useful in bare machine code!

There was an important decision at the beginning of the last chapter, though: we targeted the `sifive_e` machine in QEMU, which has a different UART to the one present in the `virt` machine. Of course, there are a bunch of _other_ UARTs
available in the wild.

Obviously, we could just target one target machine type, but is it possible to target more?

## Aims

- Create a modular method of adding swappable code to our programs
