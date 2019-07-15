package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/sgtcodfish/substratum"
)

func processAutotest(logger *log.Logger) error {
	gdbPath := os.Getenv("RISCV_PREFIX") + "gdb"
	remoteTarget := ":1234"

	conn, err := substratum.NewGdbConnection(logger, gdbPath, remoteTarget)
	if err != nil {
		return err
	}

	for {
		pcReg, err := conn.FetchPC()
		if err != nil {
			return err
		}

		if pcReg == 0x204000b0 {
			break
		}

		_, err = conn.Conn.CheckedSend("exec-step-instruction")
		if err != nil {
			return err
		}
	}

	frame, err := conn.FetchRegisterFrame()
	if err != nil {
		return err
	}

	frame.Dump(logger)

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

	logger.Printf("registers initialised as expected")

	return nil
}
