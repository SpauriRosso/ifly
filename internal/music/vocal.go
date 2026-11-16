package music

import (
	"strings"
	"unicode"

	"ifly/internal/analysis"
	"ifly/internal/semantic"
)

// shouldHum takes the measured image, tentative scene description, and musical character, then chooses whether a sparse vocal layer fits.
func shouldHum(structure analysis.StructuralAnalysis, interpretation semantic.SemanticAnalysis, character musicalCharacter) bool {
	if character.mood == "calm" {
		return structure.Complexity < 0.58
	}
	if character.mood != "neutral" {
		return false
	}
	words := strings.FieldsFunc(strings.ToLower(interpretation.Description+" "+strings.Join(interpretation.Subjects, " ")), func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	for _, word := range words {
		switch word {
		case "person", "people", "woman", "man", "child", "children", "portrait", "face", "figure", "singer":
			return true
		}
	}
	return false
}

// hummingNotes takes an interlude bar and its chord, then returns a sparse hum on a stable chord tone.
func hummingNotes(start float64, section, bar, chordRoot, chordThird int, energy float64) []MusicNote {
	if section != 3 || bar%4 == 1 || bar%4 == 3 {
		return nil
	}
	interval := 0
	if bar%4 == 2 {
		interval = chordThird
	}
	return []MusicNote{{Beat: start, Duration: 3.2, MIDI: chordRoot + 12 + interval, Velocity: 0.55 * energy, Voice: "hum"}}
}
