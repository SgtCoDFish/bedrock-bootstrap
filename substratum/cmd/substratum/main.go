package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sgtcodfish/substratum"
)

func process(rawArgs []string) ([]byte, error) {
	instructionName := strings.ToLower(rawArgs[0])

	instruction, err := substratum.GetInstructionByName(instructionName)
	if err != nil {
		return nil, err
	}

	args := rawArgs[1:]

	if len(args) != instruction.ArgumentCount() {
		return nil, fmt.Errorf("invalid arg count; got %d but need %d", len(args), instruction.ArgumentCount())
	}

	return instruction.AssembleRaw(args)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing required argument: instruction (e.g. addi)")
	}

	rawArgs := make([]string, len(os.Args[1:]))

	for i, s := range os.Args[1:] {
		rawArgs[i] = strings.TrimRight(s, ",")
	}

	out, err := process(rawArgs)
	if err != nil {
		log.Fatalf("fatal error: %v", err)
	}

	for _, b := range out {
		fmt.Printf("%08b ", b)
	}
	fmt.Print("\n")

	for _, b := range out {
		fmt.Printf("%02x       ", b)
	}

	fmt.Print("\n")
}
