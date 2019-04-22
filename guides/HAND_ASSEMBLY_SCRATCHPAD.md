# Hand Assembly Scratchpad

This file holds the "working out" for RISC-V hand assembly. It's not documented much, and might be hard to follow, but it allows you to see how different commands are hand assembled.

Of great use for this exercise was the book [The RISC-V Reader](https://www.amazon.co.uk/RISC-V-Reader-Open-Architecture-Atlas/dp/0999249118), which has a great reference for both RV32 and RV64, including all the relevant extensions.

## Register Reference

```text
x0 == x00 == 0b00000
a0 == x10 == 0b01010
a1 == x11 == 0b01011
a2 == x12 == 0b01100
a5 == x15 == 0b01111
```

## Instruction Assembly Reference

We focus on actual instructions in RV32I; pseudoinstructions are shown in their "translated" form.

Hex values following a `=>` are little endian bytes ready to be written into a `.hex` file directly.

### LUI (U-type)

#### lui a2, 0x80000

```text
1000 0000 0000 0000 0000 0110 0011 0111
80        00        06        37

=> 37 06 00 80
```

#### lui a1, 0xDEAD0

```text
1101 1110 1010 1101 0000 0101 1011 0111
de        ad        05        b7

=> b7 05 ad de
```

### ADDI (I-type)

#### addi a5, a5, 0x3c

```text
0000 0011 1100 0111 1000 0111 1001 0011
03        c7        87        93

=> 93 87 c7 03
```

#### addi a1, x0, 0x3

```text
0000 0000 0011 0000 0000 0101 1001 0011
00        30        05        93

=> 93 05 30 00
```

#### addi a1, x0, 0x35

```text
0000 0011 0101 0000 0000 0101 1001 0011
03        50        05        93

=> 93 05 50 03
```

#### addi a0, x0, 0x8a

```text
0000 1000 1010 0000 0000 0101 0001 0011
08        a0        05        13

=> 13 05 a0 08
```

### SLLI (I-type)

#### slli a1, a1, 0x10

```text
0000 0001 0000 0101 0001 0101 0001 0011
01        05        15        13

=> 13 15 05 01
```

### SUB (R-type)

#### sub a1, x0, a1 (i.e. neg a1, a1)

```text
0100 0000 1011 0000 0000 0101 1011 0011
40        b0        05        b3

=> b3 05 b0 40
```

### AND (R-type)

#### and a0, a1, a0

```text
0000 0000 1010 0101 1111 0101 0011 0011
00        a5        f5        33

=> 33 f5 a5 00
```

#### and a0, a0, a2

```text
0000 0000 1100 0101 0111 0101 0011 0011
00        c5        75        33

=> 33 75 c5 00
```

### OR (R-type)

#### or a0, a0, a1

```text
0000 0000 1011 0101 0110 0101 0011 0011
00        b5        65        33

=> 33 65 b5 00
```

### LW (S-type)

#### lw a0, 0(a5)

```text
0000 0000 0000 0111 1010 0101 0000 0011
00        07        a5        03

=> 03 a5 07 00
```

### SW (S-type)

#### sw a0, 0(a5)

```text
0000 0000 1010 0111 1010 0000 0010 0011
00        a7        a0        23

=> 23 a0 a7 00
```

#### sw a1, 0(a5)

```text
0000 0000 1011 0111 1010 0000 0010 0011
00        b7        a0        23

=> 23 a0 b7 00
```

### BEQ (B-type)

#### bnez a0, -8 == bne a0, x0, -8

```text
1111 1110 0000 0101 0001 1100 1110 0011
fe        05        1c        e3

=> e3 1c 05 fe
```

#### beqz a0, -8 == beq a0, x0, -8

```text
1111 1110 0000 0101 0000 1100 1110 0011
fe        05        0c        e3

=> e3 0c 05 fe
```

### JAL (J-type)

#### j . == jal x0, 0x00

```text
0000 0000 0000 0000 0000 0000 0110 1111
00        00        00        6f

=> 6f 00 00 00
```
