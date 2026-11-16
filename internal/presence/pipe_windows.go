package presence

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"golang.org/x/sys/windows"
)

// dialDiscord tries Discord's ten local Windows IPC pipes and returns the first open connection.
func dialDiscord() (io.ReadWriteCloser, error) {
	var firstError error
	for index := 0; index < 10; index++ {
		path := `\\?\pipe\discord-ipc-` + strconv.Itoa(index)
		name, err := windows.UTF16PtrFromString(path)
		if err != nil {
			return nil, err
		}
		handle, err := windows.CreateFile(name, windows.GENERIC_READ|windows.GENERIC_WRITE, 0, nil, windows.OPEN_EXISTING, windows.FILE_FLAG_OVERLAPPED, 0)
		if err == nil {
			return os.NewFile(uintptr(handle), path), nil
		}
		if firstError == nil {
			firstError = err
		}
	}
	if firstError == nil {
		return nil, errors.New("Discord desktop pipe is not available")
	}
	return nil, fmt.Errorf("Discord desktop pipe is not available: %w", firstError)
}
