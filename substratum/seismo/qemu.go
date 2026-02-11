package seismo

import (
	"context"
	"fmt"
	"os/exec"
)

type QEMUMachine string

const (
	SiFiveE QEMUMachine = "sifive_e"
	Virt    QEMUMachine = "virt"

	serialPort = 4444
	gdbPort    = 3333
)

type QEMU struct {
	Cmd        *exec.Cmd
	SerialPort int
	GDBPort    int
}

func StartQEMU(ctx context.Context, machine QEMUMachine, kernelPath string) (*QEMU, error) {
	args := []string{
		"-machine", string(SiFiveE),
		"-nographic",
		"-bios", "none",
		"-kernel", kernelPath,
		"-serial", fmt.Sprintf("tcp::%d,server", serialPort),
		"-gdb", fmt.Sprintf("tcp::%d", gdbPort),
		"-S",
	}

	cmd := exec.CommandContext(ctx, "qemu-system-riscv32", args...)
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &QEMU{
		Cmd:        cmd,
		SerialPort: serialPort,
		GDBPort:    gdbPort,
	}, nil
}

func (q *QEMU) Stop() error {
	if q.Cmd.Process != nil {
		_ = q.Cmd.Process.Kill()
	}
	return q.Cmd.Wait()
}
