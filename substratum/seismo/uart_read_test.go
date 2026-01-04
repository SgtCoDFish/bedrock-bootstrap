package seismo

import (
	"testing"
	"time"
)

func TestUARTHello(t *testing.T) {
	rt := NewRISCVTest(t, "uart.elf")

	if err := rt.GDB.Continue(); err != nil {
		t.Fatal(err)
	}

	b, err := rt.UART.ReadByte(2 * time.Second)
	if err != nil {
		t.Fatalf("failed to read byte from UART: %s", err)
	}

	if b != '5' {
		t.Errorf("expected '5', got %q", b)
	}

	t.Log("verified UART output")

	rt.GDB.Halt()

	a0, err := rt.GDB.ReadReg(X10)
	if err != nil {
		t.Fatalf("failed to dump a0: %s", err)
	}
	t.Logf("a0  = 0x%08x", a0)

	// Inspect all registers
	regs, err := rt.GDB.ReadAllRegs()
	if err != nil {
		t.Fatalf("failed to dump regs: %s", err)
	}

	for r, v := range regs {
		t.Logf("r%-2d = 0x%08x", r, v)
	}

	const UART_TXDATA = 0x10013000
	val, err := rt.GDB.ReadMem(UART_TXDATA, 4)
	if err != nil {
		t.Errorf("failed to read TXDATA from memory: %s", err)
	}

	t.Logf("UART_TXDATA = 0x%08x", val)
}
