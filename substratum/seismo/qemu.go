package seismo

import (
	"fmt"
	"os/exec"
)

type QEMU struct {
	Cmd        *exec.Cmd
	SerialPort int
	GDBPort    int
}

func StartQEMU(kernel string, serialPort, gdbPort int) (*QEMU, error) {
	args := []string{
		"-machine", "sifive_e",
		"-nographic",
		"-kernel", kernel,
		"-serial", fmt.Sprintf("tcp::%d,server", serialPort),
		"-gdb", fmt.Sprintf("tcp::%d", gdbPort),
		"-S",
	}

	cmd := exec.Command("qemu-system-riscv32", args...)
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
