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
		"uart-rxxd-basic":   autotest.ProcessUARTRxxdBasic,
		"uart-rxxd-comment": autotest.ProcessUARTRxxdComment,
		"uart-rxxd-full":    autotest.ProcessUARTRxxdFull,
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

	logger.Printf("processing autotest for '%s'", testName)
	logger.Printf("starting GDB")

	gdbPath := os.Getenv("RISCV_PREFIX") + "gdb"
	remoteTarget := ":1234"

	conn, err := substratum.NewGDBConnection(gdbPath, remoteTarget)
	if err != nil {
		return err
	}

	serialOptions := serial.OpenOptions{
		PortName:        serialDevice.Value.String(),
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		ParityMode:      serial.PARITY_NONE,
		MinimumReadSize: 1,
	}

	logger.Printf("starting serial connection")

	testState, err := autotest.NewState(logger, conn, serialOptions)
	if err != nil {
		return err
	}

	err = testFn(testState)
	if err != nil {
		return err
	}

	logger.Printf("'%s' test ran successfully", testName)
	return nil
}
