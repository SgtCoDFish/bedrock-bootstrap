package substratum

import (
	"fmt"
	"strings"
)

// AssembleLine takes a newline (\n) delimited line of text and attempts to translate it
// into machine code. This isn't a full assembler; there's no support for jumps or other
// luxuries. It's intended to be an extremely simple way of going from a list of raw commands
// to a raw RISC-V binary representation of those commands.
func AssembleLine(line string) ([]byte, error) {
	line = strings.ToLower(strings.TrimSpace(strings.Split(line, "\n")[0]))

	// check for trailing comments on the line
	commentSplit := strings.Split(line, "#")

	if len(commentSplit) > 1 {
		line = strings.TrimSpace(commentSplit[0])
	} else if strings.HasPrefix(line, "#") {
		// don't miss the case where the entire line is a comment
		return nil, nil
	}

	if len(line) == 0 {
		return nil, nil
	}

	parts := strings.Split(line, " ")

	insn, err := GetInstructionByName(parts[0])
	if err != nil {
		return nil, err
	}

	if len(parts) == 1 {
		return nil, fmt.Errorf("missing trailing args for '%s'", line)
	}

	rest := parts[1:]

	for i, s := range rest {
		rest[i] = strings.Trim(s, ",\n")
	}

	if len(rest) != insn.ArgumentCount() {
		if (insn.Type == IType || insn.Type == SType) && len(rest) == 2 {
			// handle instructions of the format "lw xXX, imm(xXX)"
			offset := strings.TrimRight(rest[1], ")") // of the format imm(xXX
			offsetParts := strings.Split(offset, "(")

			if len(offsetParts) != 2 {
				return nil, fmt.Errorf("malformed instruction: %s", line)
			}

			rest[1] = offsetParts[1]            // the register in the parentheses
			rest = append(rest, offsetParts[0]) // the immediate value
		} else {
			return nil, fmt.Errorf("invalid number of arguments for '%s'", insn.Name)
		}
	}

	assembled, err := insn.AssembleRaw(rest)
	if err != nil {
		return nil, fmt.Errorf("failed to assemble '%s': %w", line, err)
	}

	return assembled, nil
}
