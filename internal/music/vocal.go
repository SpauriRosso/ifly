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

// hummingNotes takes the interlude bar and tonal center, then returns two long tonic and third responses across that phrase.
func hummingNotes(start float64, section, bar, root int, scale []int, energy float64) []MusicNote {
	if section != 3 || bar%4 == 1 || bar%4 == 3 {
		return nil
	}
	degree := 0
	if bar%4 == 2 {
		degree = 2
	}
	return []MusicNote{{Beat: start, Duration: 3.2, MIDI: root + 12 + scale[degree], Velocity: 0.55 * energy, Voice: "hum"}}
}
