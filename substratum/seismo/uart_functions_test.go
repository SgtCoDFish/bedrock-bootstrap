package seismo

import (
	"encoding/binary"
	"os"
	"testing"
	"time"
)

func TestUARTFunctionsBasic(t *testing.T) {
	rt := NewRISCVTest(t, "../../07-uart-functions/BUILD/uart-functions.elf")

	if err := rt.GDB.Continue(); err != nil {
		t.Fatalf("failed to continue: %v", err)
	}

	testFile := "uart-functions-test1.hex1"
	input, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("failed to read %s: %s", testFile, err)
	}

	if err := rt.UART.Write(input); err != nil {
		t.Fatalf("failed to write UART input: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if err := rt.GDB.Halt(); err != nil {
		t.Fatalf("failed to halt CPU: %v", err)
	}

	if t.Failed() {
		return
	}

	{
		const base = 0x80000000

		expected := []uint32{
			0x00000013, // nop

			// start function call code
			0x00000463, // beq x0, x0, 0x8
			0x80001000, // <location of function>
			0x00000e97, // auipc x29, 0x0
			0xffceae83, // lw x29, -4(x29)
			0x000e80e7, // jalr ra, 0(x29)
			0x00000013, // nop
			0x00000013, // nop
			0x00000013, // nop
			// end function call code

			0x005d8d93, // addi x27, x27, 0x5

			// start function call code
			0x00000463, // beq x0, x0, 0x8
			0x80001200, // <location of function>
			0x00000e97, // auipc x29, 0x0
			0xffceae83, // lw x29, -4(x29)
			0x000e80e7, // jalr ra, 0(x29)
			0x00000013, // nop
			0x00000013, // nop
			0x00000013, // nop
			// end function call code

			0x02000063, // beq x0, x0, 0x20 // skip over $C

			// start function call code
			0x00000463, // beq x0, x0, 0x8
			0x80001400, // <location of function>
			0x00000e97, // auipc x29, 0x0
			0xffceae83, // lw x29, -4(x29)
			0x000e80e7, // jalr ra, 0(x29)
			0x00000013, // nop
			0x00000013, // nop
			0x00000013, // nop
			// end function call code

			0x00000063, // beq x0, x0, 0x0
		}

		for i, want := range expected {
			addr := base + uint32(i*4)
			data, err := rt.GDB.ReadMem(addr, 4)
			if err != nil {
				t.Fatalf("failed to read from 0x%08x: %s", addr, err)
			}

			got := binary.LittleEndian.Uint32(data)

			if got != want {
				t.Errorf(
					"instruction %d @ 0x%08x = 0x%08x, want 0x%08x",
					i, addr, got, want,
				)

				continue
			}

			t.Logf("OK: instruction %d @ 0x%08x = 0x%08x",
				i, addr, got)
		}
	}

	if t.Failed() {
		return
	}

	{
		const base = 0x80001000

		expected := []uint32{
			0x40000d93, // addi x27 x0 0x400
			0x00000013, // nop
			0x00000013, // nop
			0x00008067, // ret
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

			t.Logf("OK: instruction %d @ 0x%08x = 0x%08x", i, addr, got)
		}
	}

	{
		const base = 0x80001200

		expected := []uint32{
			0x001d8d93, // addi x27 x27 0x1
			0x00008067, // ret
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

			t.Logf("OK: instruction %d @ 0x%08x = 0x%08x", i, addr, got)
		}
	}

	{
		const base = 0x80001400

		expected := []uint32{
			0x001d8d93, // addi x27 x27 0x1
			0x00008067, // ret
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

			t.Logf("OK: instruction %d @ 0x%08x = 0x%08x", i, addr, got)
		}
	}

	x27, err := rt.GDB.ReadReg(X27)
	if err != nil {
		t.Fatalf("failed to dump x27: %s", err)
	}

	expectedX27 := uint32(0x5) + uint32(0x400) + uint32(0x1)

	if x27 != expectedX27 {
		t.Fatalf("expected x27 to be 0x%08x but got 0x%08x", expectedX27, x27)
	}

	t.Logf("OK: x27 == 0x%08x", x27)
}
