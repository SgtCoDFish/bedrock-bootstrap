package qemu

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

// QEMU wraps an exec.Cmd for running QEMU, establishing a new pseudo-tty for serial
// communication and shutting down properly when finished.
type QEMU struct {
	cmd *exec.Cmd

	kernelPath string

	logger        *log.Logger
	stdoutScanner *bufio.Scanner

	stdin io.WriteCloser

	SerialConn net.Conn
}

// NewQEMU initialises a new QEMU command and allocates a ptty for serial communications
// but doesn't start QEMU itself
func NewQEMU(ctx context.Context, qemuPath string, kernelPath string) (*QEMU, error) {
	logger := log.New(os.Stderr, "qemu: ", 0)

	qemuArgs := []string{
		"-nographic",
		"-serial",
		"tcp:127.0.0.1:4444,server,nowait",
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
		return nil, fmt.Errorf("failed to get stdin pipe for QEMU: %w", err)
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe for QEMU: %w", err)
	}

	qemu := &QEMU{
		cmd:           cmd,
		kernelPath:    kernelPath,
		logger:        logger,
		stdoutScanner: bufio.NewScanner(stdoutPipe),
		stdin:         stdinPipe,
	}

	return qemu, nil
}

func (q *QEMU) stdoutReader() {
	for q.stdoutScanner.Scan() {
		q.logger.Println(q.stdoutScanner.Text())
	}
}

// Start is analogous to exec.Cmd.Start; begins the command begins reading stdout
func (q *QEMU) Start(ctx context.Context) error {
	err := q.cmd.Start()
	if err != nil {
		return err
	}

	go q.stdoutReader()

	conn, err := waitForTCP(ctx, "127.0.0.1:4444", 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to dial serial TCP port: %w", err)
	}

	q.SerialConn = conn

	return nil
}

// waitForTCP retries until timeout
func waitForTCP(ctx context.Context, addr string, timeout time.Duration) (net.Conn, error) {
	deadline := time.Now().Add(timeout)

	dialer := net.Dialer{}

	for time.Now().Before(deadline) {
		conn, err := dialer.DialContext(ctx, "tcp", addr)
		if err == nil {
			return conn, nil
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil, fmt.Errorf("timed out waiting for connection to %s", addr)
}

// Close attempts to shut down QEMU, first gracefully and then by force if required.
func (q *QEMU) Close() error {
	var errs []string

	_, err := q.stdin.Write([]byte("\nquit\n"))
	if err != nil {
		errs = append(errs, err.Error())
	}

	err = q.SerialConn.Close()
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
