package main

import (
	"encoding/hex"
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

	out, err := process(os.Args[1:])
	if err != nil {
		log.Fatalf("fatal error: %v", err)
	}

	fmt.Printf("val: %s", hex.Dump(out))
}
