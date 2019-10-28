package autotest

import (
	"fmt"
	"io"
	"log"

	"github.com/jacobsa/go-serial/serial"
	"github.com/sgtcodfish/substratum"
)

const maxSerialWrites = 100

// State holds state which is required to run tests. Usually the state must be filled in prior to use.
type State struct {
	// Logger is the default logger to use for output during tests
	Logger *log.Logger

	// GdbConn holds the connection to GDB, which will be manipulated throughout the test
	GdbConn *substratum.GdbConnection

	// SerialConn holds a persistent open connection over UART
	SerialConn io.ReadWriteCloser

	// SerialOptions can be passed when opening a serial connection to send data
	SerialOptions serial.OpenOptions
}

// NewState returns a new State with the given options, and opens (and holds open) a serial connecton based on serialOptions
func NewState(logger *log.Logger, gdbConn *substratum.GdbConnection, serialOptions serial.OpenOptions) (*State, error) {
	serialConn, err := serial.Open(serialOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to open serial connection to '%s': %w", serialOptions.PortName, err)
	}

	return &State{
		Logger:        logger,
		GdbConn:       gdbConn,
		SerialConn:    serialConn,
		SerialOptions: serialOptions,
	}, nil
}

// SendSerial takes the given byte slice of data and attempts to write it to the serial port whose options are given in
// State.SerialOptions.
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
