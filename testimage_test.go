package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"testing"
)

// encodeTestImage takes an in-memory image and returns the PNG data URL used by the service tests.
func encodeTestImage(t *testing.T, img image.Image) string {
	t.Helper()
	var output bytes.Buffer
	if err := png.Encode(&output, img); err != nil {
		t.Fatal(err)
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(output.Bytes())
}
