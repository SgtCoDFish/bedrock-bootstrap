package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

// ASMCommand is the command which runs the basic Substratum RISC-V "assembler"
const ASMCommand = "asm"

// AutoTestCMD is the command which runs automated, GDB-backed tests of RISC-V baremetal programs
const AutoTestCMD = "autotest"

// ASMFileCommand is the command which assembles a file, line by line
const ASMFileCommand = "asm-file"

func main() {
	logger := log.New(os.Stdout, "", 0)

	asmCmd := flag.NewFlagSet(ASMCommand, flag.ExitOnError)

	asmFileCMD := flag.NewFlagSet(ASMFileCommand, flag.ExitOnError)
	asmFileCMD.String("input", "-", "The file to read")
	asmFileCMD.String("output", "-", "The file to write")
	asmFileCMD.String("output-format", "bin", "The format in which to write output - 'bin' for binary, 'hex' for ASCII hex")

	autoTestCMD := flag.NewFlagSet(AutoTestCMD, flag.ExitOnError)
	autoTestCMD.String("gdb", "", "The path to the GDB executable to use")
	autoTestCMD.String("qemu", "", "The path to the QEMU executable to use")
	autoTestCMD.String("serial", "", "The serial device to use for communication")
	autoTestCMD.String("test-name", "", "The name of the test to run")

	if len(os.Args) < 2 {
		log.Fatalf("missing required argument: command (one of '%s')", strings.Join([]string{ASMCommand, ASMFileCommand, AutoTestCMD}, ", "))
	}

	switch strings.ToLower(os.Args[1]) {
	case ASMCommand:
		if len(os.Args) < 3 {
			logger.Fatalf("missing required arguments for asm command")
		}

		err := asmCmd.Parse(os.Args[2:])
		if err != nil {
			logger.Fatalf("failed to parse 'asm' command: %s", err.Error())
		}

		err = processASM(asmCmd, logger)
		if err != nil {
			logger.Fatalf("failed to process 'asm' command: %s", err.Error())
		}

	case ASMFileCommand:
		err := asmFileCMD.Parse(os.Args[2:])
		if err != nil {
			logger.Fatalf("failed to parse '%s' command: %s", ASMFileCommand, err.Error())
		}

		err = processASMFile(asmFileCMD, logger)
		if err != nil {
			logger.Fatalf("failed to process '%s' command: %s", ASMFileCommand, err.Error())
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
