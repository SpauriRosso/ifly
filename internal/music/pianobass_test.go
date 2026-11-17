package music

import (
	"testing"

	"ifly/internal/analysis"
	"ifly/internal/semantic"
)

// TestPianoBassStaysLowAndHarmonic checks that both bass arrivals and responses remain below the left hand and fit the active chord.
func TestPianoBassStaysLowAndHarmonic(t *testing.T) {
	for _, scene := range []string{"A serene pond.", "A threatening storm.", "A lively dancing scene."} {
		interpretation := semantic.SemanticAnalysis{Description: scene}
		character := selectMusicalCharacter(interpretation)
		for center := byte(0); center < 12; center++ {
			var digest [32]byte
			digest[0] = center
			root := 48 + int(center)
			plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.5}, interpretation, digest)
			measures := planPianoMeasures(root, plan.BeatsPerLoop/4-1, plan.Sections, digest, character)
			for _, note := range plan.Notes {
				if note.Voice != "bass" {
					continue
				}
				if note.MIDI < 36 || note.MIDI > 50 {
					t.Fatalf("%q center %d bass crossed into the middle register: %+v", scene, center, note)
				}
				bar := int(note.Beat / 4)
				chord := measures[bar].chordAt(note.Beat - float64(bar*4))
				if !pitchBelongsToChord(note.MIDI, chordPitchClasses(root+chord.rootOffset, chord.third)) {
					t.Fatalf("%q center %d bass left its active chord: %+v", scene, center, note)
				}
			}
		}
	}
}
