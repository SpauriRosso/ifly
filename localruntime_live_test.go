package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestLocalRuntimeLive starts and stops the installed CPU process when explicitly requested.
func TestLocalRuntimeLive(t *testing.T) {
	if os.Getenv("IFLY_VLM_RUNTIME_TEST") != "1" {
		t.Skip("set IFLY_VLM_RUNTIME_TEST=1 to test the installed CPU runtime")
	}
	if runtimeReady() {
		t.Skip("an external VLM server already owns the test port")
	}
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	runtime := &LocalRuntime{root: filepath.Join(root, "local-vlm")}
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	if err := runtime.Start(ctx); err != nil {
		t.Fatal(err)
	}
	defer runtime.Stop()
	if runtime.cmd == nil || !runtimeReady() {
		t.Fatal("managed runtime was not ready")
	}
	if err := runtime.Stop(); err != nil {
		t.Fatal(err)
	}
}
