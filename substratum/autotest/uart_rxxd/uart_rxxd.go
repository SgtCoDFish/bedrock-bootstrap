package uart_rxxd

import (
	"fmt"
	"reflect"

	"github.com/sgtcodfish/substratum"
)

// ProcessUARTRxxdBasic verifies the execution of the given GDB target and checks that it handles input as expected
// for the "uart-rxxd" bedrock bare-metal program.
// The "basic" test has only basic UART input, whose presence is checked in memory after running the whole program
func ProcessUARTRxxdBasic(gdbConn *substratum.GdbConnection) error {
	return checkInitialization(gdbConn)
}

// ProcessUARTRxxdFull verifies the execution of the given GDB target and checks that it handles input as expected
// for the "uart-rxxd" bedrock bare-metal program.
// The "full" test includes comments, invalid characters and multiple lines of text in the UART input
func ProcessUARTRxxdFull(gdbConn *substratum.GdbConnection) error {
	return fmt.Errorf("uart-rxxd-full NYI")
}

// checkInitialization advances execution until UART input is read and asserts that the registers
// were initialized as expected.
func checkInitialization(gdbConn *substratum.GdbConnection) error {
	err := gdbConn.AdvancePC(0x204000b0, 100)
	if err != nil {
		return err
	}

	frame, err := gdbConn.FetchRegisterFrame()
	if err != nil {
		return err
	}

	// frame.Dump(gdbConn.Logger)

	expectedInitialFrame := substratum.GDBRegisterFrame{
		T2: 0x4,
		A2: 0x80000000,
		A4: 0x80001000,
		A5: 0x10013000,
		A6: 0x10013004,
		A7: 0xa,
		S2: 0x20,
		PC: 0x204000b0,
	}

	if !reflect.DeepEqual(expectedInitialFrame, frame) {
		return fmt.Errorf("registers were not initialised correctly.\ngot : %#v\nwant: %#v", frame, expectedInitialFrame)
	}

	gdbConn.Logger.Printf("registers initialised as expected")

	return nil
}
