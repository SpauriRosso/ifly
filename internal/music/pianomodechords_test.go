package music

import (
	"math"
	"testing"

	"ifly/internal/analysis"
	"ifly/internal/semantic"
)

// TestPianoModeHarmonyChecksEverySeedPalette checks that each written triad belongs to the mode selected for its scene.
func TestPianoModeHarmonyChecksEverySeedPalette(t *testing.T) {
	for mood, base := range musicalStyles {
		for variant := byte(0); variant < 3; variant++ {
			var digest [32]byte
			digest[3] = variant
			character := choosePianoPalette(base, digest)
			for bar := 0; bar < len(formEpisodes)*4; bar++ {
				chord := harmonyForBar(character, bar, formEpisodes[bar/4].role, len(formEpisodes)*4-1)
				for _, interval := range []int{0, chord.third, 7} {
					pitchClass := (chord.rootOffset + interval) % 12
					found := false
					for _, degree := range character.scale {
						if degree == pitchClass {
							found = true
						}
					}
					if !found {
						t.Fatalf("%s variant %d bar %d has chord tone %d outside %s", mood, variant, bar, pitchClass, character.mode)
					}
				}
			}
		}
	}
}

// TestPianoStrongBeatsAgreeWithHarmony checks that both hands land on the active chord at structural beats.
func TestPianoStrongBeatsAgreeWithHarmony(t *testing.T) {
	for _, description := range []string{"A serene pond.", "A lonely figure mourning.", "A threatening storm.", "A joyful dance."} {
		meaning := semantic.SemanticAnalysis{Description: description}
		for variant := byte(0); variant < 3; variant++ {
			var digest [32]byte
			digest[3] = variant
			root := 48 + int(digest[0]%12)
			plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.5}, meaning, digest)
			character := choosePianoPalette(selectMusicalCharacter(meaning), digest)
			measures := planPianoMeasures(root, plan.BeatsPerLoop/4-1, plan.Sections, digest, character)
			for _, note := range plan.Notes {
				if note.Voice != "lead" || math.Abs(math.Remainder(note.Beat, 2)) > 0.05 {
					continue
				}
				bar := int(note.Beat / 4)
				chord := measures[bar].chordAt(note.Beat - float64(bar*4))
				if !pitchBelongsToChord(note.MIDI, chordPitchClasses(root+chord.rootOffset, chord.third)) {
					t.Fatalf("%s variant %d beat %.2f lands on MIDI %d outside active chord %+v", description, variant, note.Beat, note.MIDI, chord)
				}
			}
		}
	}
}

// TestPianoAttacksFitModeOrChord checks that each right-hand attack is explained by the scene's mode or its active harmony.
func TestPianoAttacksFitModeOrChord(t *testing.T) {
	for _, description := range []string{"A serene pond.", "A lonely figure mourning.", "A threatening storm.", "A joyful dance."} {
		meaning := semantic.SemanticAnalysis{Description: description}
		for variant := byte(0); variant < 3; variant++ {
			var digest [32]byte
			digest[3] = variant
			root := 48 + int(digest[0]%12)
			plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.5}, meaning, digest)
			character := choosePianoPalette(selectMusicalCharacter(meaning), digest)
			measures := planPianoMeasures(root, plan.BeatsPerLoop/4-1, plan.Sections, digest, character)
			for _, note := range plan.Notes {
				if note.Voice != "lead" {
					continue
				}
				bar := int(note.Beat / 4)
				chord := measures[bar].chordAt(note.Beat - float64(bar*4))
				inChord := pitchBelongsToChord(note.MIDI, chordPitchClasses(root+chord.rootOffset, chord.third))
				inMode := false
				for _, degree := range character.scale {
					if note.MIDI%12 == (root+degree)%12 {
						inMode = true
					}
				}
				if !inChord && !inMode {
					t.Fatalf("%s variant %d beat %.2f has MIDI %d outside mode %s and chord %+v", description, variant, note.Beat, note.MIDI, character.mode, chord)
				}
			}
		}
	}
}
