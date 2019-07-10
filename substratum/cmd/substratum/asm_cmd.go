package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/sgtcodfish/substratum"
)

func processASM(flags *flag.FlagSet) error {
	rawArgs := flags.Args()
	instructionName := strings.ToLower(rawArgs[0])

	instruction, err := substratum.GetInstructionByName(instructionName)
	if err != nil {
		return err
	}

	if err = CheckTailArgs(rawArgs[1:], instruction.ArgumentCount()); err != nil {
		return err
	}

	args := make([]string, len(rawArgs[1:]))
	for i, s := range rawArgs[1:] {
		args[i] = strings.TrimRight(s, ",")
	}

	out, err := instruction.AssembleRaw(args)
	if err != nil {
		return err
	}

	for _, b := range out {
		fmt.Printf("%08b ", b)
	}
	fmt.Print("\n")

	for _, b := range out {
		fmt.Printf("%02x       ", b)
	}

	fmt.Print("\n")

	return nil
}
