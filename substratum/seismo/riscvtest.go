package seismo

import (
	"testing"
	"time"
)

type RISCVTest struct {
	QEMU *QEMU
	UART *UART
	GDB  *GDB
}

func NewRISCVTest(t *testing.T, kernelPath string) *RISCVTest {
	ctx := t.Context()

	qemu, err := StartQEMU(ctx, SiFiveE, kernelPath)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("started QEMU")

	t.Cleanup(func() {
		qemu.Stop()
	})

	uart, err := ConnectUART(ctx, serialPort, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("connected UART")

	t.Cleanup(func() {
		uart.Close()
	})

	gdb, err := ConnectGDB(ctx, gdbPort)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("connected GDB")

	t.Cleanup(func() {
		gdb.Close()
	})

	return &RISCVTest{
		QEMU: qemu,
		UART: uart,
		GDB:  gdb,
	}
}
