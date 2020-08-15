# UART

Next thing is enable a connection and communication over [UART](https://en.wikipedia.org/wiki/Universal_asynchronous_receiver-transmitter) which will let us upload to and download from our target device.

The UART on the `sifive_e` QEMU board matches the UART on the HiFive1 so we should be able to write once and run on both platforms.

References:

- [Freedom Metal UART Driver](https://github.com/sifive/freedom-metal/blob/6d69e6d48babe4472a6f4671b832cb7df941f274/src/drivers/sifive%2Cuart0.c)
- [dwelch UART Driver](https://github.com/dwelch67/sifive_samples/blob/e93a68e4dfed9f0cc5e3d23cc4ac7c4176f15b98/hifive1/uart02/notmain.c)

## Hand Assembling

We're not going to list how to hand assemble every instruction, as it gets boring quickly and it's not hard to reason about how it works. There's a reason people use assemblers!

There are a couple of examples at [1] to help with the basic concepts, and there's a scratchpad of the "working out" in [HAND_ASSEMBLY_SCRATCHPAD.md](../guides/HAND_ASSEMBLY_SCRATCHPAD.md) with fuller working for many types of instructions.

## Running `uart.hex`

`uart.hex` is the end result of our UART example; it does nothing except output a single ASCII "5". We need to "compile" `BUILD/uart.elf` by adding an ELF header with the size of the program customised, removing the comments and reverse-hexdumping the file.

We can run on QEMU using the following slightly modified command:

```bash
$ qemu-system-riscv32 -nographic -serial pty -s -S -M sifive_e -kernel BUILD/uart.elf
QEMU 3.1.0 monitor - type 'help' for more information
char device redirected to /dev/ttys007 (label serial0)

# Note that the char device might change depending on your system.
# In any case, you can connect to it in a different terminal using:
screen /dev/ttys007 115200
```

After you've connected to the emulated serial port, go back to the `(qemu)` prompt in the QEMU window and type `c` for `(c)ontinue`, and then press enter. You should see a `5` output in the screen session.

Use `Ctrl+A` followed by `k` to quit the screen session, and then the `quit` command to exit QEMU.

## Understanding `uart.hex`

The hex file itself is commented heavily, to show both the reasoning behind each section and the individual assembly instructions which the machine code segments represent. These comments are stripped by `sed` in the Makefile.

(NB: This UART initialisation process was mostly reverse-engineered from existing code. If you know of a good guide which contains the details please do raise a PR!)

The program initialises the UART by:

- writing a device specific UART mask to a GPIO selector register
- writing the inverse of that mask to a GPIO enabler register to actually enable UART
- writing 138 to a UART divider register to select a baud rate of 115200 [2]
- writing 0x1 to a UART "txctrl" register to enable transmits [3]
- writing 0x1 to a UART "rxctrl" register to enable receives [4]

Finally we wait until the UART\_TXDATA register has its highest bit cleared and then write our value into the UART\_TXDATA register. The [datasheet](https://sifive.cdn.prismic.io/sifive%2F4d063bf8-3ae6-4db6-9843-ee9076ebadf7_fe310-g000.pdf) has details of this part:

> Writing to the `txdata` register enqueues the character contained in the data field to the transmit FIFO if the FIFO is able to accept new entries. Reading from `txdata` returns the current value of the full flag and zero in the data field.
> The full flag indicates whether the transmit FIFO is able to accept new entries; when set, writes to data are ignored.

There's also a hint in there:

> A RISC-V `amoswap` instruction can be used to both read the full status and attempt to enqueue data, with a non-zero return value indicating the character was not accepted.

We don't use this trick because we're only using RV32I instructions, and `amoswap` is in the `A` (atomic) extension.

## Next

Now that we've written a more complicated program - including a re-usable UART initialisation "driver", we can do something more complicated and write a program which takes input from UART.

## Notes

[1] Hand assembly notes:

### nop (I-type)

Our first attempt at assembling instructions by hand is to create a NOP. We'll use them liberally to pad out instructions and leave space for ourselves to add further instructions later.

From the [RISC-V spec](https://content.riscv.org/wp-content/uploads/2017/05/riscv-spec-v2.2.pdf) we see that:

> NOPs can be used to align code segments to microarchitecturally significant address boundaries, or to leave space for inline code modifications. Although there are many possible ways to encode a NOP, we define a canonical NOP encoding to  allow microarchitectural optimizations as well as for more readable disassembly output.

That canonical NOP encoding is `addi x0, x0, 0`, which is an "I-type" instruction. In chapter 19 of the spec, there's a handy chart for looking up the encodings of different instruction types and the opcodes required for different instructions.

`addi` has an opcode of `0010011` (note that's 7 bits, not 8). Since everything else in `addi x0, x0, 0` is 0 for a NOP (since x0 is represented as `0b00000` and the immediate value is zero too), the instruction is easy to represent: `0x0000_0013`, which we write in little endian as `13 00 00 00`.

### lui a5, 0x10012 (U-type)

`lui` is an instruction we've encountered before. The spec gives us a U-type encoding, and we want to load the value `0x10012000` into `a5`, the 5th argument register (for reasons that will become clear later). The assembly is `lui a5, 0x10012`.

From the guide we see the opcode `0110111`, and we need a 5-bit register (a5 is x15, so has the 5-bit encoding `01111`). Combined together the lowest 12 bits are `0111 1011 0111` (note that the lowest bit of the destination register is joined with the upper 3 bits of the opcode to create the nibble `1011`). The highest 20 bits are the immediate value, `0x10012`, so we have the complete instruction `0x100127b7` or `b7 27 01 10`.

[2] We write a value one less than the one we want (139). The whole thing is a bit confusing, but the details are in the data sheet.

[3] Some guides write `0x3` here, to use a second UART stop bit. The choice is explained in this [StackOverflow question](https://electronics.stackexchange.com/questions/29945/one-or-two-uart-stop-bits). We use one stop bit to increase throughput.

[4] We don't actually use the receive functionality in this example, though.
