package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sgtcodfish/substratum"
	"github.com/sgtcodfish/substratum/autotest/uart_rxxd"
)

func processAutotest(flags *flag.FlagSet, logger *log.Logger) error {
	testMap := map[string]func(connection *substratum.GdbConnection) error{
		"uart-rxxd-basic": uart_rxxd.ProcessUARTRxxdBasic,
		"uart-rxxd-full":  uart_rxxd.ProcessUARTRxxdFull,
	}

	allTests := make([]string, 0)
	for k := range testMap {
		allTests = append(allTests, k)
	}

	if flags.NArg() == 0 {
		return fmt.Errorf("missing required argument <test-name>: must be one of: %s", strings.Join(allTests, " | "))
	}

	testName := strings.ToLower(flags.Arg(0))

	testFn, ok := testMap[testName]
	if !ok {
		return fmt.Errorf("unknown test name '%s'", testName)
	}

	gdbPath := os.Getenv("RISCV_PREFIX") + "gdb"
	remoteTarget := ":1234"

	conn, err := substratum.NewGdbConnection(logger, gdbPath, remoteTarget)
	if err != nil {
		return err
	}

	err = testFn(conn)
	if err != nil {
		return err
	}

	logger.Printf("'%s' test ran successfully", testName)
	return nil
}
