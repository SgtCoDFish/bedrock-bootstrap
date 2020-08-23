package autotest

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sgtcodfish/substratum"
)

const (
	memoryStart = 0x80000000
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

	initialMemoryLoc := uint32(0x20400000)

	word, err := state.GdbConn.ReadMemoryWord(initialMemoryLoc)
	if err != nil {
		return err
	}

	fmt.Printf("word at 0x%8.8X: %s\n", initialMemoryLoc, word)

	for i := 0; i < len(msg); i++ {
		err = state.GdbConn.AdvancePC(0x204000cc, 1000)
		if err != nil {
			return err
		}

		a0, err := state.GdbConn.FetchRegister("a0")
		if err != nil {
			return err
		}

		expected := uint32(msg[i])
		if a0 != expected {
			return fmt.Errorf("a0 == 0x%8.8X but expected 0x%8.8X", a0, expected)
		}

		state.Logger.Printf("a0 was set correctly to 0x%8.8X after a read from UART", msg[i])

		err = state.GdbConn.AdvancePC(0x204000b0, 1000)
		if err != nil {
			return err
		}
	}

	frame, err := state.GdbConn.FetchRegisterFrame()
	if err != nil {
		return err
	}

	frame.Dump(state.Logger)

	for i := uint32(memoryStart); i < memoryStart+0x10; i += 4 {
		word, err := state.GdbConn.ReadMemoryWord(i)
		if err != nil {
			return err
		}

		fmt.Printf("word at 0x%8.8X: %s\n", i, word)

		// we only want to write our single word at the initial memory location
		// and we don't want to touch any of the surrounding memory
		if i == memoryStart {
			if word != "13000000" {
				return fmt.Errorf("wanted memory at 0x%8.8X == 0x13000000 but got %s", i, word)
			}
		} else {
			if word != "00000000" {
				return fmt.Errorf("wanted memory at 0x%8.8X == 0x00000000 but got %s", i, word)
			}
		}
	}

	return nil
}

// ProcessUARTRxxdComment verifies the execution of the given GDB target and checks that it handles input as expected
// for the "uart-rxxd" bedrock bare-metal program.
// The "comment" test has only basic UART input, whose presence is checked in memory after running the whole program, but
// crucially also includes a leading comment which should be ignored
func ProcessUARTRxxdComment(state *State) error {
	_ = state.GdbConn.StepOnce()

	err := checkInitialization(state)
	if err != nil {
		return err
	}

	msg := []byte("# comment\n13000000")
	err = state.SendSerial(msg)
	if err != nil {
		return err
	}

	initialMemoryLoc := uint32(0x20400000)

	word, err := state.GdbConn.ReadMemoryWord(initialMemoryLoc)
	if err != nil {
		return err
	}

	fmt.Printf("word at 0x%8.8X: %s\n", initialMemoryLoc, word)

	for i := 0; i < len(msg); i++ {
		err = state.GdbConn.AdvancePC(0x204000cc, 1000)
		if err != nil {
			return err
		}

		a0, err := state.GdbConn.FetchRegister("a0")
		if err != nil {
			return err
		}

		expected := uint32(msg[i])
		if a0 != expected {
			return fmt.Errorf("a0 == 0x%8.8X but expected 0x%8.8X", a0, expected)
		}

		state.Logger.Printf("a0 was set correctly to 0x%8.8X after a read from UART", msg[i])

		err = state.GdbConn.AdvancePC(0x204000b0, 1000)
		if err != nil {
			return err
		}
	}

	frame, err := state.GdbConn.FetchRegisterFrame()
	if err != nil {
		return err
	}

	frame.Dump(state.Logger)

	for i := uint32(memoryStart); i < memoryStart+0x10; i += 4 {
		word, err := state.GdbConn.ReadMemoryWord(i)
		if err != nil {
			return err
		}

		fmt.Printf("word at 0x%8.8X: %s\n", i, word)

		// we only want to write our single word at the initial memory location
		// and we don't want to touch any of the surrounding memory
		if i == memoryStart {
			if word != "13000000" {
				return fmt.Errorf("wanted memory at 0x%8.8X == 0x13000000 but got %s", i, word)
			}
		} else {
			if word != "00000000" {
				return fmt.Errorf("wanted memory at 0x%8.8X == 0x00000000 but got %s", i, word)
			}
		}
	}

	return nil
}

// ProcessUARTRxxdFull verifies the execution of the given GDB target and checks that it handles input as expected
// for the "uart-rxxd" bedrock bare-metal program.
// The "full" test includes comments, invalid characters and multiple lines of text in the UART input
func ProcessUARTRxxdFull(state *State) error {
	_ = state.GdbConn.StepOnce()

	err := checkInitialization(state)
	if err != nil {
		return err
	}

	msg := []byte("13000000 # nop dummy command\n13 00 00 00 # endless screaming\n\n    13 0f 30 12 # addi x30 x00 0x123 #\nJ")
	err = state.SendSerial(msg)
	if err != nil {
		return err
	}

	initialMemoryLoc := uint32(0x20400000)

	word, err := state.GdbConn.ReadMemoryWord(initialMemoryLoc)
	if err != nil {
		return err
	}

	fmt.Printf("word at 0x%8.8X: %s\n", initialMemoryLoc, word)

	for i := 0; i < len(msg); i++ {
		// advance to 2040_00CC which is just after UART has been read from
		err = state.GdbConn.AdvancePC(0x204000cc, 1000)
		if err != nil {
			return err
		}

		a0, err := state.GdbConn.FetchRegister("a0")
		if err != nil {
			return err
		}

		expected := uint32(msg[i])
		if a0 != expected {
			return fmt.Errorf("a0 == 0x%8.8X but expected 0x%8.8X", a0, expected)
		}

		state.Logger.Printf("a0 was set correctly to 0x%8.8X after a read from UART", msg[i])

		if strings.ToLower(string(msg[i])) == "j" {
			fmt.Println("found a J in input")
			break
		}

		// advance to 2040_00B0 which is the start of reading from UART
		err = state.GdbConn.AdvancePC(0x204000b0, 1000)
		if err != nil {
			return err
		}
	}

	err = state.GdbConn.WalkPC(0x20400160, 50)
	if err != nil {
		return err
	}

	frame, err := state.GdbConn.FetchRegisterFrame()
	if err != nil {
		return err
	}

	frame.Dump(state.Logger)

	for i := uint32(memoryStart); i < memoryStart+0x10; i += 4 {
		word, err := state.GdbConn.ReadMemoryWord(i)
		if err != nil {
			return err
		}

		fmt.Printf("word at 0x%8.8X: %s\n", i, word)

		// we only want to write our single word at the initial memory location
		// and we don't want to touch any of the surrounding memory
		if i == memoryStart+0x00 {
			if word != "13000000" {
				return fmt.Errorf("wanted memory at 0x%8.8X == 0x13000000 but got %s", i, word)
			}
		} else if i == memoryStart+0x04 {
			if word != "13000000" {
				return fmt.Errorf("wanted memory at 0x%8.8X == 0x13000000 but got %s", i, word)
			}
		} else if i == memoryStart+0x08 {
			if word != "130f3012" {
				return fmt.Errorf("wanted memory at 0x%8.8X == 0x130f3012 but got %s", i, word)
			}
		} else if i == memoryStart+0x0C {
			if word != "63000000" {
				return fmt.Errorf("wanted memory at 0x%8.8X == 0x63000000 but got %s", i, word)
			}
		} else {
			if word != "00000000" {
				return fmt.Errorf("wanted memory at 0x%8.8X == 0x00000000 but got %s", i, word)
			}
		}
	}

	err = state.GdbConn.WalkPC(memoryStart+0xC, 50)
	if err != nil {
		return err
	}

	err = state.GdbConn.StepOnce()
	if err != nil {
		return err
	}

	frame, err = state.GdbConn.FetchRegisterFrame()
	if err != nil {
		return err
	}

	frame.Dump(state.Logger)

	if frame.T5 != 0x123 {
		return fmt.Errorf("wanted x30/t5 to be 0x123 but got %8.8X", frame.T5)
	}

	return nil
}

// checkInitialization advances execution until UART input is read and asserts that the registers
// were initialized as expected.
func checkInitialization(state *State) error {
	err := state.GdbConn.AdvancePC(0x204000b0, 1000)
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
		A4: memoryStart,
		A5: 0x10013000,
		A6: 0x10013004,
		A7: 0xa,
		S2: 0x20,
		S5: memoryStart,
		PC: 0x204000b0,
	}

	if !reflect.DeepEqual(expectedInitialFrame, frame) {
		return fmt.Errorf("registers were not initialised correctly.\ngot : %#v\nwant: %#v", frame, expectedInitialFrame)
	}

	state.Logger.Printf("registers initialised as expected")
	return nil
}
