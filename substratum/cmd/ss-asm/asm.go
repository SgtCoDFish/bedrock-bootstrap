package asm

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sgtcodfish/substratum"
)

type outputFormatter func(cmd string, out []byte) []byte

func formatHex(cmd string, out []byte) []byte {
	return []byte(fmt.Sprintf("%02x %02x %02x %02x  # %s", out[0], out[1], out[2], out[3], cmd))
}

func formatBin(_ string, out []byte) []byte {
	return out
}

var outputFormatterMap map[string]outputFormatter = map[string]outputFormatter{
	"hex": formatHex,
	"bin": formatBin,
}

// Invocation holds state for a given invocation of the ss-asm command
type Invocation struct {
	Input  *bufio.Reader
	Output io.WriteCloser

	OutputFormatter outputFormatter

	underlyingInput io.Closer
}

// Close closes both input and output
func (a *Invocation) Close() error {
	inputErr := a.underlyingInput.Close()
	outputErr := a.Output.Close()

	if inputErr != nil && outputErr == nil {
		return inputErr
	} else if inputErr == nil && outputErr == nil {
		return outputErr
	} else if inputErr != nil && outputErr != nil {
		return fmt.Errorf(`failed to close input ("%v") and output ("%v")`, inputErr, outputErr)
	}

	return nil
}

// Invoke parses and runs an invocation of ss-asm, where name is the name
// of the binary being run and flags are a slice of command-line flags to be parsed
func Invoke(ctx context.Context, name string, flags []string) error {
	invocation, err := ParseInvocation(name, flags)
	if err != nil {
		return err
	}

	return invocation.Run(ctx)
}

// ParseInvocation builds an invocation for a run of the the ss-asm command, where name is the name
// of the binary being run and flags are a slice of command-line flags to be parsed
func ParseInvocation(name string, flags []string) (*Invocation, error) {
	asmCMD := flag.NewFlagSet(name, flag.ExitOnError)

	inputFilename := asmCMD.String("input", "-", "File to read, defaults to stdin")
	outputFilename := asmCMD.String("output", "-", "File to write, defaults to stdout")
	outputFormatFlag := asmCMD.String("output-format", "bin", "The format in which to write output - 'bin' for binary, 'hex' for ASCII hex")

	if err := asmCMD.Parse(flags); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	inputFile, err := openFlag(inputFilename, os.Stdin)
	if err != nil {
		return nil, err
	}

	inputReader := bufio.NewReader(inputFile)

	outputFile, err := openFlag(outputFilename, os.Stdout)
	if err != nil {
		return nil, err
	}

	outputFormat := strings.ToLower(*outputFormatFlag)

	formatterFunc, ok := outputFormatterMap[outputFormat]
	if !ok {
		return nil, fmt.Errorf("invalid value for 'output-format'; valid values are 'hex' and 'bin'")
	}

	return &Invocation{
		Input:           inputReader,
		Output:          outputFile,
		OutputFormatter: formatterFunc,

		underlyingInput: inputFile,
	}, nil
}

// Run executes ss-asm for the configured invocation
func (a *Invocation) Run(ctx context.Context) error {
	defer func() {
		err := a.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to close files: %v", err)
		}
	}()

	done := false
	for {
		if done {
			break
		}

		line, err := a.Input.ReadString(byte(0xa))
		if err != nil {
			if err != io.EOF {
				return err
			}

			done = true
			continue
		}

		assembled, err := substratum.AssembleLine(line)
		if err != nil {
			return fmt.Errorf("failed to assemble line: %w", err)
		}

		if assembled == nil {
			// line was a comment or was empty
			continue
		}

		insnOut := a.OutputFormatter(line, assembled)

		_, err = a.Output.Write(insnOut)
		if err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
	}

	return nil
}

func openFlag(flag *string, def *os.File) (*os.File, error) {
	if flag == nil || *flag == "-" {
		return def, nil
	}

	value := *flag

	openedFile, err := os.Open(value)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", value, err)
	}

	return openedFile, nil
}
