# Substratum

A very lightweight assistant tool for converting assembly instructions into hex, aiming to save some of the effort of writing
machine code by hand.

Substratum is explicitly NOT an assembler. It won't calculate offsets for jumps and won't be capable of handling labels;
all it can do is take an instruction such as `addi a0 x0 0x01` and return the corresponding hex machine code `13 05 10 00`.

Build with `make build`, run the CI process with `make ci` (requires staticcheck in the `$PATH`).
