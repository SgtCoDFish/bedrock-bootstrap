# uart-dump

`start.s` is a simple program which dumps itself over UART.

Run `make all` to build it then run it in QEMU to dump itself. You'll need to SIGINT to kill QEMU (Ctrl+C).

Run `make dump-asm` to see the dump of the compiled version.

Run `make dump-dumped` to see the dump of the dumped version.

This program is intentionally simple and dumb - it's a proof of concept.
