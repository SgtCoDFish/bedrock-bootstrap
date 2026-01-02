//go:build linux

package qemu

import (
	"os"
	"syscall"
	"unsafe"
)

func unlockSubPty(masterFile *os.File) error {
	// 0 means "unlock"
	unlock := 0

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(masterFile.Fd()), uintptr(tiocsptlck), uintptr(unsafe.Pointer(&unlock)))

	if errno != 0 {
		return errno
	}

	return nil
}
