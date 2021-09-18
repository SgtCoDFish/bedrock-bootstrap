package substratum

import (
	"reflect"
	"testing"
)

func TestAssembleLine(t *testing.T) {
	tests := map[string]struct {
		inputLine      string
		expectedOutput []byte
		expectErr      bool
	}{
		"simple case": {
			inputLine:      "addi x00 x00 0x00\n",
			expectedOutput: []byte{0x13, 0x00, 0x00, 0x00},
			expectErr:      false,
		},
		"without newline": {
			inputLine:      "addi x00 x00 0x00",
			expectedOutput: []byte{0x13, 0x00, 0x00, 0x00},
			expectErr:      false,
		},
		"using 'regular' register names": {
			inputLine:      "addi a0 t0 0x00\n",
			expectedOutput: []byte{0x13, 0x85, 0x02, 0x00},
			expectErr:      false,
		},
		"using decimal number": {
			inputLine:      "addi x0 x0 5\n",
			expectedOutput: []byte{0x13, 0x00, 0x50, 0x00},
			expectErr:      false,
		},
		"using octal number": {
			inputLine:      "addi x0 x0 0o5\n",
			expectedOutput: []byte{0x13, 0x00, 0x50, 0x00},
			expectErr:      false,
		},
		"using binary number": {
			inputLine:      "addi x0 x0 0b0101\n",
			expectedOutput: []byte{0x13, 0x00, 0x50, 0x00},
			expectErr:      false,
		},
		"using hex number": {
			inputLine:      "addi x0 x0 0x05\n",
			expectedOutput: []byte{0x13, 0x00, 0x50, 0x00},
			expectErr:      false,
		},
		"using commas": {
			inputLine:      "addi x0, x0, 0x00\n",
			expectedOutput: []byte{0x13, 0x00, 0x00, 0x00},
			expectErr:      false,
		},
		"using strange commas": {
			inputLine:      "addi x0 ,x0 ,0x00\n",
			expectedOutput: []byte{0x13, 0x00, 0x00, 0x00},
			expectErr:      false,
		},
		"weird spacing fails": {
			inputLine:      "addi x0   \t x0\t\t\t\t 0x00\t\t\t\n",
			expectedOutput: nil,
			expectErr:      true,
		},
		"newline is significant": {
			inputLine:      "addi x0 x0\n 0x00\n",
			expectedOutput: nil,
			expectErr:      true,
		},
		"comments are trimmed": {
			inputLine:      "addi x0 x0 0 # literally anything\n",
			expectedOutput: []byte{0x13, 0x00, 0x00, 0x00},
			expectErr:      false,
		},
		"full line comment is ignored": {
			inputLine:      " # literally anything\n",
			expectedOutput: nil,
			expectErr:      false,
		},
		"only one line processed": {
			inputLine:      "addi x00 x00 0x00\naddi x00 x00 0x00\n",
			expectedOutput: []byte{0x13, 0x00, 0x00, 0x00},
			expectErr:      false,
		},
		"I-type regular form supported": {
			inputLine:      "lw x01 x01 0x05",
			expectedOutput: []byte{0x83, 0xa0, 0x50, 0x00},
			expectErr:      false,
		},
		"I-type brackets form supported": {
			inputLine:      "lw x01 5(x01)",
			expectedOutput: []byte{0x83, 0xa0, 0x50, 0x00},
			expectErr:      false,
		},
		"malformed brackets form error": {
			inputLine:      "lw x01 5((x01)",
			expectedOutput: nil,
			expectErr:      true,
		},
		"regular form incorrect arg count": {
			inputLine:      "lw x01 x01 0x05 0x05",
			expectedOutput: nil,
			expectErr:      true,
		},
		"brackets form incorrect arg count": {
			inputLine:      "lw x01 5(x01) 0x05",
			expectedOutput: nil,
			expectErr:      true,
		},
		"invalid instruction error is passed through": {
			inputLine:      "borked 0x00 0x00 0x00",
			expectedOutput: nil,
			expectErr:      true,
		},
		"instruction with 2 args": {
			inputLine:      "lui x00 0x10000",
			expectedOutput: []byte{0x37, 0x00, 0x00, 0x10},
			expectErr:      false,
		},
		"too few args": {
			inputLine:      "lui x00",
			expectedOutput: nil,
			expectErr:      true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := AssembleLine(test.inputLine)

			if (err != nil) != test.expectErr {
				t.Errorf("err=%v, expectErr=%v", err, test.expectErr)
			}

			if !reflect.DeepEqual(output, test.expectedOutput) {
				t.Errorf("wanted %#v but got %#v", test.expectedOutput, output)
				return
			}
		})
	}
}
