package autotest

import (
	"fmt"
	"reflect"

	"github.com/sgtcodfish/substratum"
)

// ProcessUARTRxxdBasic verifies the execution of the given GDB target and checks that it handles input as expected
// for the "uart-rxxd" bedrock bare-metal program.
// The "basic" test has only basic UART input, whose presence is checked in memory after running the whole program
func ProcessUARTRxxdBasic(state *State) error {
	_ = state.GdbConn.StepOnce()

	err := checkInitialization(state)
	if err != nil {
		return err
	}

	msg := []byte("13000000")
	err = state.SendSerial(msg)
	if err != nil {
		return err
	}

	word, err := state.GdbConn.ReadMemoryWord(0x20400000)
	if err != nil {
		return err
	}
	fmt.Printf("word up: %s\n", word)

	for i := 0; i < len(msg); i++ {
		err = state.GdbConn.AdvancePC(0x204000cc, 200)
		if err != nil {
			return err
		}

		a0, err := state.GdbConn.FetchRegister("a0")
		if err != nil {
			return err
		}

		if a0 != uint32(msg[i]) {
			return fmt.Errorf("a0 == 0x%8.8x but expected 0x%8.8x", a0, uint32('1'))
		}

		state.Logger.Printf("a0 was set correctly to 0x%8.8x after a read from UART", msg[i])

		err = state.GdbConn.AdvancePC(0x204000b0, 200)
		if err != nil {
			return err
		}
	}

	frame, err := state.GdbConn.FetchRegisterFrame()
	if err != nil {
		return err
	}

	frame.Dump(state.Logger)

	return nil
}

// ProcessUARTRxxdFull verifies the execution of the given GDB target and checks that it handles input as expected
// for the "uart-rxxd" bedrock bare-metal program.
// The "full" test includes comments, invalid characters and multiple lines of text in the UART input
func ProcessUARTRxxdFull(state *State) error {
	return fmt.Errorf("uart-rxxd-full NYI: %v", state.Logger.Flags())
}

// checkInitialization advances execution until UART input is read and asserts that the registers
// were initialized as expected.
func checkInitialization(state *State) error {
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
