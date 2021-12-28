package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	ssasm "github.com/sgtcodfish/substratum/cmd/ss-asm"
)

// ASMCommand is the command which runs the basic Substratum RISC-V "assembler"
const ASMCommand = "asm"

// AutoTestCMD is the command which runs automated, GDB-backed tests of RISC-V baremetal programs
const AutoTestCMD = "autotest"

// ASMFileCommand is the command which assembles a file, line by line
const ASMFileCommand = "asm-file"

func main() {
	ctx := context.Background()

	logger := log.New(os.Stdout, "", 0)

	autoTestCMD := flag.NewFlagSet(AutoTestCMD, flag.ExitOnError)
	autoTestCMD.String("gdb", "", "The path to the GDB executable to use")
	autoTestCMD.String("qemu", "", "The path to the QEMU executable to use")
	autoTestCMD.String("serial", "", "The serial device to use for communication")
	autoTestCMD.String("test-name", "", "The name of the test to run")

	if len(os.Args) < 2 {
		log.Fatalf("missing required argument: command (one of '%s')", strings.Join([]string{ASMCommand, ASMFileCommand, AutoTestCMD}, ", "))
	}

	subcommandName := strings.ToLower(os.Args[1])
	fullCommandName := fmt.Sprintf("%s %s", os.Args[0], subcommandName)

	switch subcommandName {
	case ASMCommand:
		err := ssasm.Invoke(ctx, fullCommandName, os.Args[2:])
		if err != nil {
			logger.Fatalf("failed to run '%s' command: %s", ASMCommand, err)
		}

	case AutoTestCMD:
		err := autoTestCMD.Parse(os.Args[2:])
		if err != nil {
			logger.Fatalf("failed to parse 'autotest' command: %s", err.Error())
		}

		err = processAutotest(autoTestCMD, logger)
		if err != nil {
			logger.Fatalf("failed to process 'autotest' command successfully: %s", err.Error())
		}

	default:
		logger.Fatalf("unrecognised command '%s'", os.Args[1])
	}
}
