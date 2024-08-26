package main

import (
	"bytes"
	"context"
	"testing"

	sstemplate "github.com/sgtcodfish/substratum/cmd/ss-template"
)

type testCase struct {
	args []string

	expectedOutput string
}

const (
	outLine = "13 00 00 00 # nop\n13 00 00 00 # nop\n13 00 00 00 # nop\n13 00 00 00 # nop\n"
)

func TestSSTemplate(t *testing.T) {
	cases := map[string]testCase{
		"size=4": {
			args:           []string{"-size", "4"},
			expectedOutput: outLine,
		},
		"size=8": {
			args:           []string{"-size", "8"},
			expectedOutput: outLine + "\n" + outLine,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			inv, err := sstemplate.ParseInvocation("ss-template-test", c.args)
			if err != nil {
				t.Errorf("failed to set up test with args (%v): %v", c.args, err)
				return
			}

			outBuf := &bytes.Buffer{}

			inv.Output = outBuf

			err = inv.Run(ctx)
			if err != nil {
				t.Errorf("got unepected error from running ss-template (%v): %v", c.args, err)
				return
			}

			output := outBuf.String()

			if output != c.expectedOutput {
				t.Errorf("wanted output: %s\n   got output: %s", c.expectedOutput, output)
				return
			}
		})
	}
}
