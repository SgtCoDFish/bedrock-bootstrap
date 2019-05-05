package types

import (
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

func NewLUI(args UTypeArgs) UTypeInstruction {
	return UTypeInstruction{
		Opcode: 0x37,
		Args:   args,
	}
}
