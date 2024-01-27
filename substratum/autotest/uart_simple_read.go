package autotest

import (
	"context"
	"fmt"
	"io"
)

// ProcessUARTSimpleReadBreak verifies the execution of the given GDB target and checks that it handles input as expected
// for the "uart" bedrock bare-metal program.
func ProcessUARTSimpleReadBreak(ctx context.Context, state *State) error {
	err := state.GDBConn.StepOnce()
	if err != nil {
		return err
	}

	target := uint32(0x20400060)

	err = state.GDBConn.AdvanceToBreak(target)
	if err != nil {
		return err
	}

	state.Logger.InfoContext(ctx, "finished advancing, fetching registers")

	frame, err := state.GDBConn.FetchRegisterFrame()
	if err != nil {
		return err
	}

	frame.Dump(state.Logger)

	//serialReader := io.LimitReader(state.SerialConn, 1)

	//data, err := io.ReadAll(serialReader)
	//if err != nil {
	//	return err
	//}

	//expectedOutut := "5"

	//if string(data) != expectedOutut {
	//	return fmt.Errorf("unexpected output from serial port\nwanted: %s\n   got: %s", expectedOutut, string(data))
	//}

	return nil
}

// ProcessUARTSimpleRead verifies the execution of the given GDB target and checks that it handles input as expected
// for the "uart" bedrock bare-metal program.
func ProcessUARTSimpleRead(ctx context.Context, state *State) error {
	err := state.GDBConn.StepOnce()
	if err != nil {
		return err
	}

	start := uint32(0x20400000)

	state.Logger.InfoContext(ctx, fmt.Sprintf("advancing PC to 0x%8.8X", start))

	err = state.GDBConn.AdvancePC(start, 1000)
	if err != nil {
		return err
	}

	target := uint32(0x20400080)

	for {
		frame, err := state.GDBConn.FetchRegisterFrame()
		if err != nil {
			return err
		}

		if frame.PC == target {
			break
		}

		state.Logger.InfoContext(ctx, fmt.Sprintf("got PC 0x%8.8X, advancing", frame.PC))

		err = state.GDBConn.StepOnce()
		if err != nil {
			return err
		}
	}

	frame, err := state.GDBConn.FetchRegisterFrame()
	if err != nil {
		return err
	}

	frame.Dump(state.Logger)

	serialReader := io.LimitReader(state.SerialConn, 1)

	data, err := io.ReadAll(serialReader)
	if err != nil {
		return err
	}

	expectedOutut := "5"

	if string(data) != expectedOutut {
		return fmt.Errorf("unexpected output from serial port\nwanted: %s\n   got: %s", expectedOutut, string(data))
	}

	return nil
}
