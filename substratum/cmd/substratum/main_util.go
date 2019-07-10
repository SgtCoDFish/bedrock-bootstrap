package main

import (
	"fmt"
	"strings"
)

// CheckTailArgs takes a slice of args and returns an error if
// any of those args start with a "-" which might indicate a misplaced
// flag argument. The argument count is also checked; pass "-1" to mean any amount of arguments.
func CheckTailArgs(args []string, argCount int) error {
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			return fmt.Errorf("argument '%s' looks like a flag; all flags should be at the beginning of the subcommand", arg)
		}
	}

	if argCount >= 0 {
		if len(args) != argCount {
			return fmt.Errorf("expected %d arguments but got %d", argCount, len(args))
		}
	}

	return nil
}
