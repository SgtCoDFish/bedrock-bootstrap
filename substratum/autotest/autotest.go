package autotest

import (
	"io"
	"log"

	"github.com/sgtcodfish/substratum"
)

// State holds state which is required to run tests. Usually the state must be filled in prior to use.
type State struct {
	// Logger is the default logger to use for output during tests
	Logger *log.Logger

	// GdbConn holds the connection to GDB, which will be manipulated throughout the test
	GdbConn *substratum.GdbConnection

	// SerialConn allows the test to send over UART
	SerialConn io.ReadWriteCloser
}
