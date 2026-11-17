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

// TestDevelopedPianoPhrasesKeepPitchVariety checks that transformed motives remain playable phrases rather than collapsing into repeated attacks.
func TestDevelopedPianoPhrasesKeepPitchVariety(t *testing.T) {
	for _, scene := range []string{"A serene pond with water lilies.", "A threatening storm.", "A lively dancing scene."} {
		for center := byte(0); center < 12; center++ {
			var digest [32]byte
			digest[0] = center
			plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.5}, semantic.SemanticAnalysis{Description: scene}, digest)
			for _, episode := range []int{5, 6, 9, 12} {
				counts := map[int]int{}
				attacks := 0
				most := 0
				start := float64(episode * 16)
				for _, note := range plan.Notes {
					if note.Voice != "lead" || note.Beat < start || note.Beat >= start+16 {
						continue
					}
					attacks++
					counts[note.MIDI]++
					if counts[note.MIDI] > most {
						most = counts[note.MIDI]
					}
				}
				if len(counts) < 5 || most*5 > attacks*2 {
					t.Fatalf("%q center %d episode %d has a flat lead: %d pitches, most frequent %d/%d", scene, center, episode, len(counts), most, attacks)
				}
			}
		}
	}
}
