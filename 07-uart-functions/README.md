# BB1: Adding Support for Functions

It's not a great deal of fun to have to calculate offsets for jumps by hand, and even ignoring the lack of fun it's error prone and makes our programs brittle; if we add extra instructions in between the start and end of a jump, we might need to recalculate offsets.

Being able to label sections of code - and then later jump to those labels and have the assembler calculate the offset for us - is a staple of assembly languge, and adopting it for our code will make it much easier to write larger programs. To do this, we'll add the ability to define up to 26 functions, each with a single ASCII letter as a label.

We'll start by adopting a convention: that registers will be used as in the RISC-V ABI with regards to their being preserved across calls. There's no compelling reason to deviate from that ABI.

We also need to make a choice about what kind of program layout we'll adopt.

## Functions First

One option is to have everything be defined in functions. We can choose a convention whereby the `main` or `_start` function is labelled `z`, and insert an unconditional jump to that function when we write the program. This approach broadly mirrors how a higher-level language such as C or assembly would work.

The downside to this approach is that no code written for BB0 - i.e. "raw" ASCII machine code - would work without changes. That said, who's writing a lot of code in BB0?

## Code First

If a function `A` is defined by a prelude such as `.A`, we can have BB0 code optionally interspersed with function calls such as "xA". Functions - if there are any - can be defined after the initial block. That is, the "main" block, such that there is one, would consist of any code up until the first function definition.

This allows BB0 code to run as-is in BB1. This property of preserving the ability to run code from previous bootstrapping stages is a neat element of purity, and also by default provides us with several ready-to-run tests of any new stage - that is, the code of a previous stage!

## Which to Use?

There's not really a correct choice here - both approaches would work and have different advantages. In practise, the differences in testability between the different approaches are minimal, since we'd use an automated testbed in `substratum`. The tradeoff then is between the approach which has a neat sense of purity - "code first" - and the "functions first" code which more closely matches the type of programs that we are likely to write in future bootstrap stages.

On balance, we choose to go with "functions first". Working pragmatically, the only purpose of any given bootstrap step (except the last one!) is to enable the next step in the process. Being able to write programs in any given step is a neat distraction, but practically there's little point writing in BB0 when we can do a full bootstrap and use assembly or a high-level language. As such, preserving backwards compatibility - while cool - isn't hugely useful, and it's better to prefer a more easy-to-read approach which minimises implicit assumptions.

## Implementation

We define some new control characters:

- `.` denotes a function. It must be followed by a single letter, which will be converted to lowercase and will be the name of the function. For example, `.Z` defines a function called `z`.
- `x` is a function call - it should be followed by a single letter which is converted into lowercase and which denotes the function to be called. For example, `xX` calls the function called `x`

We define `z` to be the main function. Every BB1 program must have `z`. The output of the program will start with an unconditional jump to the start of `z`, followed by an infinite loop in case `z` returns.

To simplify further, we _do not allow_ function calls from one function to another _except_ for `z`, which may call any function. So, `a` cannot call `b` or `z`, but `z` can call `a` and `b`. This simplifies the calculation of offsets, since we only need to do it for `z`.

In addition, we require that `z` is defined last in the file, and is followed by a `j` character which will cause execution to jump to the start of the program received over UART.

Otherwise, the syntax is the same as BB0.
