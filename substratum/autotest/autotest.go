package autotest

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/jacobsa/go-serial/serial"
	"github.com/sgtcodfish/substratum"
)

const maxSerialWrites = 100

// State holds state which is required to run tests. Must be initialised prior to use.
type State struct {
	// Logger is the default logger to use for output during tests
	Logger *log.Logger

	// GDBConn holds the connection to GDB, which will be manipulated throughout the test
	GDBConn *substratum.GDBConnection

	// SerialConn holds a persistent open connection over UART
	serialConn io.ReadWriteCloser
}

var _ io.Closer = (*State)(nil)

// NewState returns a new State with the given options, and opens (and holds open) a serial connecton based on serialOptions
func NewState(logger *log.Logger, gdbConn *substratum.GDBConnection, serialOptions serial.OpenOptions) (*State, error) {
	serialConn, err := serial.Open(serialOptions)
	if err != nil {
		return nil, err
	}

	return &State{
		Logger:     logger,
		GDBConn:    gdbConn,
		serialConn: serialConn,
	}, nil
}

// Close terminates any open connections held by the State, gracefully if possible
func (s *State) Close() error {
	var closeErrors []string

	err := s.serialConn.Close()
	if err != nil {
		closeErrors = append(closeErrors, fmt.Sprintf("failed to close serial connection cleanly: %s", err.Error()))
	}

	err = s.GDBConn.Close()
	if err != nil {
		closeErrors = append(closeErrors, fmt.Sprintf("failed to terminate GDB and close connection cleanly: %s", err.Error()))
	}

	if len(closeErrors) > 0 {
		return fmt.Errorf("failed to shutdown test state cleanly: %s", strings.Join(closeErrors, " | "))
	}

	return nil
}

// SendSerial attempts to fully write the given data to the serial port corresponding to this State.
func (s *State) SendSerial(data []byte) error {
	bytesRemaining := len(data)
	for i := 0; i < maxSerialWrites; i++ {
		bytesWritten, err := s.serialConn.Write(data)
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
