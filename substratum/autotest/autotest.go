package autotest

import (
	"fmt"
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

	// SerialOptions can be passed when opening a serial connection to send data
	SerialOptions serial.OpenOptions
}

// SendSerial takes the given byte slice of data and attempts to write it to the serial port whose options are given in
// State.SerialOptions.
func (s *State) SendSerial(data []byte) error {
	conn, err := serial.Open(s.SerialOptions)
	if err != nil {
		return err
	}

	bytesRemaining := len(data)
	for i := 0; i < maxSerialWrites; i++ {
		bytesWritten, err := conn.Write(data)
		if err != nil {
			return err
		}

		bytesRemaining -= bytesWritten
		if bytesRemaining <= 0 {
			err = conn.Close()
			if err != nil {
				s.Logger.Printf("got an error closing serial connection: %v", err)
			}

			s.Logger.Printf("wrote %d bytes over UART", len(data))

			return nil
		}
	}

	return fmt.Errorf("couldn't write all %d bytes of message (wrote %d bytes)", len(data), len(data)-bytesRemaining)
}
