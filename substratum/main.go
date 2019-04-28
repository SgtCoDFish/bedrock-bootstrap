package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/sgtcodfish/substratum/types"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing required argument: instruction (e.g. addi)")
	}

	instruction := strings.ToLower(os.Args[1])

	if instruction != "addi" {
		log.Fatalf("invalid instruction: %s", instruction)
	}

	if len(os.Args) != 5 {
		log.Fatalf("invalid argument count %d. required format: %s %s rd rs1 imm", len(os.Args), os.Args[0], instruction)
	}

	rd, err := types.GetRegisterValue(os.Args[2])
	if err != nil {
		log.Fatalf("invalid register: %v", err)
	}

	rs1, err := types.GetRegisterValue(os.Args[3])
	if err != nil {
		log.Fatalf("invalid register: %v", err)
	}

	imm, err := strconv.ParseInt(os.Args[4], 0, 12)
	if err != nil {
		log.Fatalf("invalid immediate value: %s", os.Args[4])
	}

	i := types.NewADDI(types.ITypeArgs{
		Rd:        rd,
		Rs1:       rs1,
		Immediate: uint16(imm & 0xFFF),
	})

	out := i.Assemble()
	fmt.Printf("%v\n", hex.Dump(out[:]))
}
