package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sgtcodfish/substratum/autotest"

	"github.com/jacobsa/go-serial/serial"

	"github.com/sgtcodfish/substratum"
	"github.com/sgtcodfish/substratum/autotest/uart_rxxd"
)

func processAutotest(flags *flag.FlagSet, logger *log.Logger) error {
	testMap := map[string]func(state *autotest.State) error{
		"uart-rxxd-basic": uart_rxxd.ProcessUARTRxxdBasic,
		"uart-rxxd-full":  uart_rxxd.ProcessUARTRxxdFull,
	}

	allTests := make([]string, 0)
	for k := range testMap {
		allTests = append(allTests, k)
	}

	serialDevice := flags.Lookup("serial")
	if serialDevice == nil {
		return fmt.Errorf("missing required flag %s", "serial")
	}

	command := flags.Lookup("test-name")
	if command == nil {
		return fmt.Errorf("missing required flag %s", "test-name")
	}

	testName := strings.ToLower(command.Value.String())

	testFn, ok := testMap[testName]
	if !ok {
		return fmt.Errorf("unknown test name '%s'; must be one of %s", testName, strings.Join(allTests, " | "))
	}

	gdbPath := os.Getenv("RISCV_PREFIX") + "gdb"
	remoteTarget := ":1234"

	conn, err := substratum.NewGdbConnection(gdbPath, remoteTarget)
	if err != nil {
		return err
	}

	serialPort, err := serial.Open(serial.OpenOptions{
		PortName:        serialDevice.Value.String(),
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		ParityMode:      serial.PARITY_NONE,
		MinimumReadSize: 4,
	})

	if err != nil {
		return err
	}

	defer func() { _ = serialPort.Close() }()

	testState := &autotest.State{
		Logger:     logger,
		GdbConn:    conn,
		SerialConn: serialPort,
	}

	err = testFn(testState)
	if err != nil {
		return err
	}

	logger.Printf("'%s' test ran successfully", testName)
	return nil
}
