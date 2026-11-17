package music

import (
	"testing"

	"ifly/internal/analysis"
	"ifly/internal/semantic"
)

// TestPianoInnerVoicesLeaveRoomForGestures checks that chord voices release between attacks and thin out in the interlude.
func TestPianoInnerVoicesLeaveRoomForGestures(t *testing.T) {
	plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.5}, semantic.SemanticAnalysis{Description: "A serene pond."}, [32]byte{})
	for _, note := range plan.Notes {
		if note.Voice == "pad" && note.Beat < float64(plan.BeatsPerLoop-4) && note.Duration > 1.6 {
			t.Fatalf("an inner chord sustains over the next gesture: %+v", note)
		}
	}
	if countVoiceInSection(plan, "pad", 48, 64) >= countVoiceInSection(plan, "pad", 32, 48) {
		t.Fatal("inner chords do not thin out during the interlude")
	}
}
