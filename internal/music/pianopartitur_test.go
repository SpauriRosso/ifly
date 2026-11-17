package music

import (
	"math"
	"reflect"
	"testing"

	"ifly/internal/analysis"
	"ifly/internal/semantic"
)

// TestPianoWrittenOpeningResolvesModalNeighbor checks that the opening figure visits the mode's lower neighbor and returns.
func TestPianoWrittenOpeningResolvesModalNeighbor(t *testing.T) {
	for center := byte(0); center < 12; center++ {
		var digest [32]byte
		digest[0] = center
		plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.1}, semantic.SemanticAnalysis{Description: "A serene pond."}, digest)
		lead := upperMelody(plan)
		if len(lead) < 3 {
			t.Fatalf("center %d has too few written notes: %+v", center, lead)
		}
		if lead[0].Beat > 0.05 || math.Abs(lead[1].Beat-0.5) > 0.05 || math.Abs(lead[2].Beat-1) > 0.05 || lead[0].MIDI-lead[1].MIDI < 1 || lead[0].MIDI-lead[1].MIDI > 2 || lead[2].MIDI != lead[0].MIDI {
			t.Fatalf("center %d lost its modal neighbor and resolution: %+v", center, lead[:3])
		}
	}
}

// TestPianoThemesVaryByArtwork checks that three image seeds select distinct opening phrases in the same key and mood.
func TestPianoThemesVaryByArtwork(t *testing.T) {
	type attack struct {
		beat  float64
		pitch int
	}
	var themes [3][]attack
	for variant := byte(0); variant < 3; variant++ {
		var digest [32]byte
		digest[1] = variant
		plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.5}, semantic.SemanticAnalysis{Description: "A serene pond."}, digest)
		for _, note := range upperMelody(plan) {
			if note.Beat >= 8 {
				break
			}
			themes[variant] = append(themes[variant], attack{beat: math.Round(note.Beat*4) / 4, pitch: note.MIDI})
		}
	}
	for first := 0; first < len(themes); first++ {
		for second := first + 1; second < len(themes); second++ {
			if reflect.DeepEqual(themes[first], themes[second]) {
				t.Fatalf("artwork variants %d and %d repeat the same opening", first, second)
			}
		}
	}
}

// TestPianoScoreStaysConnected checks that written notes sustain every formal transition without an unintended empty measure.
func TestPianoScoreStaysConnected(t *testing.T) {
	for _, scene := range []string{"A serene pond.", "A threatening storm.", "A lively dancing scene."} {
		plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.5}, semantic.SemanticAnalysis{Description: scene}, [32]byte{})
		activeUntil := 0.0
		for _, note := range plan.Notes {
			if note.Beat >= float64(plan.BeatsPerLoop-4) {
				break
			}
			if note.Beat > activeUntil+0.45 {
				t.Fatalf("%q has an unplanned %.2f-beat gap before beat %.2f", scene, note.Beat-activeUntil, note.Beat)
			}
			if note.Beat+note.Duration > activeUntil {
				activeUntil = note.Beat + note.Duration
			}
		}
	}
}

// TestPianoScoreAvoidsUnplayableLeadLeaps checks connected right-hand writing across scene styles and every tonal center.
func TestPianoScoreAvoidsUnplayableLeadLeaps(t *testing.T) {
	for _, scene := range []string{"A serene pond.", "A threatening storm.", "A lively dancing scene."} {
		for center := byte(0); center < 12; center++ {
			var digest [32]byte
			digest[0] = center
			plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.5}, semantic.SemanticAnalysis{Description: scene}, digest)
			lead := upperMelody(plan)
			for index := 1; index < len(lead); index++ {
				if absInt(lead[index].MIDI-lead[index-1].MIDI) > 12 {
					t.Fatalf("%q center %d jumps from %d to %d at beat %.2f", scene, center, lead[index-1].MIDI, lead[index].MIDI, lead[index].Beat)
				}
			}
		}
	}
}
