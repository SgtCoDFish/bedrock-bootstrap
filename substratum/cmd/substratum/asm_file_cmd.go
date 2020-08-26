package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/sgtcodfish/substratum"
)

func openFlag(flag *flag.Flag, def *os.File) (*os.File, error) {
	if flag == nil {
		return def, nil
	}

	value := flag.Value.String()

	if value == "-" {
		return def, nil
	}

	openedFile, err := os.Open(value)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", value, err)
	}

	return openedFile, nil
}

func processASMFile(flags *flag.FlagSet, logger *log.Logger) error {
	inputFile, err := openFlag(flags.Lookup("input"), os.Stdin)
	if err != nil {
		return err
	}

	outputFile, err := openFlag(flags.Lookup("output"), os.Stdout)
	if err != nil {
		return err
	}

	outputFormatFlag := flags.Lookup("output-format")
	outputFormat := strings.ToLower(outputFormatFlag.Value.String())

	if outputFormat != "hex" && outputFormat != "bin" {
		return fmt.Errorf("invalid value for 'output-format'; valid values are 'hex' and 'bin'")
	}

	inputReader := bufio.NewReader(inputFile)

	done := false
	for {
		if done {
			break
		}

		line, err := inputReader.ReadString(byte(0xa))
		if err != nil {
			if err != io.EOF {
				return err
			}

			done = true
		}

		line = strings.ToLower(strings.TrimSpace(line))

		if len(line) == 0 {
			break
		}

		parts := strings.Split(line, " ")

		insn, err := substratum.GetInstructionByName(parts[0])
		if err != nil {
			return err
		}

		if len(parts) == 1 {
			return fmt.Errorf("missing trailing args for '%s'", line)
		} else if len(parts[1:]) != insn.ArgumentCount() {
			return fmt.Errorf("invalid number of arguments for '%s'", insn.Name)
		}

		for i, s := range parts[1:] {
			parts[i+1] = strings.TrimRight(s, ",\n")
		}

		assembled, err := insn.AssembleRaw(parts[1:])
		if err != nil {
			return fmt.Errorf("failed to assemble '%s': %w", line, err)
		}

		if outputFormat == "bin" {
			_, _ = outputFile.Write(assembled)
		} else {
			_, _ = outputFile.WriteString(fmt.Sprintf("%02x %02x %02x %02x  # %s\n", assembled[0], assembled[1], assembled[2], assembled[3], line))
		}
	}

	return nil
}
