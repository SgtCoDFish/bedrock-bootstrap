package substratum

import (
	"fmt"
	"strconv"
)

// Args contains all possible arguments to any instruction type. Not all arguments
// will be used for a given instruction
type Args struct {
	set uint8

	Rd  uint8
	Rs1 uint8
	Rs2 uint8

	Immediate uint32
}

func parseBTypeArgs(args []string) (Args, error) {
	a := Args{}

	rs1, err := GetRegisterValue(args[0])
	if err != nil {
		return a, err
	}

	rs2, err := GetRegisterValue(args[1])
	if err != nil {
		return a, err
	}

	imm, err := strconv.ParseInt(args[2], 0, 13)
	if err != nil {
		return a, fmt.Errorf("invalid immediate value: %s", args[2])
	}

	return Args{
		set:       1,
		Rs1:       rs1,
		Rs2:       rs2,
		Immediate: uint32(imm & 0x1FFF),
	}, nil
}

func parseITypeArgs(args []string) (Args, error) {
	a := Args{}

	rd, err := GetRegisterValue(args[0])
	if err != nil {
		return a, err
	}

	rs1, err := GetRegisterValue(args[1])
	if err != nil {
		return a, err
	}

	imm, err := strconv.ParseInt(args[2], 0, 12)
	if err != nil {
		return a, fmt.Errorf("invalid immediate value: %s", args[2])
	}

	return Args{
		set:       1,
		Rd:        rd,
		Rs1:       rs1,
		Immediate: uint32(imm & 0xFFF),
	}, nil
}

func parseJTypeArgs(args []string) (Args, error) {
	a := Args{}

	rd, err := GetRegisterValue(args[0])
	if err != nil {
		return a, err
	}

	imm, err := strconv.ParseInt(args[2], 0, 21)
	if err != nil {
		return a, fmt.Errorf("invalid immediate value: %s", args[2])
	}

	return Args{
		set:       1,
		Rd:        rd,
		Immediate: uint32(imm & 0x2FFFFFF),
	}, nil
}

func parseRTypeArgs(args []string) (Args, error) {
	a := Args{}

	rd, err := GetRegisterValue(args[0])
	if err != nil {
		return a, err
	}

	rs1, err := GetRegisterValue(args[1])
	if err != nil {
		return a, err
	}

	rs2, err := GetRegisterValue(args[2])
	if err != nil {
		return a, err
	}

	return Args{
		set: 1,
		Rd:  rd,
		Rs1: rs1,
		Rs2: rs2,
	}, nil

}

func parseSTypeArgs(args []string) (Args, error) {
	a := Args{}

	rs1, err := GetRegisterValue(args[0])
	if err != nil {
		return a, err
	}

	rs2, err := GetRegisterValue(args[1])
	if err != nil {
		return a, err
	}

	imm, err := strconv.ParseInt(args[2], 0, 12)
	if err != nil {
		return a, fmt.Errorf("invalid immediate value: %s", args[2])
	}

	return Args{
		set:       1,
		Rs1:       rs1,
		Rs2:       rs2,
		Immediate: uint32(imm & 0xFFF),
	}, nil
}

func parseUTypeArgs(args []string) (Args, error) {
	a := Args{}

	rd, err := GetRegisterValue(args[0])
	if err != nil {
		return a, err
	}

	imm, err := strconv.ParseInt(args[1], 0, 32)
	if err != nil {
		return a, fmt.Errorf("invalid immediate value: %s", args[1])
	}

	return Args{
		set:       1,
		Rd:        rd,
		Immediate: uint32(imm & 0xFFFFF),
	}, nil
}
