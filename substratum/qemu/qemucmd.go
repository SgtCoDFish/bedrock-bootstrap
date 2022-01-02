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

	stdin *io.WriteCloser

	pty *PTY
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

	cmd := exec.CommandContext(ctx, binaryPath, qemuArgs...)

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	qemu := &QEMU{
		cmd:           cmd,
		kernelPath:    kernelPath,
		stdoutScanner: bufio.NewScanner(stdoutPipe),
		stdin:         &stdinPipe,
		pty:           pty,
	}

	go qemu.stdoutReader()

	return qemu, nil
}

// SerialDevice returns the name of the file for the subordinate console
func (q *QEMU) SerialDevice() string {
	return q.pty.SubFilename
}

func (q *QEMU) stdoutReader() {
	for q.stdoutScanner.Scan() {
		q.stdout.Write(q.stdoutScanner.Bytes())
	}
}

// Start is analogous to exec.Cmd.Start; begins the command begins reading stdout
func (q *QEMU) Start() error {
	return fmt.Errorf("NYI")
}

// Close attempts to shut down QEMU, first gracefully and then by force if required.
func (q *QEMU) Close() error {
	return fmt.Errorf("NYI")
}
