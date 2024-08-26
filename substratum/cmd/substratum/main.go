package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	ssasm "github.com/sgtcodfish/substratum/cmd/ss-asm"
	ssautotest "github.com/sgtcodfish/substratum/cmd/ss-autotest"
	sstemplate "github.com/sgtcodfish/substratum/cmd/ss-template"
	"github.com/sgtcodfish/substratum/cmd/util"
)

const (
	// ASMCommand is the command which runs the basic Substratum RISC-V "assembler"
	ASMCommand = "asm"

	// AutoTestCommand is the command which runs automated, GDB-backed tests of RISC-V baremetal programs
	AutoTestCommand = "autotest"

	// TemplateCommand creates a template of NOP commands to aid with writing a program
	TemplateCommand = "template"
)

func run(ctx context.Context) error {
	logger := util.Logger(ctx)

	if len(os.Args) < 2 {
		logger.ErrorContext(ctx, fmt.Sprintf("missing required argument: command (one of '%s')", strings.Join([]string{ASMCommand, AutoTestCommand, TemplateCommand}, ", ")))
		os.Exit(1)
	}

	subcommandName := strings.ToLower(os.Args[1])
	fullCommandName := fmt.Sprintf("%s %s", os.Args[0], subcommandName)

	switch subcommandName {
	case ASMCommand:
		err := ssasm.Invoke(ctx, fullCommandName, os.Args[2:])
		if err != nil {
			return fmt.Errorf("%s: %s", ASMCommand, err.Error())
		}

	case AutoTestCommand:
		err := ssautotest.Invoke(ctx, fullCommandName, os.Args[2:])
		if err != nil {
			return fmt.Errorf("%s: %s", AutoTestCommand, err.Error())
		}

	case TemplateCommand:
		err := sstemplate.Invoke(ctx, fullCommandName, os.Args[2:])
		if err != nil {
			return fmt.Errorf("%s: %s", TemplateCommand, err.Error())
		}

	default:
		return fmt.Errorf("unrecognised command '%s'", os.Args[1])
	}

	return nil
}

func main() {
	ctx := util.Context()

	err := run(ctx)
	if err != nil {
		util.Logger(ctx).ErrorContext(ctx, "execution failed", "error", err)
		os.Exit(1)
	}
}
