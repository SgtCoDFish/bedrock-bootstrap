//go:build darwin

package qemu

import (
	"fmt"
	"os"
)

// unlockSubPty does nothing on macOS; PTYs are unlocked automatically.
func unlockSubPty(_ *os.File) error {
	// TMP: return an error because the whole setup needs to be different on macos
	return fmt.Errorf("NYI")
}
