package types

import (
	"encoding/binary"
	"fmt"
)

type UTypeArgs struct {
	Rd        uint8
	Immediate uint32
}

func (a *UTypeArgs) Verify() error {
	if a.Rd > 0x1F {
		return fmt.Errorf("invalid Rd register value %d", a.Rd)
	}

	if a.Immediate < 0x1000 {
		return fmt.Errorf("invalid immediate value: %d", a.Immediate)
	}

	return nil
}

func (a *UTypeArgs) Sanitize() {
	a.Rd &= 0x1F
	a.Immediate &= 0xFFFFF000
}

type UTypeInstruction struct {
	Opcode uint8
	Args   UTypeArgs
}

func (i *UTypeInstruction) Assemble() []byte {
	insn := uint32(0)

	insn |= uint32(i.Opcode & 0x7F)
	insn |= uint32(i.Args.Rd&0x1F) << 7
	insn |= i.Args.Immediate & 0xFFFFF000

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, insn)

	return b
}

func NewLUI(args UTypeArgs) UTypeInstruction {
	return UTypeInstruction{
		Opcode: 0x37,
		Args:   args,
	}
}
