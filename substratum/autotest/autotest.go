package autotest

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/tarm/serial"

	"github.com/sgtcodfish/substratum"
	"github.com/sgtcodfish/substratum/cmd/util"
	"github.com/sgtcodfish/substratum/qemu"
)

const maxSerialWrites = 100

// TestFunc is a utility type for referring to any function which implements a substratum autotest
type TestFunc func(ctx context.Context, state *State) error

// State holds state which is required to run tests. Must be initialised prior to use.
type State struct {
	// Logger is the default logger to use for output during tests
	Logger *slog.Logger

	// GDBConn holds the connection to GDB, which will be manipulated throughout the test
	GDBConn *substratum.GDBConnection

	// QEMU holds connection details for manipulating QEMU
	QEMU *qemu.QEMU

	// SerialConn holds a persistent open connection over UART
	SerialConn *serial.Port
}

var _ io.Closer = (*State)(nil)

// NewState returns a new State with the given options, and opens (and holds open) a serial connecton based on serialOptions
func NewState(ctx context.Context, qemuPath string, gdbPath string, gdbPort string, kernelPath string) (*State, error) {
	logger := util.Logger(ctx)
	logger.InfoContext(ctx, "creating new QEMU instance and PTY")

	qemu, err := qemu.NewQEMU(ctx, qemuPath, kernelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create QEMU instance: %w", err)
	}

	err = qemu.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start QEMU instance: %w", err)
	}

	serialDevice := qemu.SerialDevice()
	logger.InfoContext(ctx, "initialised QEMU", "serialDevice", serialDevice)

	serialOptions := &serial.Config{
		Name:        serialDevice,
		Baud:        115200,
		ReadTimeout: 2 * time.Second,
	}

	serialConn, err := serial.OpenPort(serialOptions)
	if err != nil {
		_ = qemu.Close()
		return nil, err
	}

	gdbConn, err := substratum.NewGDBConnection(gdbPath, gdbPort)
	if err != nil {
		_ = qemu.Close()
		_ = serialConn.Close()
		return nil, err
	}

	logger.InfoContext(ctx, "connected with GDB", "port", gdbPort)

	return &State{
		Logger:     logger,
		GDBConn:    gdbConn,
		QEMU:       qemu,
		SerialConn: serialConn,
	}, nil
}

// Run invokes the given test function using this state
func (s *State) Run(ctx context.Context, testFunc TestFunc) error {
	err := testFunc(ctx, s)
	if err != nil {
		return err
	}

	return nil
}

// Close terminates any open connections held by the State, gracefully if possible
func (s *State) Close() error {
	var closeErrors []string

	err := s.SerialConn.Close()
	if err != nil {
		closeErrors = append(closeErrors, fmt.Sprintf("failed to close serial connection cleanly: %s", err.Error()))
	}

	err = s.GDBConn.Close()
	if err != nil {
		closeErrors = append(closeErrors, fmt.Sprintf("failed to terminate GDB and close connection cleanly: %s", err.Error()))
	}

	err = s.QEMU.Close()
	if err != nil {
		closeErrors = append(closeErrors, fmt.Sprintf("failed to terminate QEMU: %s", err.Error()))
	}

	if len(closeErrors) > 0 {
		return fmt.Errorf("failed to shutdown test state cleanly: %s", strings.Join(closeErrors, " | "))
	}

	return nil
}

// ReadAllSerial attempts to read everything available from the serial connection
func (s *State) ReadAllSerial() ([]byte, error) {
	return io.ReadAll(s.SerialConn)
}

// SendSerial attempts to fully write the given data to the serial port corresponding to this State.
func (s *State) SendSerial(data []byte) error {
	bytesRemaining := len(data)
	for i := 0; i < maxSerialWrites; i++ {
		bytesWritten, err := s.SerialConn.Write(data)
		if err != nil {
			return err
		}

		bytesRemaining -= bytesWritten

		if bytesRemaining <= 0 {
			//s.Logger.Printf("wrote %d bytes over UART", len(data))
			return nil
		}
	}

	return fmt.Errorf("couldn't write all %d bytes of message (wrote %d bytes)", len(data), len(data)-bytesRemaining)
}
