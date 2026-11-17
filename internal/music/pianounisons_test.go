package music

import (
	"math"
	"testing"

	"ifly/internal/analysis"
	"ifly/internal/semantic"
)

// TestPianoStrikesEachKeyOnce checks that simultaneous inner chords and melody notes share one physical piano attack.
func TestPianoStrikesEachKeyOnce(t *testing.T) {
	for _, scene := range []string{"A serene pond.", "A threatening storm.", "A lively dancing scene."} {
		plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.5}, semantic.SemanticAnalysis{Description: scene}, [32]byte{})
		seen := map[pianoStrikeKey]bool{}
		for _, note := range plan.Notes {
			key := pianoStrikeKey{onset: int(math.Round(note.Beat * 20)), pitch: note.MIDI}
			if seen[key] {
				t.Fatalf("%q strikes MIDI %d twice near beat %.2f", scene, note.MIDI, note.Beat)
			}
			seen[key] = true
		}
	}
}
