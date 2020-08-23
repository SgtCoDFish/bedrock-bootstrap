# BB1: Adding Support for Functions

It's not a great deal of fun to have to calculate offsets for jumps by hand, and even ignoring the lack of fun it's error prone and makes our programs brittle; if we add extra instructions in between the start and end of a jump, we might need to recalculate offsets.

Being able to label sections of code - and then later jump to those labels and have the assembler calculate the offset for us - is a staple of assembly language, and adopting it for our code will make it much easier to write larger programs. To do this, we'll add the ability to define functions, each with a single ASCII letter as a label.

We'll start by adopting a convention: that registers will be used as in the RISC-V ABI with regards to their being preserved across calls. There's no compelling reason to deviate from that ABI.

We also need to make a choice about what kind of program layout we'll adopt.

## Functions First

One option is to have everything be defined in functions. We can choose a convention whereby the `main` or `_start` function has a label by convention, say `m`, and unconditionally jump to that function when we jump to the program. This approach broadly mirrors how a higher-level language such as C or assembly would work.

The downside to this approach is that no code written for BB0 - i.e. "raw" ASCII machine code - would work without changes. That said, who's writing a lot of code in BB0?

## Code First

If a function `A` is defined by a prelude such as `.A`, we can have BB0 code optionally interspersed with function calls such as "xA". Functions - if there are any - can be defined after the initial block. That is, the "main" block, such that there is one, would consist of any code up until the first function definition.

This allows BB0 code to run as-is in BB1. This property of preserving the ability to run code from previous bootstrapping stages is a neat element of purity, and also by default provides us with several ready-to-run tests of any new stage - that is, the code of a previous stage!

## Which to Use?

There's not really a correct choice here - both approaches would work and have different advantages. In practise, the differences in testability between the different approaches are minimal, since we'd use an automated test-bed in `substratum`. The trade-off then is between the approach which has a neat sense of purity - "code first" - and the "functions first" code which more closely matches the type of programs that we are likely to write in future bootstrap stages.

On balance, we choose to go with "functions first". Working pragmatically, the only purpose of any given bootstrap step (except the last one!) is to enable the next step in the process. Being able to write programs in any given step is a neat distraction, but practically there's little point writing in BB0 when we can do a full bootstrap and use assembly or a high-level language. As such, preserving backwards compatibility - while cool - isn't hugely useful, and it's better to prefer a more easy-to-read approach which minimises implicit assumptions.

## Implementation

We define some new control characters:

- `.` denotes a function. It must be followed by a single letter, which will be converted to lowercase and will be the name of the function. For example, `.M` defines a function called `m`.
- `x` is a function call - it should be followed by a single letter which is converted into lowercase and which denotes the function to be called. For example, `xM` calls the function called `m`

We define `m` to be the main function. Every BB1 program must have `m`, and the bootstrapper will jump to `m` to start the program.

Given our desire to target the HiFive1, and the advantage of being able to stay within the on-chip memory for now (which removes the need to write any complicated "driver" code for the flash memory), we have a hard limit of 0x1000 (4096) bytes of instructions, which in RV32I translates to 1024 instructions. To avoid having to build a hash table which points to the start of each function on disk, we blindly assign each function from A-L 0x138 (312) bytes each, which gives an upper limit of 78 instructions in each function.

The main function is therefore located at the end, and has an extra 40 bytes (10 instructions) worth of space.
