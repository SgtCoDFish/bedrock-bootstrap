package types

import (
	"encoding/binary"
	"fmt"
)

type RTypeArgs struct {
	Rd  uint8
	Rs1 uint8
	Rs2 uint8
}

func (a *RTypeArgs) Verify() error {
	if a.Rd > 0x1F {
		return fmt.Errorf("invalid Rd register value %d", a.Rd)
	}

	if a.Rs1 > 0x1F {
		return fmt.Errorf("invalid Rs1 register value %d", a.Rs1)
	}

	if a.Rs2 > 0x1F {
		return fmt.Errorf("invalid Rs2 register value %d", a.Rs2)
	}

	return nil
}

func (a *RTypeArgs) Sanitize() {
	a.Rd &= 0x1F
	a.Rs1 &= 0x1F
	a.Rs2 &= 0x1F
}

type RTypeInstruction struct {
	Opcode uint8
	Funct3 uint8
	Funct7 uint8
	Args   RTypeArgs
}

func (i *RTypeInstruction) Assemble() []byte {
	insn := uint32(0)

	insn |= uint32(i.Opcode & 0x7F)
	insn |= uint32(i.Args.Rd&0x1F) << 7
	insn |= uint32(i.Funct3&0x7) << 12
	insn |= uint32(i.Args.Rs1&0x1F) << 15
	insn |= uint32(i.Args.Rs2&0x1F) << 20
	insn |= uint32(i.Funct7&0x7f) << 25

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, insn)

	return b
}

func NewADD(args RTypeArgs) RTypeInstruction {
	return RTypeInstruction{
		Opcode: 0x33,
		Funct3: 0x7,
		Funct7: 0x0,
		Args:   args,
	}
}
