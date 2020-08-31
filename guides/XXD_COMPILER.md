# "Compiling" with `xxd`

`xxd` is a hexdump tool, useful for seeing the raw bytes which make up a program:

```bash
$ echo -n "0123" | xxd -p
30313233
```

We echo the ASCII string "0123" and `xxd` shows us the bytes in hex: `30 31 32 33`.

`-p` is used to show minimal, "plain" output.

Hexdumps are useful on their own, but our main usage comes from another option: `-r`, which performs a reverse hexdump. That is it turns ASCII hex into raw bytes:

```bash
$ echo -n "30 31 32 33" | xxd -p -r
0123
```

Note that "30" here is converted into the hex value `0x30` which is the ASCII character `0` (which is as we saw above). Use `man ascii` if you need an ASCII table.

A useful feature for our bare metal machine code programming is the ability to add comments. `xxd -p -r` doesn't handle comments natively:

```bash
$ echo -n "30 31 32 33 # 34 digits are easy in ASCII" | xxd -p -r
01234
```

Much of the comment was stripped, but the "34" after the comment was still picked up. To strip comments, we can use sed, which is also widely available:

```bash
$ echo -n "30 31 32 33 # 34 digits are easy in ASCII" | sed "s/#.*$//g" | xxd -r -p
0123
```

The `sed` command replaces (the `s` command) the pattern `#.*$` (which matches a `#` and then any characters up to the end of the line) with nothing, `g`lobally.

## Summary

Our initial "compiler" is a mix of `sed` and `xxd -r -p` which together turn `.hex` files into raw binaries.

Of course, if we wanted to assert that we'd written every single line of code ourselves we could write our own "reverse hexdump" tool, which is very easy to write in any major programming language. We choose to use `xxd -r -p` since it's quite common (i.e. available on macOS and Linux) and simple.

(Writing our own `sed` would be a much bigger ask, but if you don't want comments, it's easy enough to strip them manually!)
