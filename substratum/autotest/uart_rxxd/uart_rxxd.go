package uart_rxxd

import (
	"fmt"
	"reflect"

	"github.com/sgtcodfish/substratum/autotest"

	"github.com/sgtcodfish/substratum"
)

// ProcessUARTRxxdBasic verifies the execution of the given GDB target and checks that it handles input as expected
// for the "uart-rxxd" bedrock bare-metal program.
// The "basic" test has only basic UART input, whose presence is checked in memory after running the whole program
func ProcessUARTRxxdBasic(state *autotest.State) error {
	err := checkInitialization(state)
	if err != nil {
		return err
	}

	msg := "13 00 00 00"
	bytesWritten, err := state.SerialConn.Write([]byte(msg))

	if err != nil {
		return err
	}

	if bytesWritten != len(msg) {
		return fmt.Errorf("couldn't write whole test message '%s', only wrote %d/%d bytes", msg, bytesWritten, len(msg))
	}

	state.Logger.Printf("successfully wrote %d bytes over UART", bytesWritten)

	err = state.GdbConn.AdvancePC(0x204000cc, 10)
	if err != nil {
		return err
	}

	frame, err := state.GdbConn.FetchRegisterFrame()
	if err != nil {
		return err
	}

	if frame.A0 != uint32('1') {
		return fmt.Errorf("a0 == 0x%8.8x but expected 0x%8.8x", frame.A0, uint32('1'))
	}

	state.Logger.Printf("a0 was set correctly after a read from UART")

	frame.Dump(state.Logger)

	return nil
}

// ProcessUARTRxxdFull verifies the execution of the given GDB target and checks that it handles input as expected
// for the "uart-rxxd" bedrock bare-metal program.
// The "full" test includes comments, invalid characters and multiple lines of text in the UART input
func ProcessUARTRxxdFull(state *autotest.State) error {
	return fmt.Errorf("uart-rxxd-full NYI: %v", state.Logger.Flags())
}

// checkInitialization advances execution until UART input is read and asserts that the registers
// were initialized as expected.
func checkInitialization(state *autotest.State) error {
	err := state.GdbConn.AdvancePC(0x204000b0, 100)
	if err != nil {
		return err
	}

	frame, err := state.GdbConn.FetchRegisterFrame()
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

	state.Logger.Printf("registers initialised as expected")
	return nil
}
