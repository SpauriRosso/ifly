//go:build windows

package semantic

import (
	"os/exec"
	"syscall"
)

// hideProcessWindow prevents the managed CPU runtime from opening a second console on Windows.
func hideProcessWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
