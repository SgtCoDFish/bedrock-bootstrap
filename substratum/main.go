package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sgtcodfish/substratum/types"
)

func process(rawArgs []string) ([]byte, error) {
	instructionName := strings.ToLower(rawArgs[0])

	instruction, err := types.GetInstructionByName(instructionName)
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

	// s := types.NewSW(types.STypeArgs{
	// 	Rs1:       0xa,
	// 	Rs2:       0xf,
	// 	Immediate: 0x00,
	// })

	// out_s := s.Assemble()
	// fmt.Printf("sw\n--\n%v\n", hex.Dump(out_s))

	// u := types.NewLUI(types.UTypeArgs{
	// 	Immediate: 0x10013000,
	// 	Rd:        0xF,
	// })

	// out_u := u.Assemble()
	// fmt.Printf("lui\n---\n%v\n", hex.Dump(out_u))

	// r := types.NewADD(types.RTypeArgs{
	// 	Rd:  0xa,
	// 	Rs1: 0xb,
	// 	Rs2: 0xa,
	// })

	// out_r := r.Assemble()
	// fmt.Printf("add\n---\n%v\n", hex.Dump(out_r))

	// b1 := types.NewBEQ(types.BTypeArgs{
	// 	Immediate: -0x8,
	// 	Rs1:       0x0a,
	// 	Rs2:       0x00,
	// })

	// out_b1 := b1.Assemble()
	// fmt.Printf("beq -ve\n------\n%v\n", hex.Dump(out_b1))

	// b2 := types.NewBEQ(types.BTypeArgs{
	// 	Immediate: 0x20,
	// 	Rs1:       0x08,
	// 	Rs2:       0x00,
	// })

	// out_b2 := b2.Assemble()
	// fmt.Printf("beq +ve\n------\n%v\n", hex.Dump(out_b2))
}
