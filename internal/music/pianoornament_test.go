package music

import "testing"

// TestPianoOrnamentPassesBetweenStructuralNotes checks that a developed gesture fills a third with a stepwise scale tone.
func TestPianoOrnamentPassesBetweenStructuralNotes(t *testing.T) {
	measure := pianoMeasure{bar: 8, role: 2, chord: chordSpec{rootOffset: 0, third: 4}}
	lead := []MusicNote{
		{Beat: 32, Duration: 1.1, MIDI: 72, Velocity: 0.7, Voice: "lead"},
		{Beat: 33, Duration: 1.1, MIDI: 76, Velocity: 0.7, Voice: "lead"},
	}
	got := ornamentPianoBar(lead, measure, 0.5, 48, []int{0, 2, 4, 5, 7, 9, 11})
	if len(got) != 3 || got[1].MIDI != 74 || got[1].Beat != 32.5 || got[0].Duration >= 0.5 {
		t.Fatalf("the structural third was not elaborated into playable steps: %+v", got)
	}
}
