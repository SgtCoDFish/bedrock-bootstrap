package substratum

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	nopString = "13 00 00 00 # nop\n"
)

type Invocation struct {
	Size int

	Output io.StringWriter
}

func (inv *Invocation) Run(ctx context.Context) error {
	buf := &strings.Builder{}

	expectedCap := len(nopString) * inv.Size
	expectedCap += inv.Size / 4

	buf.Grow(expectedCap)

	for i := 0; i < inv.Size; i++ {
		if i > 0 && i%4 == 0 {
			_, _ = buf.WriteString("\n")
		}

		_, _ = buf.WriteString(nopString)
	}

	_, err := inv.Output.WriteString(buf.String())
	if err != nil {
		return err
	}

	return nil
}

func Invoke(ctx context.Context, fullCommandName string, flags []string) error {
	inv, err := ParseInvocation(fullCommandName, flags)
	if err != nil {
		return err
	}

	return inv.Run(ctx)
}

func ParseInvocation(name string, flags []string) (*Invocation, error) {
	templateCommand := flag.NewFlagSet(name, flag.ExitOnError)

	size := templateCommand.Int("size", 1024, "Size of template to generate, in number of NOP instructions")

	if err := templateCommand.Parse(flags); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %s", err)
	}

	if *size%4 != 0 || *size == 0 {
		return nil, fmt.Errorf("size must be a multiple of 4 and greater than 0")
	}

	return &Invocation{
		Size: *size,

		Output: os.Stdout,
	}, nil
}
