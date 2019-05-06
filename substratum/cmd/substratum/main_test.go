package main

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"
)

type testCase struct {
	args     []string
	expected uint32
}

func TestITypes(t *testing.T) {
	var cases = []testCase{
		{
			args:     []string{"addi", "a1", "x0", "0x3"},
			expected: 0x00300593,
		},
		{
			args:     []string{"addi", "a5", "a5", "0x3C"},
			expected: 0x03c78793,
		},
		{
			args:     []string{"ADDI", "x16", "x15", "0x4"},
			expected: 0x00478813,
		},
		{
			args:     []string{"slli", "a0", "a0", "0x10"},
			expected: 0x01051513,
		},
		{
			args:     []string{"slli", "a1", "a1", "16"},
			expected: 0x01059593,
		},
		{
			args:     []string{"andi", "x10", "x10", "0xFF"},
			expected: 0x0ff57513,
		},
	}

	for _, c := range cases {
		out, err := process(c.args)
		if err != nil {
			t.Errorf("got an error response from process(%v): %v", c.args, err)
			continue
		}

		expBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(expBytes, c.expected)

		if !bytes.Equal(out, expBytes) {
			t.Errorf("%s: calculated value did not match expected value: %x != %x", strings.Join(c.args, " "),
				out, expBytes)
		}
	}
}
