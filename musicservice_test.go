package main

import (
	"image"
	"reflect"
	"testing"

	"changeme/internal/semantic"
)

// TestMusicServiceCompose checks that one artwork has a repeatable score and opposing scenes produce different scores.
func TestMusicServiceCompose(t *testing.T) {
	dataURL := encodeTestImage(t, image.NewRGBA(image.Rect(0, 0, 16, 16)))
	service := &MusicService{}
	calm := semantic.SemanticAnalysis{Description: "A serene landscape at sunset with a deer."}
	first, err := service.Compose(dataURL, calm)
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Compose(dataURL, calm)
	if err != nil || !reflect.DeepEqual(first, second) {
		t.Fatalf("same artwork changed score: %v", err)
	}
	if first.Mood != "calm" || first.TempoBPM < 55 || first.BeatsPerLoop != 32 || len(first.Notes) < 40 {
		t.Fatalf("invalid calm score: %+v", first)
	}
	threat, err := service.Compose(dataURL, semantic.SemanticAnalysis{Description: "A threatening violent scene."})
	if err != nil {
		t.Fatal(err)
	}
	if threat.Mood != "tense" || threat.Mode == first.Mode || threat.TempoBPM == first.TempoBPM || reflect.DeepEqual(threat.Notes, first.Notes) {
		t.Fatal("opposing scene produced the same score")
	}
	for _, note := range first.Notes {
		if note.Beat < 0 || note.Beat >= float64(first.BeatsPerLoop) || note.Duration <= 0 || note.MIDI < 40 || note.MIDI > 84 || note.Velocity <= 0 || note.Velocity > 1 {
			t.Fatalf("unsafe note: %+v", note)
		}
	}
}

// TestMusicServiceComposeRejectsInvalidImage checks that the score cannot be generated from an unbounded or malformed image.
func TestMusicServiceComposeRejectsInvalidImage(t *testing.T) {
	if _, err := (&MusicService{}).Compose("data:image/png;base64,AA==", semantic.SemanticAnalysis{}); err == nil {
		t.Fatal("expected malformed image to be rejected")
	}
}
