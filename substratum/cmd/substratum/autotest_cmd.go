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
)

func processAutotest(flags *flag.FlagSet, logger *log.Logger) error {
	testMap := map[string]func(state *autotest.State) error{
		"uart-rxxd-basic": autotest.ProcessUARTRxxdBasic,
		"uart-rxxd-full":  autotest.ProcessUARTRxxdFull,
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

	serialOptions := serial.OpenOptions{
		PortName:        serialDevice.Value.String(),
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        2,
		ParityMode:      serial.PARITY_NONE,
		MinimumReadSize: 1,
	}

	testState := &autotest.State{
		Logger:        logger,
		GdbConn:       conn,
		SerialOptions: serialOptions,
	}

	err = testFn(testState)
	if err != nil {
		return err
	}

	logger.Printf("'%s' test ran successfully", testName)
	return nil
}
