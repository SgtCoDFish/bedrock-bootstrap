# `xxd -r` Over UART

Our aim is to (eventually!) bootstrap a simple compiler for a higher-level language. A sensible place to start is the "compiler" we use to convert `.hex` files to binaries - `xxd -r -p`, which will let us self-host the "0th" stage of our bootstrapping process.

We'll accept a stream of ASCII hex chars over UART and later convert them into binary, storing them in RAM as we go.

- If we read a pound `#` (`0x23`) we'll ignore everything until we read a newline `\n` (`0x0A`).
- If we read an `0xFF` we'll:
- - Write, at the end of the RAM area to which we've written code an unconditional jump back to the beginning of this the uart-rxxd program (likely at `0x2040_0000`)
- - Jump to the start of the area in RAM to which we've written code
- If we read anything which isn't mentioned above or else correspond to a hexadecimal character (`[0-9a-fA-F]`), we'll ignore it.
- When we receive the first hexadecimal nibble of a byte, we'll convert it to a number, shift left by 4 bits and store it in a temp register.
- When we receive the second hexadecimal nibble of a byte, we'll:
- - Convert it to a number
- - Add it to the temp register
- - Store that register in RAM at a pre-defined pointer
- - Increase the pointer
- - Clear the register

## Aims

- Write a more complicated program in machine code which can "compile" itself

## Substratum

Assembling instructions by hand is no fun at all beyond a few instructions, and we're starting to need quite a lot more instructions for our code. We don't want to go the whole way of using an actual assembler, so we've provided a very small assembly-to-machine-code converter called [substratum](../substratum/README.md).

Substratum _isn't_ an assembler. It won't support calculating offsets, it doesn't support storing data or labels or anything. It takes an RV32I instruction and, assuming it knows about that instruction, returns the 4-byte machine code representation. The cleverest thing it does is support both types of register naming (e.g. `a0` == `x10`).

We can use substratum to help us write machine code quicker. This comes in handy not just here, but also in future "compilers" where we'll still have to write some machine code.
