# `xxd -r` Over UART

Our aim is to (eventually!) bootstrap a simple compiler for a higher-level language. A sensible place to start is the "compiler" we use to convert `.hex` files to binaries - `xxd -r -p`, which will let us self-host the "0th" stage of our bootstrapping process.

We'll accept a stream of ASCII hex chars over UART and then convert them into binary, storing them in RAM as we go. For example, given the UART input `13 00 00 00` (a no-op instruction) we'll end up writing `0x13000000` at `0x80001000` in memory. Given further input of `12 34 56 78` we'll write `0x12345678` at `0x80001004`.

## Aims

- Write a more complicated program in machine code which can "compile" itself

## Substratum

Assembling instructions by hand is no fun at all beyond a few instructions, and we're starting to need quite a lot more instructions for our code. We don't want to go the whole way of using an actual assembler, so we've provided a very small assembly-to-machine-code converter called [substratum](../substratum/README.md).

Substratum _isn't_ an assembler. It won't support calculating offsets for jump instructions, it doesn't support storing data or labels or generally anything which is provided by a regular assembler. It takes an RV32I instruction and, assuming it knows about that instruction, returns the 4-byte machine code representation. The cleverest thing it does is support both types of register naming (e.g. `a0` == `x10`).

We can use substratum to help us write machine code quicker. This comes in handy not just here, but also in future "steps" where we'll still have to write some machine code.

## Debugging What We Wrote

It's fine to write some instructions into a file, but eventually you're going to want to run them and, inevitably, debug them when things go wrong. Writing in machine code is error prone, after all!

You'll want 3 terminal windows; one with riscv-qemu, started something like this:

```bash
$ qemu-system-riscv32 -nographic -serial pty -s -S -M sifive_e -kernel BUILD/uart-rxxd.elf
...
char device redirected to /dev/ttys005 (label serial0)
```

Note the output about the serial port, which will likely differ depending on your OS and could change between runs.

Second, you'll need to connect to that serial port using GNU screen for sending input over the serial device:

```bash
$ screen /dev/ttys005 115200 # note the same serial port as above!
# no output expected, but will be used later
```

Finally, you'll want to connect to qemu using gdb:

```bash
$ riscv32-unknown-elf-gdb -q -ex "set architecture riscv:rv32" -ex "target remote :1234"
Remote debugging using :1234
warning: No executable has been specified and target does not support
determining executable automatically.  Try using the "file" command.
0x00001000 in ?? ()
```

The gdb warning is just telling us that there's no debugging information for the target. It would take a non-trivial amount of work to add that to our ELF headers and it might not even work properly if we did add it, so we'll have to take a more manual, analytical approach.

### Navigating in GDB

The main command you'll need is `si` or `stepi` which means step instruction and will jump to the next instruction, following jumps if needed. It takes an optional argument which will allow repeats, so `si 50` will jump forward 50 instructions.

The other main instructions have been mentioned before: `i r` to dump registers and `x <addr>` to dump what's in memory at the given address.

### Sending Over Serial

We'll need a way to send data over UART to be read by our program. Say we have a file, `/tmp/t.hex` which contains only a single no-op instruction (NOP) which is encoded as follows: `13000000`.

To send it, we go to our `screen` session, and press `Ctrl+A` followed by `:` which opens a command prompt. We type `readreg p /tmp/t.hex` which loads the file into a register called `p` but doesn't send it.

We then do `Ctrl+A :` again, and type `paste p` which sends the contents of `p` over UART to our running program.

### Debugging an Actual Bug

As an example of the debugging process, we'll run through how one of the bugs in one of the initial versions of uart-rxxd.hex was fixed.

It was noticed that the program entered "comment mode" incorrectly when no `#` character had been sent over UART. We open gdb and step using `si 180`  to reach `0x204000b0` which is where the program reads from UART. At that point we start single stepping using `si` to watch the flow.

We send `t.hex` above, which contains a single no-op, and run `si` until we get to around `0x204000f0` which we see from `make dump` is where we handle comment mode.

We check the registers with `i r` and see that `s0` is 0, as expected. We `si` twice, run `i r` again and see the value has been set!

We inspect the source and see that here was an incorrect instruction at `0x204000f4`, which caused us to turn on comment mode only when the character read was _not_ `#`. It's easily fixed:

```text
63 08 95 00  # beq x10, x09, +0x10 - incorrect!
63 1a 95 00  # bne x10 x09 0x14    - correct
 ```

If this seems clunky, that's because it is! But pretty quickly it becomes natural to browse using gdb and check the control flow of the program.

## Autotest

Of course, we don't need to manually debug every program. We can run it, and assert that it does what we expect programmatically. This is another feature of substratum!

```bash
substratum autotest -serial /dev/serialdevicename -test-name uart-rxxd-basic
substratum autotest -serial /dev/serialdevicename -test-name uart-rxxd-comment
substratum autotest -serial /dev/serialdevicename -test-name uart-rxxd-full
```

autotest will connect over a serial connection (which you'll need to specify) and then use GDB to run the program. For this program - uart-rxxd - it'll confirm that the correct input is being read by the program (i.e. that an ASCII value of `1` is correctly loaded into `x05` as `0x31`) and that the input is being written correctly when a full word is received. That's all the `basic` test does.

The `comment` test is similar to the basic test, but also confirms that comments are handled correctly.

The `full` test is more rigorous: it has multiple comments, multiple words and more checks.

If you want to try writing your own uart-rxxd implementation, you can still use autotest to confirm that it works, although some underlying assumptions - e.g. that received words are written to `0x8000_1000` - are baked in.
