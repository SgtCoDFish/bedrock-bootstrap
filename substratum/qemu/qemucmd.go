package qemu

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

// QEMU wraps an exec.Cmd for running QEMU, establishing a new pseudo-tty for serial
// communication and shutting down properly when finished.
type QEMU struct {
	cmd *exec.Cmd

	kernelPath string

	logger        *log.Logger
	stdoutScanner *bufio.Scanner

	stdin io.WriteCloser

	pty *PTY
}

// NewQEMU initialises a new QEMU command and allocates a ptty for serial communications
// but doesn't start QEMU itself
func NewQEMU(ctx context.Context, qemuPath string, kernelPath string) (*QEMU, error) {
	logger := log.New(os.Stderr, "qemu: ", 0)

	pty, err := NewPTY()
	if err != nil {
		return nil, err
	}

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

	cmd := exec.CommandContext(ctx, qemuPath, qemuArgs...)

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
		logger:        logger,
		stdoutScanner: bufio.NewScanner(stdoutPipe),
		stdin:         stdinPipe,
		pty:           pty,
	}

	return qemu, nil
}

// SerialDevice returns the name of the file for the subordinate console
func (q *QEMU) SerialDevice() string {
	return q.pty.SubFilename
}

func (q *QEMU) stdoutReader() {
	for q.stdoutScanner.Scan() {
		q.logger.Println(q.stdoutScanner.Text())
	}
}

// Start is analogous to exec.Cmd.Start; begins the command begins reading stdout
func (q *QEMU) Start() error {
	err := q.cmd.Start()
	if err != nil {
		return err
	}

	go q.stdoutReader()

	return nil
}

// Close attempts to shut down QEMU, first gracefully and then by force if required.
func (q *QEMU) Close() error {
	var errs []string

	_, err := q.stdin.Write([]byte("\nquit\n"))
	if err != nil {
		errs = append(errs, err.Error())
	}

	err = q.stdin.Close()
	if err != nil {
		errs = append(errs, err.Error())
	}

	err = q.cmd.Wait()
	if err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to close QEMU instance: %s", strings.Join(errs, " | "))
	}

	return nil
}
