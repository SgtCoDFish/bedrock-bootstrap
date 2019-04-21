# UART

Next thing is enable a connection and communication over [UART](https://en.wikipedia.org/wiki/Universal_asynchronous_receiver-transmitter) which will let us upload to and download from our target device.

The UART on the `sifive_e` QEMU board matches the UART on the HiFive1 so we should be able to write once and run on both platforms.

References:
- [Freedom Metal UART Driver](https://github.com/sifive/freedom-metal/blob/6d69e6d48babe4472a6f4671b832cb7df941f274/src/drivers/sifive%2Cuart0.c)
- [dwelch UART Driver](https://github.com/dwelch67/sifive_samples/blob/e93a68e4dfed9f0cc5e3d23cc4ac7c4176f15b98/hifive1/uart02/notmain.c)

## Hand Assembling

We're not going to list how to hand assemble every instruction, as it get boring quickly and it's not hard to reason about how it works - which is why people normally use assemblers for this kind of task!

It does help to have a couple of examples which we can use to get a feel for hand-assembly, though.

### NOP (I-type)

Our first attempt at assembling instructions by hand is to create a NOP. We'll use them liberally to pad out instructions and leave space in case we need to jump to a different address.

From the [RISC-V spec](https://content.riscv.org/wp-content/uploads/2017/05/riscv-spec-v2.2.pdf) we see that:

> NOPs can be used to align code segments to microarchitecturally significant address boundaries, or to leave space for inline code modifications. Although there are many possible ways to encode a NOP, we define a canonical NOP encoding to  allow microarchitectural optimizations as well as for more readable disassembly output.

That canonical NOP encoding is `addi x0, x0, 0`, which is an "I-type" instruction. In chapter 19 of the spec, there's a handy chart for looking up the encodings of different instruction types and the opcodes required for different instructions.

`addi` has an opcode of `0010011` (note that's 7 bits, not 8). Since everything else in `addi x0, x0, 0` is 0 for a NOP, the instruction is easy to represent: `0x0000_0013`, which we write in little endian as `13 00 00 00`.

### LUI a5, 0x10012 (U-type)

`lui` is an instruction we've encountered before. The spec gives us a U-type encoding, and we want to load the value `0x10012000` into `a5`, the 5th argument register (for reasons that will become clear later). The assembly is `lui a5, 0x10012`.

From the guide we see the opcode `0110111`, and we need a 5-bit register (a5 is x15, so has the 5-bit encoding `01111`). Combined together the lowest 12 bits are `0111 1011 0111` (note that the lowest bit of the destination register is joined with the upper 3 bits of the opcode to create the nibble `1011`). The highest 20 bits are the immediate value, `0x10012`, so we have the complete instruction `0x100127b7` or `b7 27 01 10`.

#### LUI a1, 0xDEAD0

1101 1110 1010 1101 0000 0101 1011 0111
de        ad        05        b7

b7 05 ad de

#### LUI a2, 0x80000
1000 0000 0000 0000 0000 0110 0011 0111
80        00        06        37

37 06 00 80

### ADDI (I-type)

#### addi a5, a5, 0x3c
a5 == x15 == 0b01111

0000 0011 1100 0111 1000 0111 1001 0011
imm----------- rs1---000 rd----opcode--
03        c7        87        93

93 87 c7 03

#### addi a1, x0, 0x3
a1 = x11 == 0b01011

0000 0000 0011 0000 0000 0101 1001 0011
00        30        05        93

93 05 30 00

#### addi a0, x0, 0x8a

0000 1000 1010 0000 0000 0101 0001 0011
08        a0        05        13

13 05 a0 08


### LW a0, 0(a5)

a0 == x10 == 0b01010

0000 0000 0000 0111 1010 0101 0000 0011
00        07        a5        03

03 a5 07 00

### SLLI (I-type)

#### slli a1, a1, 0x10

0000 0001 0000 0101 0001 0101 0001 0011
01        05        15        13

13 15 05 01

### SUB (R-type)

#### sub a1, x0, a1 (i.e. neg a1, a1)

0100 0000 1011 0000 0000 0101 1011 0011
40        b0        05        b3

b3 05 b0 40

### AND (R-type)

#### and a0, a1, a0

0000 0000 1010 0101 1111 0101 0011 0011
00        a5        f5        33

33 f5 a5 00

#### and a0, a0, a2
a2 = x12 = 01100

0000 0000 1100 0101 0111 0101 0011 0011
00        c5        75        33

### SW (S-type)

#### sw a0, 0(a5)

0000 0000 1010 0111 1010 0000 0010 0011
00        a7        a0        23

23 a0 a7 00

### OR (R-type)

#### or a0, a0, a1
0000 0000 1011 0101 0110 0101 0011 0011
00        b5        65        33

33 65 b5 00

### BEQ (B-type)

### beqz a0, -8 == beq a0, x0, -8

1111 1110 0000 0101 0000 1100 1110 0011
fe        05        0c        e3


### JAL (J-type)

#### j . == jal x0, 0x00

0000 0000 0000 0000 0000 0000 0110 1111
00        00        00        6f

6f 00 00 00
