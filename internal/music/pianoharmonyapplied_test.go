package music

import "testing"

// TestPianoAppliedDominantsResolveByFifth checks that written chromatic approach chords lead into their intended next harmony.
func TestPianoAppliedDominantsResolveByFifth(t *testing.T) {
	for mood, character := range musicalStyles {
		sections := make([]MusicSection, len(formEpisodes))
		measures := planPianoMeasures(48, len(formEpisodes)*4-1, sections, [32]byte{}, character)
		approaches := 0
		for bar, measure := range measures {
			if measure.lateChord == nil {
				continue
			}
			approaches++
			if bar+1 >= len(measures) || (measures[bar+1].chord.rootOffset-measure.lateChord.rootOffset+12)%12 != 5 || measure.lateChord.third != 4 {
				t.Fatalf("%s bar %d has an applied chord without a fifth-resolution: %+v then %+v", mood, bar, measure.lateChord, measures[bar+1].chord)
			}
		}
		if approaches < 3 {
			t.Fatalf("%s never moves through enough applied harmonies: %d", mood, approaches)
		}
	}
}

// TestPianoCalmPeriodTonicizesRelativeMinor checks that the opening period prepares its relative minor with a chromatic dominant.
func TestPianoCalmPeriodTonicizesRelativeMinor(t *testing.T) {
	character := musicalStyles["calm"]
	sections := make([]MusicSection, len(formEpisodes))
	measures := planPianoMeasures(48, len(formEpisodes)*4-1, sections, [32]byte{}, character)
	if measures[0].lateChord == nil || measures[0].lateChord.rootOffset != 4 || measures[1].chord.rootOffset != 9 {
		t.Fatalf("the opening lacks V/vi to vi: %+v then %+v", measures[0].lateChord, measures[1].chord)
	}
}
