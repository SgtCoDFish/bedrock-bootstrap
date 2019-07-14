package substratum

import (
	"fmt"
	"strconv"
	"strings"
)

// registerList is a list of all the integer register names in RISC-V 32-bit
var registerList = []string{
	"zero",
	"ra",
	"sp",
	"gp",
	"tp",
	"t0",
	"t1",
	"t2",
	"fp",
	"s1",
	"a0",
	"a1",
	"a2",
	"a3",
	"a4",
	"a5",
	"a6",
	"a7",
	"s2",
	"s3",
	"s4",
	"s5",
	"s6",
	"s7",
	"s8",
	"s9",
	"s10",
	"s11",
	"t3",
	"t4",
	"t5",
	"t6",
}

var registerMap = map[string]uint8{
	"zero": 0x00,
	"fp":   0x08,
	"ra":   0x01,
	"sp":   0x02,
	"gp":   0x03,
	"tp":   0x04,
	"t0":   0x05,
	"t1":   0x06,
	"t2":   0x07,
	"s0":   0x08,
	"s1":   0x09,
	"a0":   0x0a,
	"a1":   0x0b,
	"a2":   0x0c,
	"a3":   0x0d,
	"a4":   0x0e,
	"a5":   0x0f,
	"a6":   0x10,
	"a7":   0x11,
	"s2":   0x12,
	"s3":   0x13,
	"s4":   0x14,
	"s5":   0x15,
	"s6":   0x16,
	"s7":   0x17,
	"s8":   0x18,
	"s9":   0x19,
	"s10":  0x1a,
	"s11":  0x1b,
	"t3":   0x1c,
	"t4":   0x1d,
	"t5":   0x1e,
	"t6":   0x1f,
}

var registerNumberToABIName = map[string]string{
	"x0":  "zero",
	"x1":  "ra",
	"x2":  "sp",
	"x3":  "gp",
	"x4":  "tp",
	"x5":  "t0",
	"x6":  "t1",
	"x7":  "t2",
	"x8":  "fp",
	"x9":  "s1",
	"x10": "a0",
	"x11": "a1",
	"x12": "a2",
	"x13": "a3",
	"x14": "a4",
	"x15": "a5",
	"x16": "a6",
	"x17": "a7",
	"x18": "s2",
	"x19": "s3",
	"x20": "s4",
	"x21": "s5",
	"x22": "s6",
	"x23": "s7",
	"x24": "s8",
	"x25": "s9",
	"x26": "s10",
	"x27": "s11",
	"x28": "t3",
	"x29": "t4",
	"x30": "t5",
	"x31": "t6",
}

// GetRegisterList returns a slice containing all RISC-V rv32 integer register names
func GetRegisterList() []string {
	return registerList[:]
}

// GetRegisterValue returns the hex value of a register from a register name
// For example, x15 returns 15, and t6 returns 31
func GetRegisterValue(name string) (uint8, error) {
	name = strings.ToLower(name)
	if strings.HasPrefix(name, "x") {
		name = strings.TrimLeft(name, "x")

		a, err := strconv.Atoi(name)
		return uint8(a), err
	}

	val, ok := registerMap[name]
	if !ok {
		return 0, fmt.Errorf("unknown register %s", name)
	}

	return val, nil
}

// GetABINameForNumberRegister returns the ABI name e.g. "t6" corresponding to the given number register
// such as "x31".
func GetABINameForNumberRegister(numberRegister string) (string, error) {
	name, ok := registerNumberToABIName[strings.ToLower(numberRegister)]
	if !ok {
		return "x0", fmt.Errorf("invalid number register: %s", numberRegister)
	}

	return name, nil
}
