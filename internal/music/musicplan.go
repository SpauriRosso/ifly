package music

import (
	"math"
	"strings"
	"unicode"

	"ifly/internal/analysis"
	"ifly/internal/semantic"
)

type musicalCharacter struct {
	mood        string
	mode        string
	baseTempo   int
	scale       []int
	progression []int
	triad       []int
}

// selectMusicalCharacter takes the model's unverified words and returns a conservative mood and tonal vocabulary.
func selectMusicalCharacter(interpretation semantic.SemanticAnalysis) musicalCharacter {
	words := strings.FieldsFunc(strings.ToLower(interpretation.Description+" "+strings.Join(interpretation.Subjects, " ")), func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	scores := [3]int{}
	for _, word := range words {
		switch word {
		case "violent", "violence", "threat", "threatening", "torture", "suffering", "blood", "fear", "terrifying", "ominous", "tense", "danger", "storm":
			scores[0] += 2
		case "serene", "peaceful", "calm", "tranquil", "gentle", "quiet", "relaxed":
			scores[1] += 2
		case "sunset", "lullaby":
			scores[1]++
		case "dynamic", "energetic", "playful", "lively", "dancing", "vibrant", "motion":
			scores[2] += 2
		}
	}
	if interpretation.Uncertain {
		for index, score := range scores {
			if score < 2 {
				scores[index] = 0
			}
		}
	}
	if scores[0] > 0 && scores[0] >= scores[1] && scores[0] >= scores[2] {
		return musicalCharacter{"tense", "minor", 102, []int{0, 2, 3, 5, 7, 8, 10}, []int{0, 8, 10, 7}, []int{0, 3, 7}}
	}
	if scores[1] > 0 && scores[1] >= scores[2] {
		return musicalCharacter{"calm", "pentatonic", 65, []int{0, 2, 4, 7, 9}, []int{0, 5, 7, 0}, []int{0, 4, 7}}
	}
	if scores[2] > 0 {
		return musicalCharacter{"energetic", "major", 106, []int{0, 2, 4, 5, 7, 9, 11}, []int{0, 5, 7, 9}, []int{0, 4, 7}}
	}
	return musicalCharacter{"neutral", "dorian", 78, []int{0, 2, 3, 5, 7, 9, 10}, []int{0, 5, 10, 7}, []int{0, 3, 7}}
}

// composeMusic takes measured structure, tentative semantics, and an image digest and returns an eight-bar score.
func composeMusic(structure analysis.StructuralAnalysis, interpretation semantic.SemanticAnalysis, digest [32]byte) MusicPlan {
	character := selectMusicalCharacter(interpretation)
	complexity := math.Max(0, math.Min(1, structure.Complexity))
	tempo := character.baseTempo + int(math.Round(14*complexity))
	root := 48 + int(digest[0]%12)
	density := 2 + int(math.Round(2*complexity))
	motif := [4]int{0, 2, 4, 2}
	for i := range motif {
		motif[i] = (motif[i] + int(digest[i+1]%3)) % len(character.scale)
	}
	notes := make([]MusicNote, 0, 8*(3+4))
	for bar := 0; bar < 8; bar++ {
		start := float64(bar * 4)
		chordRoot := root + character.progression[bar%4]
		for _, interval := range character.triad {
			notes = append(notes, MusicNote{start, 3.7, chordRoot + interval, 0.32, "pad"})
		}
		for step := 0; step < density; step++ {
			degree := motif[step]
			if bar >= 4 && step == density-1 {
				degree = (degree + 1) % len(character.scale)
			}
			beat := start + float64(step)*4/float64(density)
			velocity := 0.52
			if step == 0 {
				velocity = 0.62
			}
			notes = append(notes, MusicNote{beat, 1.2, root + 12 + character.scale[degree], velocity, "lead"})
		}
	}
	return MusicPlan{Mood: character.mood, Mode: character.mode, TempoBPM: tempo, BeatsPerLoop: 32, Notes: notes}
}
