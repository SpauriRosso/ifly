package music

import (
	"testing"

	"ifly/internal/analysis"
	"ifly/internal/semantic"
)

// TestPianoUsesSecondaryDominantApproaches checks that a phrase changes harmony before its dominant arrival in both piano hands.
func TestPianoUsesSecondaryDominantApproaches(t *testing.T) {
	for _, scene := range []string{"A serene pond.", "A threatening storm.", "A simple sketch."} {
		var digest [32]byte
		plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.5}, semantic.SemanticAnalysis{Description: scene}, digest)
		character := selectMusicalCharacter(semantic.SemanticAnalysis{Description: scene})
		measures := planPianoMeasures(48, plan.BeatsPerLoop/4-1, plan.Sections, digest, character)
		approaches := 0
		for bar := 0; bar+1 < len(measures); bar++ {
			measure := measures[bar]
			if measure.lateChord == nil {
				continue
			}
			approaches++
			if measure.lateChord.rootOffset != 2 || measures[bar+1].chord.rootOffset != 7 {
				t.Fatalf("%q lacks the secondary-dominant to dominant motion in bar %d", scene, bar)
			}
			bassFound := false
			for _, note := range plan.Notes {
				if note.Voice == "bass" && note.Beat == float64(bar*4+2) && note.MIDI == 38 {
					bassFound = true
					break
				}
			}
			if !bassFound {
				t.Fatalf("%q changes harmony on paper but not in the bass in bar %d", scene, bar)
			}
		}
		if approaches < 2 {
			t.Fatalf("%q never develops a mid-bar harmonic approach", scene)
		}
	}
}
