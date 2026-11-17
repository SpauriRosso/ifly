package music

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"unicode/utf8"

	"ifly/internal/analysis"
	"ifly/internal/semantic"
)

// MusicNote is one scheduled instrument event measured in beats from the start of a score.
type MusicNote struct {
	Beat     float64 `json:"beat"`
	Duration float64 `json:"duration"`
	MIDI     int     `json:"midi"`
	Velocity float64 `json:"velocity"`
	Voice    string  `json:"voice"`
}

// MusicSection marks one phrase of the generated form with its relative intensity and tempo.
type MusicSection struct {
	Label     string  `json:"label"`
	StartBeat float64 `json:"startBeat"`
	EndBeat   float64 `json:"endBeat"`
	Energy    float64 `json:"energy"`
	TempoBPM  int     `json:"tempoBpm"`
}

// MusicPlan is a bounded, repeatable score generated from one image and its optional model description.
type MusicPlan struct {
	Seed         string         `json:"seed"`
	Mood         string         `json:"mood"`
	Mode         string         `json:"mode"`
	Timbre       string         `json:"timbre"`
	TempoBPM     int            `json:"tempoBpm"`
	BeatsPerLoop int            `json:"beatsPerLoop"`
	Sections     []MusicSection `json:"sections"`
	Notes        []MusicNote    `json:"notes"`
}

type MusicService struct{}

// scoreDigestFromImageDigest takes an image digest and returns a score digest derived only from its displayed eight-byte seed.
func scoreDigestFromImageDigest(imageDigest [32]byte) [32]byte {
	return sha256.Sum256(imageDigest[:8])
}

// Compose takes an image and model result, then returns the deterministic score for its default seed and inferred theme.
func (s *MusicService) Compose(dataURL string, interpretation semantic.SemanticAnalysis) (MusicPlan, error) {
	return s.ComposeVariant(dataURL, interpretation, "", "")
}

// ComposeVariant takes an image, model result, optional 16-digit hex seed, and optional theme, then returns a repeatable piano score.
func (s *MusicService) ComposeVariant(dataURL string, interpretation semantic.SemanticAnalysis, seed, theme string) (MusicPlan, error) {
	image, err := analysis.DecodeImage(dataURL)
	if err != nil {
		return MusicPlan{}, err
	}
	if utf8.RuneCountInString(interpretation.Description) > 300 || utf8.RuneCountInString(interpretation.Atmosphere) > 120 || len(interpretation.Subjects) > 8 {
		return MusicPlan{}, errors.New("invalid semantic input")
	}
	for _, subject := range interpretation.Subjects {
		if utf8.RuneCountInString(subject) > 60 {
			return MusicPlan{}, errors.New("invalid semantic subject")
		}
	}
	if theme != "" && theme != "auto" {
		if _, ok := musicalStyles[theme]; !ok {
			return MusicPlan{}, errors.New("invalid musical theme")
		}
	}
	imageDigest := sha256.Sum256([]byte(dataURL))
	seedBytes := imageDigest[:8]
	if seed != "" {
		if len(seed) != 16 {
			return MusicPlan{}, errors.New("seed must contain 16 hexadecimal digits")
		}
		seedBytes, err = hex.DecodeString(seed)
		if err != nil {
			return MusicPlan{}, errors.New("seed must contain 16 hexadecimal digits")
		}
	}
	character := selectMusicalCharacter(interpretation)
	if theme != "" && theme != "auto" {
		character = musicalStyles[theme]
	}
	plan := composeMusicForCharacter(analysis.MeasureImage(image), interpretation, sha256.Sum256(seedBytes), character)
	plan.Seed = hex.EncodeToString(seedBytes)
	return plan, nil
}
