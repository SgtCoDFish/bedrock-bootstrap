package seismo

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/sgtcodfish/substratum"
)

func TestUARTLoadsNOPInstruction(t *testing.T) {
	rt := NewRISCVTest(t, "uart-rxxd.elf")

	if err := rt.GDB.Continue(); err != nil {
		t.Fatalf("failed to continue: %v", err)
	}

	input := []byte("13 00 00 00")
	if err := rt.UART.Write(input); err != nil {
		t.Fatalf("failed to write UART input: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	if err := rt.GDB.Halt(); err != nil {
		t.Fatalf("failed to halt CPU: %v", err)
	}

	const RAM_BASE = 0x80000000
	data, err := rt.GDB.ReadMem(RAM_BASE, 4)
	if err != nil {
		t.Fatalf("failed to read RAM_BASE from memory: %s", err)
	}

	got := binary.LittleEndian.Uint32(data)
	want := uint32(0x00000013)

	if got != want {
		t.Fatalf("memory[0x%08x] = 0x%08x, want 0x%08x",
			RAM_BASE, got, want)
	}

	t.Logf("PASS: wrote NOP instruction 0x%08x to RAM", got)
}

func TestUARTLoadsMultipleInstructions(t *testing.T) {
	rt := NewRISCVTest(t, "uart-rxxd.elf")

	if err := rt.GDB.Continue(); err != nil {
		t.Fatalf("failed to continue: %v", err)
	}

	expectedX2 := uint32(rand.Int31n(255))
	x2ASM := fmt.Sprintf("addi x2 x0 %d", expectedX2)
	assembledX2Raw, err := substratum.AssembleLine(x2ASM)
	if err != nil {
		t.Fatalf("failed to assemble %s: %s", x2ASM, err)
	}

	assembledX2 := []byte(fmt.Sprintf("%02x %02x %02x %02x", assembledX2Raw[0], assembledX2Raw[1], assembledX2Raw[2], assembledX2Raw[3]))

	t.Logf("assembled %q (random = 0x%02x) to %s", x2ASM, expectedX2, assembledX2)

	x2Memory := binary.LittleEndian.Uint32(assembledX2Raw)

	// Program to send:
	// 5x NOP
	// clear x1
	// addi x1, x1, 5
	input := []byte(
		"13 00 00 00 " + // nop
			"# a comment\n" +
			"13 00 00 00 \n" +
			"13 00 00 00 " +
			"13 00 00 00 " +
			"13 00 00 00 " +
			"93 00 00 00 " + // addi x1, x0, 0
			"93 80 50 00 " + // addi x1, x1, 5
			string(assembledX2) + // addi x2, x0, <rand>
			"j",
	)

	if err := rt.UART.Write(input); err != nil {
		t.Fatalf("failed to write UART input: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if err := rt.GDB.Halt(); err != nil {
		t.Fatalf("failed to halt CPU: %v", err)
	}

	const base = 0x80000000

	expected := []uint32{
		0x00000013,
		0x00000013,
		0x00000013,
		0x00000013,
		0x00000013,
		0x00000093,
		0x00508093,
		x2Memory,
	}

	for i, want := range expected {
		addr := base + uint32(i*4)
		data, err := rt.GDB.ReadMem(addr, 4)
		if err != nil {
			t.Fatalf("failed to read from 0x%08x: %s", addr, err)
		}

		got := binary.LittleEndian.Uint32(data)

		if got != want {
			t.Fatalf(
				"instruction %d @ 0x%08x = 0x%08x, want 0x%08x",
				i, addr, got, want,
			)
		}

		t.Logf("OK: instruction %d @ 0x%08x = 0x%08x",
			i, addr, got)
	}

	x1, err := rt.GDB.ReadReg(X1)
	if err != nil {
		t.Fatalf("failed to dump x1: %s", err)
	}

	expectedX1 := uint32(5)

	if x1 != expectedX1 {
		t.Fatalf("expected x1 to be 0x%08x but got 0x%08x", expectedX1, x1)
	}

	t.Logf("OK: x1 == 0x%08x", x1)

	x2, err := rt.GDB.ReadReg(X2)
	if err != nil {
		t.Fatalf("failed to dump x2: %s", err)
	}

	if x2 != expectedX2 {
		t.Fatalf("expected x2 to be 0x%08x but got 0x%08x", expectedX2, x2)
	}

	t.Logf("OK: x2 == 0x%08x", x2)
}
