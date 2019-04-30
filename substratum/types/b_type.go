package types

import (
	"encoding/binary"
	"fmt"
)

type BTypeArgs struct {
	Rs1       uint8
	Rs2       uint8
	Immediate int32
}

func (a *BTypeArgs) Verify() error {
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

func (a *BTypeArgs) Sanitize() {
	a.Rs1 &= 0x1F
	a.Rs2 &= 0x1F
	a.Immediate &= 0xFFF
}

type BTypeInstruction struct {
	Opcode uint8
	Funct3 uint8
	Args   BTypeArgs
}

func (i *BTypeInstruction) Assemble() []byte {
	insn := uint32(0)
	fmt.Printf("[  11]: %d\n[ 1:4]: %d\n[5:11]: %d\n[  12]: %d\n",
		(uint32(i.Args.Immediate)&0x400)>>10,
		(i.Args.Immediate&0x1E)>>1,
		(i.Args.Immediate&0x7E0)>>5,
		(i.Args.Immediate&0x800)>>11,
	)

	// TODO: Remove debug above and fix

	insn |= uint32(i.Opcode & 0x7F)
	insn |= uint32((i.Args.Immediate&0x400)>>10) << 7
	insn |= uint32((i.Args.Immediate&0x1E)>>1) << 8
	insn |= uint32(i.Funct3&0x7) << 12
	insn |= uint32(i.Args.Rs1&0x1F) << 15
	insn |= uint32(i.Args.Rs2&0x1F) << 20
	insn |= uint32((i.Args.Immediate&0x7E0)>>5) << 25
	insn |= uint32((i.Args.Immediate&0x800)>>11) << 31

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, insn)

	return b
}

func NewBEQ(args BTypeArgs) BTypeInstruction {
	return BTypeInstruction{
		Opcode: 0x63,
		Funct3: 0x0,
		Args:   args,
	}
}
