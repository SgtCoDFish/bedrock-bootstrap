package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
			input:          "addi x0 x0 0\n",
			args:           []string{"-input", "-", "-output", "-"},
			outputFormat:   "hex",
			expectedOutput: []byte("13 00 00 00  # addi x0 x0 0\n"),
		},
		"bin output": {
			input:          "addi x0 x0 0\n",
			args:           []string{"-input", "-", "-output", "-"},
			outputFormat:   "bin",
			expectedOutput: []byte{byte(0x13), byte(0x00), byte(0x00), byte(0x00)},
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

			invocation.Input = io.NopCloser(strings.NewReader(c.input))
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
