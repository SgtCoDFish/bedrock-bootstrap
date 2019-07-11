package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/sgtcodfish/substratum"
)

func processASM(flags *flag.FlagSet, logger *log.Logger) error {
	rawArgs := flags.Args()
	instructionName := strings.ToLower(rawArgs[0])

	instruction, err := substratum.GetInstructionByName(instructionName)
	if err != nil {
		return err
	}

	if err = CheckTailArgs(rawArgs[1:], instruction.ArgumentCount()); err != nil {
		return fmt.Errorf("invalid argument list or count for '%s': %v", instructionName, err)
	}

	args := make([]string, len(rawArgs[1:]))
	for i, s := range rawArgs[1:] {
		args[i] = strings.TrimRight(s, ",")
	}

	out, err := instruction.AssembleRaw(args)
	if err != nil {
		return err
	}

	builder := new(strings.Builder)
	for _, b := range out {
		builder.WriteString(fmt.Sprintf("%08b ", b))
	}

	builder.WriteString("\n")

	for _, b := range out {
		builder.WriteString(fmt.Sprintf("%02x       ", b))
	}

	builder.WriteString("\n")

	logger.Print(builder.String())

	return nil
}
