package qemu

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
)

// QEMU wraps an exec.Cmd for running QEMU, establishing a new pseudo-tty for serial
// communication and shutting down properly when finished.
type QEMU struct {
	cmd *exec.Cmd

	kernelPath string

	stdout        bytes.Buffer
	stdoutScanner *bufio.Scanner

	stdin *io.PipeWriter

	pty *PTY
}

func (q *QEMU) stdoutReader() {
	for q.stdoutScanner.Scan() {
		q.stdout.Write(q.stdoutScanner.Bytes())
	}
}

// NewQEMU initialises a new QEMU command and allocates a ptty for serial communications
// but doesn't start QEMU itself
func NewQEMU(ctx context.Context, kernelPath string) (*QEMU, error) {
	pty, err := NewPTY()
	if err != nil {
		return nil, err
	}

	binaryPath := "/usr/bin/qemu-system-riscv32"

	qemuArgs := []string{
		"-nographic",
		"-serial",
		pty.SubFilename,
		"-s",
		"-S",
		"-M",
		"sifive_e",
		"-monitor",
		"stdio",
		"-kernel",
		kernelPath,
	}

	qemuStdoutReader, qemuStdoutWriter := io.Pipe()
	qemuStdinReader, qemuStdinWriter := io.Pipe()

	cmd := exec.CommandContext(ctx, binaryPath, qemuArgs...)
	cmd.Stdout = qemuStdoutWriter
	cmd.Stdin = qemuStdinReader

	qemu := &QEMU{
		cmd:           cmd,
		kernelPath:    kernelPath,
		stdoutScanner: bufio.NewScanner(qemuStdoutReader),
		stdin:         qemuStdinWriter,
		pty:           pty,
	}

	go qemu.stdoutReader()

	return qemu, nil
}

// Start is analogous to exec.Cmd.Start; begins the command begins reading stdout
func (q *QEMU) Start() error {
	return fmt.Errorf("NYI")
}

// Close attempts to shut down QEMU, first gracefully and then by force if required.
func (q *QEMU) Close() error {
	return fmt.Errorf("NYI")
}
