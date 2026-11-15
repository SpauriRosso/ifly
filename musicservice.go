package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"unicode/utf8"

	"changeme/internal/analysis"
	"changeme/internal/semantic"
)

// MusicNote is one scheduled synthesizer event measured in beats from the start of a loop.
type MusicNote struct {
	Beat     float64 `json:"beat"`
	Duration float64 `json:"duration"`
	MIDI     int     `json:"midi"`
	Velocity float64 `json:"velocity"`
	Voice    string  `json:"voice"`
}

// MusicPlan is a bounded, repeatable score generated from one image and its optional model description.
type MusicPlan struct {
	Seed         string      `json:"seed"`
	Mood         string      `json:"mood"`
	Mode         string      `json:"mode"`
	TempoBPM     int         `json:"tempoBpm"`
	BeatsPerLoop int         `json:"beatsPerLoop"`
	Notes        []MusicNote `json:"notes"`
}

type MusicService struct{}

// Compose validates a bounded image and model result, then returns a deterministic score tied to the image bytes.
func (s *MusicService) Compose(dataURL string, interpretation semantic.SemanticAnalysis) (MusicPlan, error) {
	image, err := analysis.DecodeImage(dataURL)
	if err != nil {
		return MusicPlan{}, err
	}
	if utf8.RuneCountInString(interpretation.Description) > 300 || len(interpretation.Subjects) > 8 {
		return MusicPlan{}, errors.New("invalid semantic input")
	}
	for _, subject := range interpretation.Subjects {
		if utf8.RuneCountInString(subject) > 60 {
			return MusicPlan{}, errors.New("invalid semantic subject")
		}
	}
	digest := sha256.Sum256([]byte(dataURL))
	plan := composeMusic(analysis.MeasureImage(image), interpretation, digest)
	plan.Seed = hex.EncodeToString(digest[:8])
	return plan, nil
}
