# The Qemu/Hardware Shim

## Bridging the Differences

We know that our code will start running at `0x2040_0000` on Qemu and `0x2000_0000` on the HiFive1. So, we can add a couple of instructions to force compatibility between on _both_ systems (so that we only have to build one output binary that can run on both). We'll write code at `0x2000_0000` which always jumps to `0x2040_0000` and then write our actual code at `0x2040_0000`.

The code to do the jump is similar to what we've seen already:

```assembly
li t0, 0x20400000
jr t0
```

So all we need to do is ensure that those two instructions are placed at `0x2000_0000` and then our code is placed at `0x2040_0000`. Great; how do we do that? The answer (at least for now) is a linker script, which lets us specify where code ends up. [This blog](https://shobhitsharda.wordpress.com/2011/03/11/linker-scripts-inner-concepts/) is a great introduction to the basic concepts.

In fact, we already used a linker script earlier, in the form of `linker.ld` from the foundational section which put our program into RAM at `0x8000_0000`.
