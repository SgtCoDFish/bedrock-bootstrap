# Hand Assembly Scratchpad

This file holds the "working out" for some RISC-V hand assembly. It's not documented much, and might be hard to follow, but it allows you to see how different commands are hand assembled. Most of the instructions are used in the basic UART machine code example; most instructions after that are "assembled" using substratum.

Of great use for this exercise was the book [The RISC-V Reader](https://www.amazon.co.uk/RISC-V-Reader-Open-Architecture-Atlas/dp/0999249118), which has a great reference for both RV32 and RV64, including all the relevant extensions.

## Register Reference

```text
x0 == x0 == 0b00000
a0 == x10 == 0b01010
a1 == x11 == 0b01011
a2 == x12 == 0b01100
a5 == x15 == 0b01111
```

## Instruction Assembly Reference

We focus on actual instructions in RV32I; pseudoinstructions are shown in their "translated" form.

Hex values following a `=>` are little endian bytes ready to be written into a `.hex` file directly.

### `LUI` (U-type)

#### `lui a2, 0x80000`

```text
1000 0000 0000 0000 0000 0110 0011 0111
80        00        06        37

=> 37 06 00 80
```

#### `lui a1, 0xDEAD0`

```text
1101 1110 1010 1101 0000 0101 1011 0111
de        ad        05        b7

=> b7 05 ad de
```

#### `lui x14, 0x80010`

```text

1000 0000 0000 0001 0000 0111 0011 0111
80        01        07        37

=> 37 07 01 80
```

#### `lui x15, 0x10013`

```text

0001 0000 0000 0001 0011 0111 1011 0111
10        01        37        b7

=> b7 37 01 10
```

### `ADDI` (I-type)

#### `addi a5, a5, 0x3c`

```text
0000 0011 1100 0111 1000 0111 1001 0011
03        c7        87        93

=> 93 87 c7 03
```

#### `addi a1, x0, 0x3`

```text
0000 0000 0011 0000 0000 0101 1001 0011
00        30        05        93

=> 93 05 30 00
```

#### `addi a1, x0, 0x35`

```text
0000 0011 0101 0000 0000 0101 1001 0011
03        50        05        93

=> 93 05 50 03
```

#### `addi a0, x0, 0x8a`

```text
0000 1000 1010 0000 0000 0101 0001 0011
08        a0        05        13

=> 13 05 a0 08
```

#### `addi x5, x0, 0x0`

```text
0000 0000 0000 0000 0000 0010 1001 0011
00        00        02        93

=> 93 02 00 00
```

#### `addi x6, x0, 0x0`

```text
0000 0000 0000 0000 0000 0011 0001 0011
00        00        03        13

=> 13 03 00 00
```

#### `addi x8, x0, 0x0`

```text
0000 0000 0000 0000 0000 0100 0001 0011
00        00        04        13

=> 13 04 00 00
```

#### `addi x16, x15, 0x4`

```text
0000 0000 0010 0111 1000 1000 0001 0011
00        47        88        13

=> 13 88 47 00
```

#### `addi x17, x0, 0xa`

```text
0000 0000 1010 0000 0000 1000 1001 0011
00        a0        08        93

=> 93 08 a0 00
```

### `SLLI` (I-type)

#### `slli a0, a0, 0x10`

```text
0000 0001 0000 0101 0001 0101 0001 0011
01        05        15        13

=> 13 15 05 01
```

### `SUB` (R-type)

#### `sub a1, x0, a1` (i.e. `neg a1, a1`)

```text
0100 0000 1011 0000 0000 0101 1011 0011
40        b0        05        b3

=> b3 05 b0 40
```

### `AND` (R-type)

#### `and a0, a1, a0`

```text
0000 0000 1010 0101 1111 0101 0011 0011
00        a5        f5        33

=> 33 f5 a5 00
```

#### `and a0, a0, a2`

```text
0000 0000 1100 0101 0111 0101 0011 0011
00        c5        75        33

=> 33 75 c5 00
```

### `ANDI` (I-type)

#### `andi x10, x10, 0xFF`

```text
0000 1111 1111 0101 0111 0101 0001 0011
0f        f5        75        13

=> 13 75 f5 0f
```

### `OR` (R-type)

#### `or a0, a0, a1`

```text
0000 0000 1011 0101 0110 0101 0011 0011
00        b5        65        33

=> 33 65 b5 00
```

### `LW` (S-type)

#### `lw a0, 0(a5)`

```text
0000 0000 0000 0111 1010 0101 0000 0011
00        07        a5        03

=> 03 a5 07 00
```

#### `lw a0, 0(a6)`

```text
0000 0000 0000 1000 0010 0101 0000 0011
00        08        25        03

=> 03 25 08 00
```

### `SW` (S-type)

#### `sw a0, 0(a5)`

```text
0000 0000 1010 0111 1010 0000 0010 0011
00        a7        a0        23

=> 23 a0 a7 00
```

#### `sw a1, 0(a5)`

```text
0000 0000 1011 0111 1010 0000 0010 0011
00        b7        a0        23

=> 23 a0 b7 00
```

### `BEQ` (B-type)

#### `bnez a0, -8` == `bne a0, x0, -8`

```text
1111 1110 0000 0101 0001 1100 1110 0011
fe        05        1c        e3

=> e3 1c 05 fe
```

#### `beqz a0, -8` == `beq a0, x0, -8`

```text
1111 1110 0000 0101 0000 1100 1110 0011
fe        05        0c        e3

=> e3 0c 05 fe
```

#### `beq x0, x0, -0x2c`

```text
1111 1100 0000 0000 0000 1010 1110 0011
fc        00        0a        e3

=> e3 0a 00 fc
```

#### `beq x8, x0, 0x20`

```text
0000 0010 0000 0100 0000 0010 0110 0011
02        04        02        63

=> 63 02 04 02
```

### `JAL` (J-type)

#### `j .` == `jal x0, 0x0`

```text
0000 0000 0000 0000 0000 0000 0110 1111
00        00        00        6f

=> 6f 00 00 00
```
