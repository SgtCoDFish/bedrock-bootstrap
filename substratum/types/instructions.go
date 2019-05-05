package types

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type InstructionType int

const (
	BType InstructionType = iota
	IType
	JType
	RType
	SType
	UType
)

// Opcodes holds possible opcodes which can appear in an instruction.
// All instructions have an "Opcode" but only some have "Funct3" and "Funct7"
type Opcodes struct {
	Opcode uint8
	Funct3 uint8
	Funct7 uint8
}

// Instruction is a full representation of a given instruction minus the arguments
// which can vary between different invocations of a given instruction
type Instruction struct {
	Name    string
	Opcodes Opcodes
	Type    InstructionType

	RequiredArgCount uint8
}

func (i Instruction) verifyArgs(args Args) error {
	if args.set == 0 {
		return fmt.Errorf("invalid args: must be created with a constructor")
	}

	return nil
}

// ArgumentFormat returns the required format of the arguments taken by the instruction
func (i Instruction) ArgumentFormat() string {
	var format string

	switch i.Type {
	case BType:
		format = "rs1 rs2 imm"
	case IType:
		format = "rd rs1 imm"
	case JType:
		format = "rd imm"
	case RType:
		format = "rd rs1 rs2"
	case SType:
		format = "rs1 rs2 imm"
	case UType:
		format = "rd imm"
	default:
		panic(fmt.Sprintf("unknown instruction type in ArgumentFormat: %+v", i))
	}

	return fmt.Sprintf("%s %s", i.Name, format)
}

// ArgumentCount returns the number of arguments which are expected for a given named argument
func (i Instruction) ArgumentCount() int {
	switch i.Type {
	case BType:
		return 3
	case IType:
		return 3
	case JType:
		return 2
	case RType:
		return 3
	case SType:
		return 3
	case UType:
		return 2
	default:
		panic(fmt.Sprintf("unknown instruction type in ArgumentCount: %+v", i))
	}
}

// Assemble returns an assembled version of the given Instruction using the provided
// Args, as a little endian byte slice.
func (i Instruction) Assemble(args Args) []byte {
	switch i.Type {
	case BType:
		return i.assembleBType(args)
	case IType:
		return i.assembleIType(args)
	case JType:
		return i.assembleJType(args)
	case RType:
		return i.assembleRType(args)
	case SType:
		return i.assembleSType(args)
	case UType:
		return i.assembleUType(args)
	default:
		panic(fmt.Sprintf("unknown instruction type in Assemble: %+v", i))
	}
}

// AssembleRaw takes raw arguments, parses them and returns an assembled instruction
func (i Instruction) AssembleRaw(rawArgs []string) ([]byte, error) {
	var args Args
	var err error

	switch i.Type {
	case BType:
		args, err = parseBTypeArgs(rawArgs)
	case IType:
		args, err = parseITypeArgs(rawArgs)
	case JType:
		args, err = parseJTypeArgs(rawArgs)
	case RType:
		args, err = parseRTypeArgs(rawArgs)
	case SType:
		args, err = parseSTypeArgs(rawArgs)
	case UType:
		args, err = parseUTypeArgs(rawArgs)
	default:
		panic(fmt.Sprintf("unknown instruction type in Assemble: %+v", i))
	}

	if err != nil {
		return nil, err
	}

	return i.Assemble(args), nil
}

func (i Instruction) assembleBType(args Args) []byte {
	insn := uint32(0)
	fmt.Printf("[  11]: %d\n[ 1:4]: %d\n[5:11]: %d\n[  12]: %d\n",
		(uint32(args.Immediate)&0x400)>>10,
		(args.Immediate&0x1E)>>1,
		(args.Immediate&0x7E0)>>5,
		(args.Immediate&0x800)>>11,
	)

	// TODO: Remove debug above and fix

	insn |= uint32(i.Opcodes.Opcode & 0x7F)
	insn |= uint32((args.Immediate&0x400)>>10) << 7
	insn |= uint32((args.Immediate&0x1E)>>1) << 8
	insn |= uint32(i.Opcodes.Funct3&0x7) << 12
	insn |= uint32(args.Rs1&0x1F) << 15
	insn |= uint32(args.Rs2&0x1F) << 20
	insn |= uint32((args.Immediate&0x7E0)>>5) << 25
	insn |= uint32((args.Immediate&0x800)>>11) << 31

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, insn)

	return b
}

func (i Instruction) assembleIType(args Args) []byte {
	insn := uint32(0)

	insn |= uint32(i.Opcodes.Opcode & 0x7F)
	insn |= uint32(args.Rd&0x1F) << 7
	insn |= uint32(i.Opcodes.Funct3&0x7) << 12
	insn |= uint32(args.Rs1&0x1F) << 15
	insn |= uint32(args.Immediate&0xFFF) << 20

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, insn)

	return b
}

func (i Instruction) assembleJType(args Args) []byte {
	insn := uint32(0)

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, insn)

	panic("NYI: j type")

	return b
}

func (i Instruction) assembleRType(args Args) []byte {
	insn := uint32(0)

	insn |= uint32(i.Opcodes.Opcode & 0x7F)
	insn |= uint32(args.Rd&0x1F) << 7
	insn |= uint32(i.Opcodes.Funct3&0x7) << 12
	insn |= uint32(args.Rs1&0x1F) << 15
	insn |= uint32(args.Rs2&0x1F) << 20
	insn |= uint32(i.Opcodes.Funct7&0x7f) << 25

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, insn)

	return b
}

func (i Instruction) assembleSType(args Args) []byte {
	insn := uint32(0)

	insn |= uint32(i.Opcodes.Opcode & 0x7F)
	insn |= uint32(args.Immediate&0x1F) << 7
	insn |= uint32(i.Opcodes.Funct3&0x7) << 12
	insn |= uint32(args.Rs2&0x1F) << 15
	insn |= uint32(args.Rs1&0x1F) << 20
	insn |= uint32(args.Immediate&0xFE0) << 25

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, insn)

	return b
}

func (i Instruction) assembleUType(args Args) []byte {
	insn := uint32(0)

	insn |= uint32(i.Opcodes.Opcode & 0x7F)
	insn |= uint32(args.Rd&0x1F) << 7
	insn |= args.Immediate & 0xFFFFF000

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, insn)

	return b
}

// GetInstructionByName returns an Instruction struct for the given name if one exists
func GetInstructionByName(name string) (Instruction, error) {
	insn, ok := instructionMap[strings.ToLower(name)]
	if !ok {
		return Instruction{}, fmt.Errorf("unknown instruction type: %s", name)
	}

	return insn, nil
}

var instructionMap = map[string]Instruction{
	"add": {
		Name: "add",
		Type: RType,
		Opcodes: Opcodes{
			Opcode: 0x33,
			Funct3: 0x07,
			Funct7: 0x00,
		},
	},
	"addi": {
		Name: "addi",
		Type: IType,
		Opcodes: Opcodes{
			Opcode: 0x13,
			Funct3: 0x00,
		},
	},
	"beq": {
		Name: "beq",
		Type: BType,
		Opcodes: Opcodes{
			Opcode: 0x63,
			Funct3: 0x00,
		},
	},
	"lui": {
		Name: "lui",
		Type: UType,
		Opcodes: Opcodes{
			Opcode: 0x37,
		},
	},
	"sw": {
		Name: "sw",
		Type: SType,
		Opcodes: Opcodes{
			Opcode: 0x23,
			Funct3: 0x02,
		},
	},
}
