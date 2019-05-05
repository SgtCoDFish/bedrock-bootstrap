package types

import (
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

func NewBEQ(args BTypeArgs) BTypeInstruction {
	return BTypeInstruction{
		Opcode: 0x63,
		Funct3: 0x0,
		Args:   args,
	}
}
