package types

import (
	"encoding/binary"
	"fmt"
)

type STypeArgs struct {
	Immediate uint16
	Rs1       uint8
	Rs2       uint8
}

func (a *STypeArgs) Verify() error {
	if a.Rs1 > 0x1F {
		return fmt.Errorf("invalid Rs1 register value %d", a.Rs1)
	}

	if a.Rs2 > 0x1F {
		return fmt.Errorf("invalid Rs2 register value %d", a.Rs2)
	}

	if a.Immediate > 0xFFF {
		return fmt.Errorf("invalid immediate value: %d", a.Immediate)
	}

	return nil
}

func (a *STypeArgs) Sanitize() {
	a.Rs1 &= 0x1F
	a.Rs2 &= 0x1F
	a.Immediate &= 0xFFF
}

type STypeInstruction struct {
	Opcode uint8
	Funct3 uint8
	Args   STypeArgs
}

func (i *STypeInstruction) Assemble() []byte {
	insn := uint32(0)

	insn |= uint32(i.Opcode & 0x7F)
	insn |= uint32(i.Args.Immediate&0x1F) << 7
	insn |= uint32(i.Funct3&0x7) << 12
	insn |= uint32(i.Args.Rs2&0x1F) << 15
	insn |= uint32(i.Args.Rs1&0x1F) << 20
	insn |= uint32(i.Args.Immediate&0xFE0) << 25

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, insn)

	return b
}

func NewSW(args STypeArgs) STypeInstruction {
	return STypeInstruction{
		Opcode: 0x23,
		Funct3: 0x02,
		Args:   args,
	}
}
