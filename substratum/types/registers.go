package types

import (
	"fmt"
	"strconv"
	"strings"
)

var registerMap = map[string]uint8{
	"fp":   0x08,
	"zero": 0x00,
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
