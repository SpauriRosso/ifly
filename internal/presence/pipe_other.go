//go:build !windows

package presence

import (
	"errors"
	"io"
)

// dialDiscord reports that the Windows IPC transport is unavailable on this platform.
func dialDiscord() (io.ReadWriteCloser, error) {
	return nil, errors.New("Discord presence currently supports Windows only")
}
