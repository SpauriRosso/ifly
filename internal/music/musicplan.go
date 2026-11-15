package music

import (
	"math"
	"sort"
	"strings"
	"unicode"

	"ifly/internal/analysis"
	"ifly/internal/semantic"
)

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
		return musicalStyles["tense"]
	}
	if scores[1] > 0 && scores[1] >= scores[2] {
		return musicalStyles["calm"]
	}
	if scores[2] > 0 {
		return musicalStyles["energetic"]
	}
	return musicalStyles["neutral"]
}

// composeMusic takes measured structure, tentative semantics, and an image digest and returns a four-phrase score.
func composeMusic(structure analysis.StructuralAnalysis, interpretation semantic.SemanticAnalysis, digest [32]byte) MusicPlan {
	character := selectMusicalCharacter(interpretation)
	complexity := math.Max(0, math.Min(1, structure.Complexity))
	tempo := character.baseTempo + int(math.Round(10*complexity))
	root := 48 + int(digest[0]%12)
	motif := [4]int{0, 2, 4, 2}
	for i := range motif {
		motif[i] = (motif[i] + int(digest[i+1]%2)) % len(character.scale)
	}
	sections := []MusicSection{
		{Label: "A", StartBeat: 0, EndBeat: 16, Energy: 0.72},
		{Label: "A'", StartBeat: 16, EndBeat: 32, Energy: 0.88},
		{Label: "B", StartBeat: 32, EndBeat: 48, Energy: 1},
		{Label: "Conclusion", StartBeat: 48, EndBeat: 64, Energy: 0.76},
	}
	notes := make([]MusicNote, 0, 176)
	var previousVoicing [3]int
	for bar := 0; bar < 16; bar++ {
		start := float64(bar * 4)
		section := bar / 4
		chord := character.progression[bar%4]
		if section == 2 || (section == 1 && bar%4 == 3) {
			chord = character.bridge[bar%4]
		}
		if bar == 15 {
			third := 4
			if character.mode == "minor" || character.mode == "dorian" {
				third = 3
			}
			chord = chordSpec{rootOffset: 0, third: third}
		}
		chordRoot := root + chord.rootOffset
		energy := sections[section].Energy
		padLength := 3.75
		if bar == 15 {
			padLength = 2.3
		}
		voicing := voiceChord(chordRoot, chord.third, previousVoicing, bar > 0)
		previousVoicing = voicing
		for _, pitch := range voicing {
			notes = append(notes, MusicNote{start, padLength, pitch, 0.3 * energy, "pad"})
		}
		for _, beat := range character.bassBeats {
			if bar == 15 && beat > 0 {
				continue
			}
			pitch := chordRoot - 12
			if beat > 0 {
				pitch += 7
			}
			notes = append(notes, MusicNote{start + beat, 1.4, pitch, 0.58 * energy, "bass"})
		}
		if bar != 15 {
			for _, beat := range character.pulseBeats {
				notes = append(notes, MusicNote{start + beat, 0.18, root, 0.32 * energy, "pulse"})
			}
		}
		rhythm := character.melodyRhythm[bar%2]
		if complexity < 0.28 && (character.mood == "calm" || character.mood == "neutral") && len(rhythm) > 2 {
			rhythm = rhythm[:2]
		}
		if bar%4 == 3 {
			rhythm = rhythm[:1]
		}
		for index, beat := range rhythm {
			degree := motif[(bar%2*2+index)%len(motif)]
			if section == 1 && index == len(rhythm)-1 {
				degree = (degree + 1) % len(character.scale)
			}
			if section == 2 {
				degree = (len(character.scale) + 4 - degree) % len(character.scale)
			}
			pitch := root + 12 + character.scale[degree]
			if index == 0 {
				pitch = chordRoot + 12
			}
			if bar == 15 {
				pitch = root + 12
			}
			duration := 1.2
			if index+1 < len(rhythm) {
				duration = math.Min(1.5, rhythm[index+1]-beat-0.12)
			} else if bar == 15 {
				duration = 1.6
			} else if bar%4 == 3 {
				duration = 2.2
			}
			velocity := 0.57 * energy
			if index == 0 {
				velocity *= 1.18
			}
			notes = append(notes, MusicNote{start + beat, duration, pitch, velocity, "lead"})
		}
	}
	sort.SliceStable(notes, func(i, j int) bool { return notes[i].Beat < notes[j].Beat })
	return MusicPlan{Mood: character.mood, Mode: character.mode, Timbre: character.timbre, TempoBPM: tempo, BeatsPerLoop: 64, Sections: sections, Notes: notes}
}
