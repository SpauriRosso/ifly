package music

import "testing"

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
