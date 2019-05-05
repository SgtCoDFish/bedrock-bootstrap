package types

import (
	"fmt"
)

type ITypeConstructor func(args ITypeArgs) ITypeInstruction

type ITypeArgs struct {
	Rd        uint8
	Rs1       uint8
	Immediate uint16
}

func (a ITypeArgs) Verify() error {
	if a.Rd > 0x1F {
		return fmt.Errorf("invalid Rd register value %d", a.Rd)
	}

	if a.Rs1 > 0x1F {
		return fmt.Errorf("invalid Rs1 register value %d", a.Rs1)
	}

	if a.Immediate > 0xFFF {
		return fmt.Errorf("invalid immediate value: %d", a.Immediate)
	}

	return nil
}

func (a *ITypeArgs) Sanitize() {
	a.Rd &= 0x1F
	a.Rs1 &= 0x1F
	a.Immediate &= 0xFFF
}

type ITypeInstruction struct {
	Opcode uint8
	Funct3 uint8
	Args   ITypeArgs
}

func NewADDI(args ITypeArgs) ITypeInstruction {
	return ITypeInstruction{
		Opcode: 0x13,
		Funct3: 0x0,
		Args:   args,
	}
}
