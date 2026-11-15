package analysis

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"testing"
)

// encodeTestImage serializes an in-memory image as the same PNG data URL sent by the frontend.
func encodeTestImage(t *testing.T, img image.Image) string {
	t.Helper()
	var output bytes.Buffer
	if err := png.Encode(&output, img); err != nil {
		t.Fatal(err)
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(output.Bytes())
}

// TestAnalyzeDistinguishesStructure verifies that flat and striped images produce different measured features.
func TestAnalyzeDistinguishesStructure(t *testing.T) {
	flat := image.NewRGBA(image.Rect(0, 0, 32, 32))
	striped := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			flat.Set(x, y, color.White)
			if x%2 == 0 {
				striped.Set(x, y, color.White)
			} else {
				striped.Set(x, y, color.Black)
			}
		}
	}
	service := &AnalysisService{}
	flatResult, err := service.Analyze(encodeTestImage(t, flat))
	if err != nil {
		t.Fatal(err)
	}
	stripedResult, err := service.Analyze(encodeTestImage(t, striped))
	if err != nil {
		t.Fatal(err)
	}
	if stripedResult.EdgeDensity <= flatResult.EdgeDensity || stripedResult.Contrast <= flatResult.Contrast || stripedResult.Complexity <= flatResult.Complexity {
		t.Fatalf("striped image should have more structure: flat=%+v striped=%+v", flatResult, stripedResult)
	}
}

// TestAnalyzeRejectsInvalidInput checks that non-image data never reaches the pixel analyzer.
func TestAnalyzeRejectsInvalidInput(t *testing.T) {
	if _, err := (&AnalysisService{}).Analyze("data:image/png;base64,Zm9v"); err == nil {
		t.Fatal("expected invalid PNG to be rejected")
	}
}
