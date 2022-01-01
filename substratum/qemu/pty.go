package qemu

import (
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

// These functions are based on goterm:
// https://github.com/google/goterm/blob/master/example/example.go

const (
	// tiocgptn is used to identify the sub pty ID
	tiocgptn = 0x80045430

	// tiocsptlck is used to unlock the sub pty
	tiocsptlck = 0x40045431
)

// PTY wraps a master file for a pseudo tty and a filename
type PTY struct {
	Master      *os.File
	SubFilename string
}

// NewPTY allocates a new PTY, unlocks the subordinate TTY and returns a PTY
// structure containing a file referring to the master and the sub filename.
// Should be closed when finished with
func NewPTY() (*PTY, error) {
	masterFile, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	err = unlockSubPty(masterFile)
	if err != nil {
		_ = masterFile.Close()
		return nil, err
	}

	subPtyName, err := getSubPtyName(masterFile)
	if err != nil {
		_ = masterFile.Close()
		return nil, err
	}

	return &PTY{
		Master:      masterFile,
		SubFilename: subPtyName,
	}, nil
}

// Close closes both the master and subordinate files
func (p *PTY) Close() error {
	if err := p.Master.Close(); err != nil {
		return err
	}

	return nil
}

// getSubPtyName returns the filename of the allocated sub pty
func getSubPtyName(masterFile *os.File) (string, error) {
	var ptyno uint64

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(masterFile.Fd()), uintptr(tiocgptn), uintptr(unsafe.Pointer(&ptyno)))
	if errno != 0 {
		return "", errno
	}

	name := "/dev/pts/" + strconv.FormatUint(ptyno, 10)
	return name, nil
}

func unlockSubPty(masterFile *os.File) error {
	// 0 means "unlock"
	unlock := 0

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(masterFile.Fd()), uintptr(tiocsptlck), uintptr(unsafe.Pointer(&unlock)))

	if errno != 0 {
		return errno
	}

	return nil
}
