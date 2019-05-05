package types

import (
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

func NewSW(args STypeArgs) STypeInstruction {
	return STypeInstruction{
		Opcode: 0x23,
		Funct3: 0x02,
		Args:   args,
	}
}
