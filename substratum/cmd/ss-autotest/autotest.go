package autotest

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"
)

// Invocation holds state for a given invocation of the ss-autotest command
type Invocation struct {
	gdbPath      string
	qemuPath     string
	serialDevice string
	testName     string
}

// ParseInvocation builds an invocation for a run of the the ss-autotest command, where name is the name
// of the binary being run and flags are a slice of command-line flags to be parsed
func ParseInvocation(name string, flags []string) (*Invocation, error) {
	autoTestCmd := flag.NewFlagSet(name, flag.ExitOnError)

	autoTestCmd.String("gdb", "", "The path to the GDB executable to use")
	autoTestCmd.String("qemu", "", "The path to the QEMU executable to use")
	autoTestCmd.String("serial", "", "The serial device to use for communication")
	autoTestCmd.String("test-name", "", "The name of the test to run")

	if err := autoTestCmd.Parse(flags); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	serialDevice := autoTestCmd.Lookup("serial")
	if serialDevice == nil {
		return nil, errors.New("missing required flag: serial")
	}

	testNameFlag := autoTestCmd.Lookup("test-name")
	if testNameFlag == nil {
		return nil, errors.New("missing required flag: test-name")
	}

	testName := strings.ToLower(testNameFlag.Value.String())

	validTestNames := map[string]bool{
		"uart-rxxd-basic":   true,
		"uart-rxxd-comment": true,
		"uart-rxxd-full":    true,
	}

	allTests := make([]string, 0, len(validTestNames))
	for k := range validTestNames {
		allTests = append(allTests, k)
	}

	if _, ok := validTestNames[testName]; !ok {
		return nil, fmt.Errorf("invalid test name %q; must be one of %s", testName, strings.Join(allTests, " | "))
	}

	return &Invocation{
		gdbPath:      "TBC",
		qemuPath:     "TBC",
		serialDevice: "TBC",
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
	return fmt.Errorf("NYI")
}
