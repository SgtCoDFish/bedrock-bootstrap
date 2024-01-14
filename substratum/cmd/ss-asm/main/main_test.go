package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
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
	outputFormat   string
	expectedOutput []byte
}

func TestSsASM(t *testing.T) {
	var cases = map[string]testCase{
		"hex output": {
			input:          "bne a0 x0 -8\n",
			args:           []string{"-input", "-", "-output", "-"},
			outputFormat:   "hex",
			expectedOutput: []byte("fe 05 1c e3 # bne a0 x0 -8\n"),
			// 0xFE051CE3
		},
		"bin output": {
			input:          "bne a0 x0 -8\n",
			args:           []string{"-input", "-", "-output", "-"},
			outputFormat:   "bin",
			expectedOutput: []byte{byte(0xE3), byte(0x1C), byte(0x05), byte(0xFE)},
			// 0xFE051CE3
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			c.args = append(c.args, "-output-format")
			c.args = append(c.args, c.outputFormat)

			invocation, err := ssasm.ParseInvocation("ss-asm-test", c.args)
			if err != nil {
				t.Errorf("failed to set up test with args '(%v)': %v", c.args, err)
				return
			}

			outBuf := &closerBuffer{}

			invocation.Input = bufio.NewReader(strings.NewReader(c.input))
			invocation.Output = outBuf

			err = invocation.Run(ctx)
			if err != nil {
				t.Errorf("got an error response from running ss-asm(%v): %v", c.args, err)
				return
			}

			output := outBuf.Bytes()

			if !bytes.Equal(output, c.expectedOutput) {
				var renderedOutput string
				var wantedOutput string

				if c.outputFormat == "bin" {
					renderedOutput = fmt.Sprintf("%+v", output)
					wantedOutput = fmt.Sprintf("%+v", c.expectedOutput)
				} else if c.outputFormat == "hex" {
					renderedOutput = fmt.Sprintf("%s", string(output))
					wantedOutput = fmt.Sprintf("%s", string(c.expectedOutput))
				} else {
					t.Errorf("test setup error: unknown output format %q", c.outputFormat)
					renderedOutput = "<unknown>"
				}

				t.Errorf("got unexpected output from input %q\nwanted: %+v\ngot   : %s\n", c.input, wantedOutput, renderedOutput)
				return
			}
		})
	}
}
