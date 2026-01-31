# Compiling Bare Metal C Code

We intentionally avoid using C in this repo, but in case you want to compile bare-metal C code:

```make
BUILD/nothing.o: nothing.c
	@mkdir -p BUILD
	$(RISCV_PREFIX)gcc -fPIC -O2 -g0 -static -fvisibility=hidden -nostdlib -nostartfiles -march=rv32i -mabi=ilp32 -c $< -o $@

BUILD:
	@mkdir -p $@
```
