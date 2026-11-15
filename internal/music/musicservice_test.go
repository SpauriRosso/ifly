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

// countVoiceInSection takes a score, voice, and beat range and returns the matching event count for arrangement checks.
func countVoiceInSection(plan MusicPlan, voice string, start, end float64) int {
	count := 0
	for _, note := range plan.Notes {
		if note.Voice == voice && note.Beat >= start && note.Beat < end {
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

// TestArpeggioBuildsAndRecedes checks that the counterline enters after the opening and withdraws before the final cadence.
func TestArpeggioBuildsAndRecedes(t *testing.T) {
	calm := composeMusic(analysis.StructuralAnalysis{}, semantic.SemanticAnalysis{Description: "A serene landscape."}, [32]byte{})
	tense := composeMusic(analysis.StructuralAnalysis{}, semantic.SemanticAnalysis{Description: "A threatening scene."}, [32]byte{})
	if countVoiceInSection(calm, "arp", 0, 16) != 0 || countVoiceInSection(calm, "arp", 32, 48) == 0 || countVoiceInSection(calm, "arp", 56, 64) != 0 {
		t.Fatal("the calm counterline did not follow the phrase arc")
	}
	if countVoiceInSection(tense, "arp", 32, 48) <= countVoiceInSection(calm, "arp", 32, 48) {
		t.Fatal("the tense bridge should carry more counterline motion")
	}
}

// TestMelodyDevelopsWithoutRegisterJumps checks that the reprise changes the theme while consecutive pitches stay nearby.
func TestMelodyDevelopsWithoutRegisterJumps(t *testing.T) {
	plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.8}, semantic.SemanticAnalysis{Description: "A serene landscape."}, [32]byte{})
	opening := []int{}
	reprise := []int{}
	previous := 0
	for _, note := range plan.Notes {
		if note.Voice != "lead" {
			continue
		}
		if previous != 0 && (note.MIDI-previous > 12 || previous-note.MIDI > 12) {
			t.Fatalf("melody jumped too far from %d to %d", previous, note.MIDI)
		}
		previous = note.MIDI
		if note.Beat < 16 {
			opening = append(opening, note.MIDI)
		} else if note.Beat < 32 {
			reprise = append(reprise, note.MIDI)
		}
	}
	if len(opening) == 0 || len(opening) != len(reprise) || reflect.DeepEqual(opening, reprise) {
		t.Fatalf("the reprise did not develop the opening motif: %v / %v", opening, reprise)
	}
}

// TestMelodyRhythmFollowsPhraseArc checks that a detailed image grows in the bridge and thins before the final cadence.
func TestMelodyRhythmFollowsPhraseArc(t *testing.T) {
	plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.8}, semantic.SemanticAnalysis{Description: "A serene landscape."}, [32]byte{})
	if countVoiceInSection(plan, "lead", 32, 48) <= countVoiceInSection(plan, "lead", 16, 32) {
		t.Fatal("the bridge should develop the melodic rhythm")
	}
	if countVoiceInSection(plan, "lead", 48, 64) >= countVoiceInSection(plan, "lead", 32, 48) {
		t.Fatal("the conclusion should leave more space than the bridge")
	}
}

// TestExpressNotesPreservesAnchors checks that humanized melody timing stays in its bar while sustained chords remain aligned.
func TestExpressNotesPreservesAnchors(t *testing.T) {
	notes := []MusicNote{
		{Beat: 0, Velocity: 0.3, Voice: "pad"},
		{Beat: 0, Velocity: 0.6, Voice: "lead"},
		{Beat: 3.99, Velocity: 0.4, Voice: "arp"},
		{Beat: 60, Velocity: 0.6, Voice: "lead"},
	}
	var digest [32]byte
	for index := range digest {
		digest[index] = 255
	}
	expressNotes(notes, digest)
	if notes[0].Beat != 0 || notes[0].Velocity != 0.3 || notes[1].Beat <= 0 || notes[2].Beat >= 4 || notes[3].Beat != 60 {
		t.Fatalf("expression moved an anchor or crossed a bar: %+v", notes)
	}
}

// TestMusicServiceComposeRejectsInvalidImage checks that the score cannot be generated from an unbounded or malformed image.
func TestMusicServiceComposeRejectsInvalidImage(t *testing.T) {
	if _, err := (&MusicService{}).Compose("data:image/png;base64,AA==", semantic.SemanticAnalysis{}); err == nil {
		t.Fatal("expected malformed image to be rejected")
	}
}
