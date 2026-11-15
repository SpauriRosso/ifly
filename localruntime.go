package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type LocalRuntime struct {
	root    string
	cmd     *exec.Cmd
	done    chan error
	logFile *os.File
}

// findLocalRuntime locates verified development assets next to the executable or its parent directory.
func findLocalRuntime() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(exe)
	for _, root := range []string{filepath.Join(dir, "local-vlm"), filepath.Join(dir, "..", "local-vlm")} {
		if _, err := os.Stat(filepath.Join(root, "runtime", "llama-server.exe")); err == nil {
			return root, nil
		}
	}
	return "", errors.New("local VLM files are missing; run scripts/setup-local-vlm.ps1")
}

// runtimeReady reports whether a loopback llama.cpp server is accepting requests on the development port.
func runtimeReady() bool {
	client := &http.Client{Timeout: 700 * time.Millisecond}
	response, err := client.Get("http://127.0.0.1:11480/health")
	if err != nil {
		return false
	}
	defer response.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 1024))
	return response.StatusCode == http.StatusOK
}

// Start launches an installed CPU model on loopback and waits for readiness without failing the desktop shell.
func (r *LocalRuntime) Start(ctx context.Context) error {
	if runtimeReady() {
		return nil
	}
	root := r.root
	if root == "" {
		var err error
		root, err = findLocalRuntime()
		if err != nil {
			return err
		}
	}
	model := filepath.Join(root, "SmolVLM2-2.2B-Instruct-Q4_K_M.gguf")
	projector := filepath.Join(root, "mmproj-SmolVLM2-2.2B-Instruct-Q8_0.gguf")
	for _, path := range []string{model, projector} {
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("local VLM file missing: %s", filepath.Base(path))
		}
	}
	logFile, err := os.OpenFile(filepath.Join(root, "server.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	r.logFile = logFile
	r.cmd = exec.CommandContext(ctx, filepath.Join(root, "runtime", "llama-server.exe"),
		"--model", model, "--mmproj", projector, "--no-mmproj-offload", "--n-gpu-layers", "0",
		"--host", "127.0.0.1", "--port", "11480", "--ctx-size", "2048", "--cors-origins", "localhost")
	r.cmd.Stdout = logFile
	r.cmd.Stderr = logFile
	hideProcessWindow(r.cmd)
	if err := r.cmd.Start(); err != nil {
		_ = logFile.Close()
		r.logFile = nil
		r.cmd = nil
		return err
	}
	r.done = make(chan error, 1)
	cmd := r.cmd
	done := r.done
	go func() { done <- cmd.Wait() }()
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	timer := time.NewTimer(30 * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-ticker.C:
			if runtimeReady() {
				return nil
			}
		case err := <-r.done:
			_ = r.logFile.Close()
			r.logFile = nil
			r.cmd = nil
			return fmt.Errorf("local VLM exited during startup: %v", err)
		case <-timer.C:
			_ = r.Stop()
			return errors.New("local VLM startup timed out")
		case <-ctx.Done():
			_ = r.Stop()
			return ctx.Err()
		}
	}
}

// Stop terminates only the inference process started by this application and closes its log.
func (r *LocalRuntime) Stop() error {
	if r.cmd == nil {
		return nil
	}
	_ = r.cmd.Process.Kill()
	select {
	case <-r.done:
	case <-time.After(3 * time.Second):
	}
	r.cmd = nil
	if r.logFile != nil {
		_ = r.logFile.Close()
		r.logFile = nil
	}
	return nil
}
