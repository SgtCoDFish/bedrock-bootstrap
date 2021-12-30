package autotest

import (
	"fmt"
	"io"
	"log"

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
	SerialConn io.ReadWriteCloser

	// SerialOptions can be passed when opening a serial connection to send data
	SerialOptions serial.OpenOptions
}

// NewState returns a new State with the given options, and opens (and holds open) a serial connecton based on serialOptions
func NewState(logger *log.Logger, gdbConn *substratum.GDBConnection, serialOptions serial.OpenOptions) (*State, error) {
	serialConn, err := serial.Open(serialOptions)
	if err != nil {
		return nil, err
	}

	return &State{
		Logger:        logger,
		GDBConn:       gdbConn,
		SerialConn:    serialConn,
		SerialOptions: serialOptions,
	}, nil
}

// SendSerial attempts to write the given data to the serial port corresponding to this State.
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
