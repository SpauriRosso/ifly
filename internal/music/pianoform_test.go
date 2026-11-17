package music

import "testing"

// TestPianoDevelopmentChangesHarmony checks that the returning bridge and climax are distinct phrases that still prepare a dominant cadence.
func TestPianoDevelopmentChangesHarmony(t *testing.T) {
	for mood, character := range musicalStyles {
		phrases := [3][4]chordSpec{}
		for phrase, firstBar := range []int{8, 24, 36} {
			for position := 0; position < 4; position++ {
				phrases[phrase][position] = harmonyForBar(character, firstBar+position, 2, 63)
			}
			if phrases[phrase][3].rootOffset != 7 {
				t.Fatalf("%s bridge %d misses its dominant ending", mood, phrase)
			}
		}
		if phrases[0] == phrases[1] || phrases[0] == phrases[2] || phrases[1] == phrases[2] {
			t.Fatalf("%s repeats a bridge instead of developing it: %v", mood, phrases)
		}
	}
}
