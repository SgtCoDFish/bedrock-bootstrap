package autotest

import (
	"context"
	"debug/elf"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sgtcodfish/substratum/cmd/util"
)

const (
	testNameFlagName = "test-name"
	kernelFlagName   = "kernel"
)

type TestFunc func()

var testMap = map[string]TestFunc{}

// Invocation holds state for a given invocation of the ss-autotest command
type Invocation struct {
	gdbPath string
	gdbPort string

	qemuPath string

	testName   string
	kernelFile string
}

// ParseInvocation builds an invocation for a run of the the ss-autotest command, where name is the name
// of the binary being run and flags are a slice of command-line flags to be parsed
func ParseInvocation(name string, flags []string) (*Invocation, error) {
	autoTestCmd := flag.NewFlagSet(name, flag.ExitOnError)

	allTestNames := getAllTestNames()

	gdbPathFlag := autoTestCmd.String("gdb", "", "Path to the GDB executable to use. Defaults to an architecture approprite value on $PATH, or else ${RISCV_PREFIX}gdb")
	qemuPathFlag := autoTestCmd.String("qemu", "", "Path to the QEMU executable to run. Defaults to an architecture appropriate system on $PATH if possible")
	kernelFlag := autoTestCmd.String(kernelFlagName, "", "ELF file containing the kernel to run using QEMU")
	testNameFlag := autoTestCmd.String(testNameFlagName, "", "Name of the test to run. Must be one of: "+allTestNames)

	if err := autoTestCmd.Parse(flags); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	if testNameFlag == nil || *testNameFlag == "" {
		return nil, fmt.Errorf("missing required flag: %s", testNameFlagName)
	}

	testName := strings.ToLower(*testNameFlag)

	if kernelFlag == nil || *kernelFlag == "" {
		return nil, fmt.Errorf("missing required flag: %s", kernelFlagName)
	}

	bitSize, err := checkKernel(*kernelFlag)
	if err != nil {
		return nil, err
	}

	gdbPath := *gdbPathFlag

	if len(gdbPath) == 0 {
		var ok bool

		gdbPath, ok = findDefaultGDB(bitSize)
		if !ok {
			return nil, fmt.Errorf("failed to find any valid GDB executable")
		}
	}

	qemuPath := *qemuPathFlag

	if len(qemuPath) == 0 {
		var ok bool

		qemuPath, ok = findDefaultQEMU(bitSize)
		if !ok {
			return nil, fmt.Errorf("failed to find any valid QEMU executable")
		}
	}

	if _, ok := testMap[testName]; !ok {
		return nil, fmt.Errorf("invalid test name %q; must be one of: %s", testName, allTestNames)
	}

	return &Invocation{
		gdbPath:    gdbPath,
		gdbPort:    ":1234",
		qemuPath:   qemuPath,
		testName:   testName,
		kernelFile: *kernelFlag,
	}, nil
}

func findDefaultGDB(bitSize int) (string, bool) {
	path, err := exec.LookPath(fmt.Sprintf("riscv%d-elf-gdb", bitSize))
	if err == nil {
		return path, true
	}

	prefixPath := os.Getenv("RISCV_PREFIX") + "gdb"
	_, err = os.Stat(prefixPath)
	if err == nil {
		return prefixPath, true
	}

	return "", false
}

func findDefaultQEMU(bitSize int) (string, bool) {
	expectedBinary := fmt.Sprintf("qemu-system-riscv%d", bitSize)

	path, err := exec.LookPath(expectedBinary)
	if err == nil {
		return path, true
	}

	prefixPath := os.Getenv("RISCV_PREFIX") + expectedBinary
	_, err = os.Stat(prefixPath)
	if err == nil {
		return prefixPath, true
	}

	return "", false
}

// Invoke parses and runs an invocation of ss-autotest, where name is the name
// of the binary being run and flags are a slice of command-line flags to be parsed
func Invoke(ctx context.Context, name string, flags []string) error {
	invocation, err := ParseInvocation(name, flags)
	if err != nil {
		return err
	}

	return invocation.Run(ctx)
}

// Run executes ss-autotest for the configured invocation
func (a *Invocation) Run(ctx context.Context) error {
	logger := util.Logger(ctx)

	logger.InfoContext(ctx, "processing autotest", "testName", a.testName, "gdb", a.gdbPath, "qemu", a.qemuPath)

	_, ok := testMap[a.testName]
	if !ok {
		panic("invalid test name when running autotest invocation")
	}

	return fmt.Errorf("NYI")
}

func getAllTestNames() string {
	allTests := make([]string, 0, len(testMap))

	for k := range testMap {
		allTests = append(allTests, k)
	}

	return strings.Join(allTests, " | ")
}

// checkKernel ensures the kernel exists and is a valid RISC-V ELF file, and returns the bit size (32 or 64)
func checkKernel(path string) (int, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, fmt.Errorf("kernel %q not found", path)
		}

		return 0, err
	}

	f, err := elf.Open(path)
	if err != nil {
		return 0, fmt.Errorf("failed to open kernel as an ELF file: %s", err)
	}

	defer f.Close()

	if f.FileHeader.Machine != elf.EM_RISCV {
		return 0, fmt.Errorf("specified kernel is not RISC-V")
	}

	class := f.FileHeader.Class

	if class == elf.ELFCLASS32 {
		return 32, nil
	} else if class == elf.ELFCLASS64 {
		return 64, nil
	}

	return 0, fmt.Errorf("unknown architecture in kernel; expected either 32 or 64 bit")
}
