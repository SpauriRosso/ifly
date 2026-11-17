package music

import (
	"fmt"
	"math"
	"testing"

	"ifly/internal/analysis"
	"ifly/internal/semantic"
)

// TestPianoLeftHandDevelopsAcrossPhrases checks that the accompaniment changes gesture and recedes during the reflective episode.
func TestPianoLeftHandDevelopsAcrossPhrases(t *testing.T) {
	for _, scene := range []string{"A serene pond.", "A threatening storm.", "A lively dancing scene."} {
		plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.5}, semantic.SemanticAnalysis{Description: scene}, [32]byte{})
		patterns := map[string]bool{}
		for bar := 0; bar < 12; bar++ {
			pattern := ""
			for _, note := range plan.Notes {
				if note.Voice == "arp" && int(note.Beat/4) == bar {
					pattern += fmt.Sprintf("%.2f,", math.Round((note.Beat-float64(bar*4))*4)/4)
				}
			}
			patterns[pattern] = true
		}
		if len(patterns) < 4 {
			t.Fatalf("%q repeats nearly the same left-hand timing in every bar: %v", scene, patterns)
		}
		if countVoiceInSection(plan, "arp", 48, 64) >= countVoiceInSection(plan, "arp", 32, 48) {
			t.Fatalf("%q does not make space for its reflective passage", scene)
		}
	}
}

// TestPianoLeftHandChangesAtFormalReturns checks that the bridge, its return, and the climax have distinct accompaniment rhythms.
func TestPianoLeftHandChangesAtFormalReturns(t *testing.T) {
	for _, scene := range []string{"A serene pond.", "A threatening storm.", "A lively dancing scene."} {
		plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.5}, semantic.SemanticAnalysis{Description: scene}, [32]byte{})
		patterns := map[int]string{}
		for _, episode := range []int{2, 6, 9} {
			start := float64(episode * 16)
			for _, note := range plan.Notes {
				if note.Voice == "arp" && note.Beat >= start && note.Beat < start+16 {
					patterns[episode] += fmt.Sprintf("%.2f,", math.Round((note.Beat-start)*4)/4)
				}
			}
		}
		if patterns[2] == "" || patterns[6] == "" || patterns[9] == "" || patterns[2] == patterns[6] || patterns[6] == patterns[9] || patterns[2] == patterns[9] {
			t.Fatalf("%q repeats its left-hand figure at a formal return: %v", scene, patterns)
		}
	}
}
