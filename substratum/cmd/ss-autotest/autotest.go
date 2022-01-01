package autotest

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sgtcodfish/substratum"
	"github.com/sgtcodfish/substratum/autotest"
)

var testMap = map[string]autotest.TestFunc{
	"uart-rxxd-basic":   autotest.ProcessUARTRxxdBasic,
	"uart-rxxd-comment": autotest.ProcessUARTRxxdComment,
	"uart-rxxd-full":    autotest.ProcessUARTRxxdFull,
}

// Invocation holds state for a given invocation of the ss-autotest command
type Invocation struct {
	gdbPath      string
	gdbPort      string
	serialDevice string
	testName     string
}

// ParseInvocation builds an invocation for a run of the the ss-autotest command, where name is the name
// of the binary being run and flags are a slice of command-line flags to be parsed
func ParseInvocation(name string, flags []string) (*Invocation, error) {
	autoTestCmd := flag.NewFlagSet(name, flag.ExitOnError)

	gdbPathFlag := autoTestCmd.String("gdb", "", "Path to the GDB executable to use. Defaults to ${RISCV_PREFIX}gdb")
	serialDeviceFlag := autoTestCmd.String("serial", "", "Device to use for serial communication with the running test")
	autoTestCmd.String("test-name", "", "Name of the test to run")

	if err := autoTestCmd.Parse(flags); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	if len(*serialDeviceFlag) == 0 {
		return nil, errors.New("missing required flag: serial")
	}

	testNameFlag := autoTestCmd.Lookup("test-name")
	if testNameFlag == nil {
		return nil, errors.New("missing required flag: test-name")
	}

	testName := strings.ToLower(testNameFlag.Value.String())

	gdbPath := *gdbPathFlag

	if len(gdbPath) == 0 {
		gdbPath = os.Getenv("RISCV_PREFIX") + "gdb"
	}

	allTests := make([]string, 0, len(testMap))
	for k := range testMap {
		allTests = append(allTests, k)
	}

	if _, ok := testMap[testName]; !ok {
		return nil, fmt.Errorf("invalid test name %q; must be one of %s", testName, strings.Join(allTests, " | "))
	}

	return &Invocation{
		gdbPath:      gdbPath,
		gdbPort:      ":1234",
		serialDevice: *serialDeviceFlag,
		testName:     testName,
	}, nil
}

// Invoke parses and runs an invocation of ss-autotest, where name is the name
// of the binary being run and flags are a slice of command-line flags to be parsed
func Invoke(ctx context.Context, name string, flags []string) error {
	invocation, err := ParseInvocation(name, flags)
	if err != nil {
		return err
	}

	return invocation.Run(ctx)
}

// Run executes ss-autotest for the configured invocation
func (a *Invocation) Run(ctx context.Context) error {
	logger := log.New(os.Stdout, "", 0)

	logger.Printf("processing autotest for '%s'", a.testName)

	logger.Printf("attempting to connect to GDB on port %s", a.gdbPort)

	gdbConn, err := substratum.NewGDBConnection(a.gdbPath, a.gdbPort)
	if err != nil {
		return err
	}

	testFn, ok := testMap[a.testName]
	if !ok {
		panic("invalid test name when running autotest invocation")
	}

	logger.Printf("starting GDB")

	testState, err := autotest.NewState(ctx, logger, gdbConn, "abc")
	if err != nil {
		return err
	}

	defer func() {
		err := testState.Close()
		if err != nil {
			logger.Printf("failed to close: %s", err.Error())
		}
	}()

	err = testFn(ctx, testState)
	if err != nil {
		return err
	}

	return nil
}
