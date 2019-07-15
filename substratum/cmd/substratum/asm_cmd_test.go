package main

import (
	"bytes"
	"flag"
	"log"
	"testing"
)

type testCase struct {
	args []string
	// expected uint32
}

func makeTestLogger() (*log.Logger, *bytes.Buffer) {
	buf := new(bytes.Buffer)

	return log.New(buf, "", 0), buf
}

func TestBTypes(t *testing.T) {
	var cases = []testCase{
		{
			args: []string{"bne", "a0", "x0", "-8"},
			// expected: 0xFE051CE3,
		},
		{
			args: []string{"beq", "a0", "x0", "-0x8"},
			// expected: 0xFE050CE3,
		},
		{
			args: []string{"beq", "x0", "x0", "-0x2C"},
			// expected: 0xFC000AE3,
		},
		{
			args: []string{"beq", "a0", "s1", "+0x10"},
			// expected: 0x00950863,
		},
	}

	logger, buffer := makeTestLogger()

	for _, c := range cases {
		flagSet := flag.NewFlagSet("", flag.ExitOnError)
		err := flagSet.Parse(c.args)
		if err != nil {
			t.Errorf("couldn't parse flags: %+v", c.args)
			continue
		}

		err = processASM(flagSet, logger)
		if err != nil {
			t.Errorf("got an error response from process(%v): %v", c.args, err)
			continue
		}

		output := buffer.String()
		if len(output) == 0 {
			t.Errorf("got a 0 length buffer response but expected some output")
		}

		// expBytes := make([]byte, 4)
		// binary.LittleEndian.PutUint32(expBytes, c.expected)

		// if !bytes.Equal([]byte(output), expBytes) {
		// 	t.Errorf("%s: calculated value did not match expected value: %x != %x", strings.Join(c.args, " "),
		// 		output, expBytes)
		// }
	}
}

//func TestITypes(t *testing.T) {
//	var cases = []testCase{
//		{
//			args:     []string{"addi", "a1", "x0", "0x3"},
//			expected: 0x00300593,
//		},
//		{
//			args:     []string{"addi", "a5", "a5", "0x3C"},
//			expected: 0x03c78793,
//		},
//		{
//			args:     []string{"ADDI", "x16", "x15", "0x4"},
//			expected: 0x00478813,
//		},
//		{
//			args:     []string{"slli", "a0", "a0", "0x10"},
//			expected: 0x01051513,
//		},
//		{
//			args:     []string{"slli", "a1", "a1", "16"},
//			expected: 0x01059593,
//		},
//		{
//			args:     []string{"andi", "x10", "x10", "0xFF"},
//			expected: 0x0ff57513,
//		},
//		{
//			args:     []string{"ori", "x10", "x00", "0x20"},
//			expected: 0x02006513,
//		},
//	}
//
//	for _, c := range cases {
//		err := processASM(c.args)
//		if err != nil {
//			t.Errorf("got an error response from process(%v): %v", c.args, err)
//			continue
//		}
//
//		expBytes := make([]byte, 4)
//		binary.LittleEndian.PutUint32(expBytes, c.expected)
//
//		if !bytes.Equal(out, expBytes) {
//			t.Errorf("%s: calculated value did not match expected value: %x != %x", strings.Join(c.args, " "),
//				out, expBytes)
//		}
//	}
//}
//
//func TestJTypes(t *testing.T) {
//	var cases = []testCase{
//		{
//			args:     []string{"jal", "x0", "0"},
//			expected: 0x0000006F,
//		},
//		{
//			args:     []string{"jal", "ra", "-0xC"},
//			expected: 0xFF5FF0EF,
//		},
//	}
//
//	for _, c := range cases {
//		flagSet := flag.NewFlagSet("", flag.ExitOnError)
//		err := flagSet.Parse(c.args)
//		if err != nil {
//			t.Errorf("couldn't parse flags: %+v", c.args)
//			continue
//		}
//
//		err = processASM(flagSet)
//		if err != nil {
//			t.Errorf("got an error response from process(%v): %v", c.args, err)
//			continue
//		}
//
//		expBytes := make([]byte, 4)
//		binary.LittleEndian.PutUint32(expBytes, c.expected)
//
//		if !bytes.Equal(out, expBytes) {
//			t.Errorf("%s: calculated value did not match expected value: %x != %x", strings.Join(c.args, " "),
//				out, expBytes)
//		}
//	}
//}
//
//func TestRTypes(t *testing.T) {
//	var cases = []testCase{
//		{
//			args:     []string{"sub", "a1", "x0", "a1"},
//			expected: 0x40b005b3,
//		},
//		{
//			args:     []string{"and", "a0", "a1", "a0"},
//			expected: 0x00A5F533,
//		},
//		{
//			args:     []string{"and", "a0", "a0", "a2"},
//			expected: 0x00C57533,
//		},
//		{
//			args:     []string{"or", "a0", "a0", "a1"},
//			expected: 0x00B56533,
//		},
//		{
//			args:     []string{"sll", "ra", "sp", "gp"},
//			expected: 0x003110b3,
//		},
//	}
//
//	for _, c := range cases {
//		out, err := processASM(c.args)
//		if err != nil {
//			t.Errorf("got an error response from process(%v): %v", c.args, err)
//			continue
//		}
//
//		expBytes := make([]byte, 4)
//		binary.LittleEndian.PutUint32(expBytes, c.expected)
//
//		if !bytes.Equal(out, expBytes) {
//			t.Errorf("%s: calculated value did not match expected value: %x != %x", strings.Join(c.args, " "),
//				out, expBytes)
//		}
//	}
//}
//
//func TestUTypes(t *testing.T) {
//	var cases = []testCase{
//		{
//			args:     []string{"lui", "a2", "0x80000"},
//			expected: 0x80000637,
//		},
//		{
//			args:     []string{"lui", "a1", "0xDEAD0"},
//			expected: 0xDEAD05B7,
//		},
//		{
//			args:     []string{"lui", "x14", "524304"}, // 0x80010
//			expected: 0x80010737,
//		},
//		{
//			args:     []string{"lui", "x15", "0x10013"},
//			expected: 0x100137B7,
//		},
//	}
//
//	for _, c := range cases {
//		out, err := processASM(c.args)
//		if err != nil {
//			t.Errorf("got an error response from process(%v): %v", c.args, err)
//			continue
//		}
//
//		expBytes := make([]byte, 4)
//		binary.LittleEndian.PutUint32(expBytes, c.expected)
//
//		if !bytes.Equal(out, expBytes) {
//			t.Errorf("%s: calculated value did not match expected value: %x != %x", strings.Join(c.args, " "),
//				out, expBytes)
//		}
//	}
//}
