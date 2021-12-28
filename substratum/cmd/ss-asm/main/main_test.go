package main

import (
	"bufio"
	"bytes"
	"context"
	"strings"
	"testing"

	ssasm "github.com/sgtcodfish/substratum/cmd/ss-asm"
)

type closerBuffer struct {
	bytes.Buffer
}

func (c *closerBuffer) Close() error {
	return nil
}

type testCase struct {
	args           []string
	input          string
	expectedOutput []byte
}

func TestSsASM(t *testing.T) {
	t.Parallel()

	var cases = []testCase{
		{
			input:          "bne a0 x0 -8",
			args:           []string{"-input", "-", "-output", "-", "-output-format", "hex"},
			expectedOutput: []byte("fe 05 1c e3 # bne a0 x0 -8\n"),
			// 0xFE051CE3
		},
		{
			input:          "bne a0 x0 -8",
			args:           []string{"-input", "-", "-output", "-", "-output-format", "hex"},
			expectedOutput: []byte{byte(0xFE), byte(0x05), byte(0x1C), byte(0xE3)},
			// 0xFE051CE3
		},
	}

	for _, c := range cases {
		ctx := context.Background()

		outBuf := &closerBuffer{}

		invocation, err := ssasm.ParseInvocation("ss-asm-test", c.args)
		if err != nil {
			t.Errorf("failed to set up test with args '(%v)': %v", c.args, err)
			continue
		}

		invocation.Input = bufio.NewReader(strings.NewReader(c.input))
		invocation.Output = outBuf

		err = invocation.Run(ctx)
		if err != nil {
			t.Errorf("got an error response from running ss-asm(%v): %v", c.args, err)
			continue
		}

		output := outBuf.Bytes()

		if bytes.Equal(output, c.expectedOutput) {
			t.Errorf("got unexpected output from input %q", c.input)
		}
	}
}
