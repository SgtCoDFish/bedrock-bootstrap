# Substratum

Collection of very lightweight assistant tools for bedrock baremetal development.

`substratum` is the overarching binary, with several commands, with each command also built as a separate binary containing just that command.

Substratum naturally depends on every dependency of every subcommand, but the individual binaries can have the truly minimal set of dependencies. This ensures that, say, the assembly helper remains easy to reimplement from scratch if desired while not limiting helper dependencies for the autotester.

Build all commands with `make build`, run the CI process with `make ci`.

## Commands

### ss-asm

`ss-asm` converts some RISC-V assembly instructions into hex, aiming to save some of the effort of writing machine code by hand.

`ss-asm` is explicitly **NOT** an assembler. It won't calculate offsets for jumps and won't be capable of handling labels;
all it can do is take an instruction such as `addi a0 x0 0x01` and return the corresponding hex machine code `13 05 10 00`.

### ss-autotest

Communicates with a running paused QEMU instance to check that the bedrock baremetal program loaded in that QEMU instance conforms to certain specifications.

For example, checks that the uart-rxxd program correctly translates hex received over serial into bytes in memory.
