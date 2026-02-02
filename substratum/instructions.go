package substratum

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// InstructionType is a marker for the different possible RISC-V instruction types.
// An example would be "I-type", which the "addi" instruction is formatted using.
type InstructionType int

const (
	// BType is used for RISC-V B-type instructions such as "beq"
	BType InstructionType = iota

	// IType is used for RISC-V I-type instructions such as "addi"
	IType

	// JType is used for RISC-V J-type instructions such as "jal"
	JType

	// RType is used for RISC-V R-type instructions such as "add"
	RType

	// SType is used for RISC-V S-type instructions such as "sw"
	SType

	// UType is used for RISC-V U-type instructions such as "lui"
	UType
)

// Opcodes holds possible opcodes which can appear in an instruction.
// All instructions have an "Opcode" but only some have "Funct3" and "Funct7"
type Opcodes struct {
	Opcode uint8
	Funct3 uint8
	Funct7 uint8

	// ImmediateMask, if set, masks the immediate value to make it smaller. This helps
	// for instructions such as "slli" which don't have a full 12-bit immediate value,
	// and rather have 6-bits of zeroes in the upper part of the immediate.
	ImmediateMask uint32
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
func (i Instruction) Assemble(args Args) ([]byte, error) {
	if err := i.verifyArgs(args); err != nil {
		return nil, err
	}

	switch i.Type {
	case BType:
		return i.assembleBType(args), nil
	case IType:
		return i.assembleIType(args), nil
	case JType:
		return i.assembleJType(args), nil
	case RType:
		return i.assembleRType(args), nil
	case SType:
		return i.assembleSType(args), nil
	case UType:
		return i.assembleUType(args), nil
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

	return i.Assemble(args)
}

func (i Instruction) assembleBType(args Args) []byte {
	// fmt.Printf("%s:\n%032b\n[  11]: %01b\n[ 4:1]: %04b\n[11:5]: %06b\n[  12]: %01b\n---------\n",
	// 	i.Name,
	// 	args.Immediate,
	// 	(uint32(args.Immediate)&0x400)>>10,
	// 	(uint32(args.Immediate&0x1E))>>1,
	// 	(uint32(args.Immediate&0x7E0))>>5,
	// 	(uint32(args.Immediate&0x800))>>11,
	// )

	insn := uint32(0)

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

	imm := args.Immediate & 0xFFF

	if i.Opcodes.ImmediateMask != 0 {
		imm &= i.Opcodes.ImmediateMask
	}

	insn |= uint32(i.Opcodes.Opcode & 0x7F)
	insn |= uint32(args.Rd&0x1F) << 7
	insn |= uint32(i.Opcodes.Funct3&0x7) << 12
	insn |= uint32(args.Rs1&0x1F) << 15
	insn |= uint32(imm) << 20

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, insn)

	return b
}

func (i Instruction) assembleJType(args Args) []byte {
	insn := uint32(0)

	insn |= uint32(i.Opcodes.Opcode & 0x7F)
	insn |= uint32(args.Rd&0x1F) << 7
	insn |= uint32((args.Immediate&0xFF000)>>12) << 12
	insn |= uint32((args.Immediate&0x800)>>11) << 20
	insn |= uint32((args.Immediate&0x7FF)>>1) << 21
	insn |= uint32((args.Immediate&0x100000)>>20) << 31

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, insn)

	// fmt.Printf("%s\n---\n[19:12] %08b\n[   11] %01b\n[10:01] %010b\n[   20] %01b\n%08x\n%032b\n\n\n",
	// 	i.Name,
	// 	uint32((args.Immediate&0xFF000)>>12),
	// 	uint32((args.Immediate&0x800)>>11),
	// 	uint32((args.Immediate&0x7FE)>>1),
	// 	uint32((args.Immediate&0x100000)>>20),
	// 	insn,
	// 	insn,
	// )

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
	insn |= (args.Immediate & 0xFFFFF) << 12

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
			Funct3: 0x00,
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
	"and": {
		Name: "and",
		Type: RType,
		Opcodes: Opcodes{
			Opcode: 0x33,
			Funct3: 0x07,
			Funct7: 0x00,
		},
	},
	"andi": {
		Name: "andi",
		Type: IType,
		Opcodes: Opcodes{
			Opcode: 0x13,
			Funct3: 0x07,
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
	"bne": {
		Name: "bne",
		Type: BType,
		Opcodes: Opcodes{
			Opcode: 0x63,
			Funct3: 0x01,
		},
	},
	"lui": {
		Name: "lui",
		Type: UType,
		Opcodes: Opcodes{
			Opcode: 0x37,
		},
	},
	"auipc": {
		Name: "auipc",
		Type: UType,
		Opcodes: Opcodes{
			Opcode: 0x17,
		},
	},
	"or": {
		Name: "or",
		Type: RType,
		Opcodes: Opcodes{
			Opcode: 0x33,
			Funct3: 0x06,
			Funct7: 0x00,
		},
	},
	"ori": {
		Name: "ori",
		Type: IType,
		Opcodes: Opcodes{
			Opcode: 0x13,
			Funct3: 0x06,
		},
	},
	"sll": {
		Name: "sll",
		Type: RType,
		Opcodes: Opcodes{
			Opcode: 0x33,
			Funct3: 0x01,
			Funct7: 0x00,
		},
	},
	"slli": {
		Name: "slli",
		Type: IType,
		Opcodes: Opcodes{
			Opcode:        0x13,
			Funct3:        0x01,
			Funct7:        0x00,
			ImmediateMask: 0x3F,
		},
	},
	"sub": {
		Name: "sub",
		Type: RType,
		Opcodes: Opcodes{
			Opcode: 0x33,
			Funct3: 0x00,
			Funct7: 0x20,
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
	"sb": {
		Name: "sb",
		Type: SType,
		Opcodes: Opcodes{
			Opcode: 0x23,
			Funct3: 0x00,
		},
	},
	"lw": {
		Name: "lw",
		Type: IType,
		Opcodes: Opcodes{
			Opcode: 0x03,
			Funct3: 0x02,
		},
	},
	"lbu": {
		Name: "lbu",
		Type: IType,
		Opcodes: Opcodes{
			Opcode: 0x03,
			Funct3: 0x04,
		},
	},
	"lb": {
		Name: "lb",
		Type: IType,
		Opcodes: Opcodes{
			Opcode: 0x03,
			Funct3: 0x00,
		},
	},
	"jal": {
		Name: "jal",
		Type: JType,
		Opcodes: Opcodes{
			Opcode: 0x6F,
		},
	},
	"jalr": {
		Name: "jalr",
		Type: IType,
		Opcodes: Opcodes{
			Opcode: 0x67,
			Funct3: 0x00,
		},
	},
	"bge": {
		Name: "bge",
		Type: BType,
		Opcodes: Opcodes{
			Opcode: 0x63,
			Funct3: 0x05,
		},
	},
	"blt": {
		Name: "blt",
		Type: BType,
		Opcodes: Opcodes{
			Opcode: 0x63,
			Funct3: 0x04,
		},
	},
	"xor": {
		Name: "xor",
		Type: BType,
		Opcodes: Opcodes{
			Opcode: 0x33,
			Funct3: 0x4,
			Funct7: 0x00,
		},
	},
	"xori": {
		Name: "xori",
		Type: IType,
		Opcodes: Opcodes{
			Opcode: 0x13,
			Funct3: 0x4,
		},
	},
}
