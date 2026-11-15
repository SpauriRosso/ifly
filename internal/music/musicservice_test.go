package music

import (
	"image"
	"reflect"
	"testing"

	"ifly/internal/analysis"
	"ifly/internal/semantic"
)

// countLeadNotes takes a score and returns the number of melodic events used to compare visual densities.
func countLeadNotes(plan MusicPlan) int {
	count := 0
	for _, note := range plan.Notes {
		if note.Voice == "lead" {
			count++
		}
	}
	return count
}

// TestMusicServiceCompose checks that one artwork has a repeatable score and opposing scenes produce different scores.
func TestMusicServiceCompose(t *testing.T) {
	dataURL := encodeTestImage(t, image.NewRGBA(image.Rect(0, 0, 16, 16)))
	service := &MusicService{}
	calm := semantic.SemanticAnalysis{Description: "A serene landscape at sunset with a deer."}
	first, err := service.Compose(dataURL, calm)
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Compose(dataURL, calm)
	if err != nil || !reflect.DeepEqual(first, second) {
		t.Fatalf("same artwork changed score: %v", err)
	}
	if first.Mood != "calm" || first.TempoBPM < 55 || first.BeatsPerLoop != 64 || len(first.Sections) != 4 || len(first.Notes) < 80 {
		t.Fatalf("invalid calm score: mood=%s tempo=%d beats=%d sections=%d notes=%d", first.Mood, first.TempoBPM, first.BeatsPerLoop, len(first.Sections), len(first.Notes))
	}
	threat, err := service.Compose(dataURL, semantic.SemanticAnalysis{Description: "A threatening violent scene."})
	if err != nil {
		t.Fatal(err)
	}
	if threat.Mood != "tense" || threat.Mode == first.Mode || threat.TempoBPM == first.TempoBPM || reflect.DeepEqual(threat.Notes, first.Notes) {
		t.Fatal("opposing scene produced the same score")
	}
	for index, note := range first.Notes {
		if note.Beat < 0 || note.Beat >= float64(first.BeatsPerLoop) || note.Duration <= 0 || note.MIDI < 36 || note.MIDI > 90 || note.Velocity <= 0 || note.Velocity > 1 {
			t.Fatalf("unsafe note: %+v", note)
		}
		if index > 0 && note.Beat < first.Notes[index-1].Beat {
			t.Fatal("score events are not sorted")
		}
		if note.Voice == "lead" && note.Beat > 60 {
			t.Fatal("the concluding bar should leave room for silence")
		}
	}
}

// TestComposeMusicUsesVisualComplexity checks that added visual detail changes melodic density for the same interpreted scene.
func TestComposeMusicUsesVisualComplexity(t *testing.T) {
	interpretation := semantic.SemanticAnalysis{Description: "A serene landscape."}
	quiet := composeMusic(analysis.StructuralAnalysis{Complexity: 0.1}, interpretation, [32]byte{})
	busy := composeMusic(analysis.StructuralAnalysis{Complexity: 0.8}, interpretation, [32]byte{})
	if countLeadNotes(busy) <= countLeadNotes(quiet) || busy.TempoBPM <= quiet.TempoBPM {
		t.Fatal("visual complexity did not alter density and pace")
	}
}

// TestComposeMusicResolvesHome checks that a tense phrase ends on its tonic instead of looping from an unresolved dominant.
func TestComposeMusicResolvesHome(t *testing.T) {
	plan := composeMusic(analysis.StructuralAnalysis{}, semantic.SemanticAnalysis{Description: "A threatening scene."}, [32]byte{})
	for _, note := range plan.Notes {
		if note.Beat == 60 && note.Voice == "pad" && note.MIDI%12 == 0 {
			return
		}
	}
	t.Fatal("final harmony did not resolve to the tonic")
}

// TestVoiceChordPrefersNearInversion checks that chord voices move less than a root-position block change.
func TestVoiceChordPrefersNearInversion(t *testing.T) {
	first := voiceChord(48, 4, [3]int{}, false)
	second := voiceChord(53, 4, first, true)
	motion := 0
	for note := range first {
		delta := first[note] - second[note]
		if delta < 0 {
			delta = -delta
		}
		motion += delta
	}
	if motion >= 15 {
		t.Fatalf("chord motion is too large: %v to %v", first, second)
	}
}

// TestMusicServiceComposeRejectsInvalidImage checks that the score cannot be generated from an unbounded or malformed image.
func TestMusicServiceComposeRejectsInvalidImage(t *testing.T) {
	if _, err := (&MusicService{}).Compose("data:image/png;base64,AA==", semantic.SemanticAnalysis{}); err == nil {
		t.Fatal("expected malformed image to be rejected")
	}
}
