package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	ssasm "github.com/sgtcodfish/substratum/cmd/ss-asm"
	ssautotest "github.com/sgtcodfish/substratum/cmd/ss-autotest"
)

// ASMCommand is the command which runs the basic Substratum RISC-V "assembler"
const ASMCommand = "asm"

// AutoTestCMD is the command which runs automated, GDB-backed tests of RISC-V baremetal programs
const AutoTestCMD = "autotest"

func main() {
	ctx := context.Background()

	logger := log.New(os.Stdout, "", 0)

	if len(os.Args) < 2 {
		logger.Fatalf("missing required argument: command (one of '%s')", strings.Join([]string{ASMCommand, AutoTestCMD}, ", "))
	}

	subcommandName := strings.ToLower(os.Args[1])
	fullCommandName := fmt.Sprintf("%s %s", os.Args[0], subcommandName)

	switch subcommandName {
	case ASMCommand:
		err := ssasm.Invoke(ctx, fullCommandName, os.Args[2:])
		if err != nil {
			logger.Fatalf("failed to run '%s' command: %s", ASMCommand, err.Error())
		}

	case AutoTestCMD:
		err := ssautotest.Invoke(ctx, fullCommandName, os.Args[2:])
		if err != nil {
			logger.Fatalf("failed to run '%s' command: %s", AutoTestCMD, err.Error())
		}

	default:
		logger.Fatalf("unrecognised command '%s'", os.Args[1])
	}
}
