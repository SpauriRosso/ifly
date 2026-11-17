package music

import (
	"reflect"
	"testing"

	"ifly/internal/analysis"
	"ifly/internal/semantic"
)

// leadOnsetsInEpisode takes a score and episode index, then returns lead attacks relative to that episode's first beat.
func leadOnsetsInEpisode(plan MusicPlan, episode int) []float64 {
	onsets := []float64{}
	start := float64(episode * 16)
	for _, note := range plan.Notes {
		if note.Voice == "lead" && note.Beat >= start && note.Beat < start+16 {
			onsets = append(onsets, note.Beat-start)
		}
	}
	return onsets
}

// TestPianoReturnsDevelopTheirRhythm checks that the second bridge, climax, and echo have audible phrase pacing distinct from their earlier material.
func TestPianoReturnsDevelopTheirRhythm(t *testing.T) {
	for _, scene := range []string{"A serene pond with water lilies.", "A threatening storm.", "A lively dancing scene."} {
		plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.5}, semantic.SemanticAnalysis{Description: scene}, [32]byte{})
		for _, pair := range [][2]int{{1, 5}, {2, 6}, {6, 9}, {0, 12}} {
			first := leadOnsetsInEpisode(plan, pair[0])
			later := leadOnsetsInEpisode(plan, pair[1])
			if len(first) == 0 || len(later) == 0 || reflect.DeepEqual(first, later) {
				t.Fatalf("%q repeats lead rhythm in episodes %d and %d: %v / %v", scene, pair[0], pair[1], first, later)
			}
		}
	}
}
