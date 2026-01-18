# `BB1`: Adding Support for Functions

It's not a great deal of fun to have to calculate offsets for jumps by hand, and even ignoring fun it's error prone and makes our programs brittle; if we add extra instructions in between the start and end of a jump, we need to recalculate offsets.

Being able to label sections of code and then later jump to those labels without needing to calculate offsets is a staple of assembly language. Adding support for functions for our code will make it much easier to write larger programs. To do this, we'll add the ability to define functions.

We'll start by adopting a convention: that registers will be used as in the RISC-V ABI with regards to their being preserved across calls. There's no compelling reason to deviate from that ABI.

## Constraints

Given our desire to target the HiFive1 and given that we're currently stuck to using RAM[1], we have a hard limit of 0x4000 (~16k) bytes of instructions which in RV32I translates to ~4k instructions in total.

[1] We could write a driver of some kind to use other kinds of memory, but that would be difficult even if we had full access to assembly. Better to stick to RAM for now.

## Function Syntax

We must add new syntax for defining and calling functions:

- `.` denotes a function definition. It must be followed by a single letter, which will be converted to lowercase and will be the name of the function. For example, `.M` defines a function called `m`.
- `$` is a function call - it should be followed by a single letter which is converted into lowercase and which denotes the function to be called. For example, `$M` calls the function called `m`

We're explicitly not defining a way to _end_ a function. Once we see e.g. `.A`, all following instructions are part of `a` until we see a different definition. This reduces complexity of parsing inputs.

## Function Addresses

Knowing the address of a function is a non-trivial topic in computer science; linkers can be quite involved, and there are a lot of ways to solve this particular problem. We have some main considerations, though:

1) We want to minimise the amount of code we write in BB0
2) We have 16kB of RAM to work with
3) Trying to avoid writing code implies that functions will tend to be smaller, and hex code is small anyway.

To avoid the world of lookup tables for function addresses or having to deal with linkers, we can at least for now define functions to have preset addresses.

If we define that functions begin at `0x8000_1000`, then we can define `a` to always have the address `0x8000_1000` and `b` to have address `0x8000_1200` and so on. That means each function could have 128 instructions in it (0x200 / 0x4 == 0x80 == 128), and we could fit 24 functions total (0x3000 / 0x200 == 0x18 == 24).

With the goal of leaving ourselves some spare capacity, we'll restrict ourselves to functions from 'a' to 'u', leaving the memory from 0x80003a00 to 0x80004000 free for whatever we might want to do later (such as a stack!)

While simple, our approach of defining static addresses for functions comes with real downsides which we need to be aware of:

1) All functions are the same size, so a function which is only a couple of instructions takes as much RAM as a complex function.
2) It's easy for one function to clobber another; with enough instructions, `a` can grow past the start of `b`.

We'll leave it to the programmer to remember these caveats and choose function names wisely.

## Code Layout

To keep implementation simple, we define the start of the program to be the start of the main function, written to the beginning of RAM at `0x8000_0000`.

Function definitions start at `0x8000_1000` as mentioned earlier. This implies a cap of 1024 (0x1000 / 4 == 1024) instructions in the "main" function before it clobbers the `a` function.
