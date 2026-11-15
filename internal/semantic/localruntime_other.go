//go:build !windows

package semantic

import "os/exec"

// hideProcessWindow does nothing on platforms without a Windows console window.
func hideProcessWindow(_ *exec.Cmd) {}
