package music

import (
	"reflect"
	"testing"

	"ifly/internal/analysis"
	"ifly/internal/semantic"
)

// TestPianoSeedChangesModesAndScore checks that the same scene can produce distinct scales and harmonies without changing its perceived theme.
func TestPianoSeedChangesModesAndScore(t *testing.T) {
	for _, scene := range []string{"A serene pond.", "A lonely figure mourning.", "A threatening storm.", "A joyful dance."} {
		modes := map[string]bool{}
		var signatures [3][]int
		mood := ""
		for variant := byte(0); variant < 3; variant++ {
			var digest [32]byte
			digest[3] = variant
			plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.5}, semantic.SemanticAnalysis{Description: scene}, digest)
			modes[plan.Mode] = true
			if mood != "" && plan.Mood != mood {
				t.Fatalf("%q changes scene theme with the tonal seed", scene)
			}
			mood = plan.Mood
			for _, note := range upperMelody(plan) {
				if note.Beat >= 32 {
					break
				}
				signatures[variant] = append(signatures[variant], note.MIDI)
			}
		}
		if len(modes) != 3 || reflect.DeepEqual(signatures[0], signatures[1]) || reflect.DeepEqual(signatures[0], signatures[2]) || reflect.DeepEqual(signatures[1], signatures[2]) {
			t.Fatalf("%q does not vary its tonal vocabulary across seeds: modes=%v", scene, modes)
		}
	}
}
