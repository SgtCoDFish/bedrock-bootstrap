package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

const ASMCommand = "asm"
const AutoTestCMD = "autotest"

func main() {
	asmCmd := flag.NewFlagSet(ASMCommand, flag.ExitOnError)

	autoTestCMD := flag.NewFlagSet(AutoTestCMD, flag.ExitOnError)
	autoTestCMD.String("gdb", "", "The path to the GDB executable to use")
	autoTestCMD.String("qemu", "", "The path to the QEMU executable to use")

	if len(os.Args) < 2 {
		log.Fatalf("missing required argument: command (one of '%s')", strings.Join([]string{ASMCommand, AutoTestCMD}, ", "))
	}

	switch strings.ToLower(os.Args[1]) {
	case ASMCommand:
		err := asmCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatalf("failed to parse 'asm' command: %s", err.Error())
		}

		err = processASM(asmCmd)
		if err != nil {
			log.Fatalf("failed to process 'asm' command: %s", err.Error())
		}

	case AutoTestCMD:
		err := autoTestCMD.Parse(os.Args[2:])
		if err != nil {
			log.Fatalf("failed to parse 'autotest' command: %s", err.Error())
		}

		log.Fatalf("autotest nyi")

	default:
		log.Fatalf("unrecognised command '%s'", os.Args[1])
	}
}
