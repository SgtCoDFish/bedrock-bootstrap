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

func getStringFlag(flags *flag.FlagSet, name string) (string, error) {
	flagVar := flags.Lookup(name)
	if flagVar == nil {
		return "", fmt.Errorf("missing required flag '%s'", name)
	}

	flagVal := flagVar.Value.String()
	if flagVal == "" {
		return "", fmt.Errorf("missing flag value for '%s'", name)
	}

	return flagVal, nil
}

func processAutotest(flags *flag.FlagSet, logger *log.Logger) error {
	testMap := map[string]func(state *autotest.State) error{
		"uart-rxxd-basic": autotest.ProcessUARTRxxdBasic,
		"uart-rxxd-full":  autotest.ProcessUARTRxxdFull,
	}

	allTests := make([]string, 0)
	for k := range testMap {
		allTests = append(allTests, k)
	}

	serialDevice, err := getStringFlag(flags, "serial")
	if err != nil {
		return err
	}

	command, err := getStringFlag(flags, "test-name")
	if err != nil {
		return err
	}

	testName := strings.ToLower(command)

	testFn, ok := testMap[testName]
	if !ok {
		return fmt.Errorf("unknown test name '%s'; must be one of %s", testName, strings.Join(allTests, " | "))
	}

	gdbPath := os.Getenv("RISCV_PREFIX") + "gdb"
	remoteTarget := ":1234"

	conn, err := substratum.NewGdbConnection(gdbPath, remoteTarget)
	if err != nil {
		return fmt.Errorf("failed to open GDB connection to '%s' using '%s': %w", remoteTarget, gdbPath, err)
	}

	serialOptions := serial.OpenOptions{
		PortName:        serialDevice,
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		ParityMode:      serial.PARITY_NONE,
		MinimumReadSize: 1,
	}

	testState, err := autotest.NewState(logger, conn, serialOptions)
	if err != nil {
		return fmt.Errorf("failed to create new autotest state: %w", err)
	}

	err = testFn(testState)
	if err != nil {
		return err
	}

	logger.Printf("'%s' test ran successfully", testName)
	return nil
}
