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

// composeMusic takes measured structure, tentative semantics, and an image digest and returns a six-phrase score.
func composeMusic(structure analysis.StructuralAnalysis, interpretation semantic.SemanticAnalysis, digest [32]byte) MusicPlan {
	character := selectMusicalCharacter(interpretation)
	useHum := shouldHum(structure, interpretation, character)
	complexity := math.Max(0, math.Min(1, structure.Complexity))
	tempo := character.baseTempo + int(math.Round(10*complexity))
	root := 48 + int(digest[0]%12)
	motif := makeMotif(digest, len(character.scale))
	bridgeMotif := makeBridgeMotif(digest, len(character.scale))
	sections := []MusicSection{
		{Label: "A", StartBeat: 0, EndBeat: 16, Energy: 0.72, TempoBPM: tempo + character.tempoArc[0]},
		{Label: "A'", StartBeat: 16, EndBeat: 32, Energy: 0.88, TempoBPM: tempo + character.tempoArc[1]},
		{Label: "B", StartBeat: 32, EndBeat: 48, Energy: 1, TempoBPM: tempo + character.tempoArc[2]},
		{Label: "Interlude", StartBeat: 48, EndBeat: 64, Energy: 0.42, TempoBPM: tempo + character.tempoArc[3]},
		{Label: "Retour", StartBeat: 64, EndBeat: 80, Energy: 0.92, TempoBPM: tempo + character.tempoArc[4]},
		{Label: "Conclusion", StartBeat: 80, EndBeat: 96, Energy: 0.68, TempoBPM: tempo + character.tempoArc[5]},
	}
	notes := make([]MusicNote, 0, 256)
	var previousVoicing [3]int
	previousLeadPitch := 0
	finalBar := len(sections)*4 - 1
	for bar := 0; bar <= finalBar; bar++ {
		start := float64(bar * 4)
		section := bar / 4
		chord := character.progression[bar%4]
		if section == 2 || (section == 1 && bar%4 == 3) || (section == 4 && bar%4 == 2) {
			chord = character.bridge[bar%4]
		}
		if section == 3 {
			chord = character.bridge[(bar+2)%4]
		}
		if bar == finalBar {
			third := 4
			if character.mode == "minor" || character.mode == "dorian" {
				third = 3
			}
			chord = chordSpec{rootOffset: 0, third: third}
		}
		chordRoot := root + chord.rootOffset
		energy := sections[section].Energy
		padLength := 3.75
		if bar == finalBar {
			padLength = 2.3
		}
		voicing := voiceChord(chordRoot, chord.third, previousVoicing, bar > 0)
		previousVoicing = voicing
		for voiceIndex, pitch := range voicing {
			if (section == 0 || section == 3) && voiceIndex == 1 {
				continue
			}
			notes = append(notes, MusicNote{start, padLength, pitch, 0.24 * energy, "pad"})
		}
		notes = append(notes, arpeggioNotes(start, section, bar, voicing, energy, character)...)
		notes = append(notes, bassNotes(start, bar, section, finalBar, chordRoot, chord.third, energy, character)...)
		if section != 3 && bar != finalBar && (section != 5 || bar%4 < 2) {
			for _, beat := range character.pulseBeats {
				notes = append(notes, MusicNote{start + beat, 0.18, root, 0.32 * energy, "pulse"})
			}
		}
		currentMotif := motif
		if section == 2 || section == 3 {
			currentMotif = bridgeMotif
		}
		lead, finalPitch := melodyNotes(start, bar, section, finalBar, root, chord, currentMotif, previousLeadPitch, energy, complexity, character)
		notes = append(notes, lead...)
		previousLeadPitch = finalPitch
		if useHum {
			notes = append(notes, hummingNotes(start, section, bar, root, character.scale, energy)...)
		}
	}
	expressNotes(notes, digest, float64(finalBar*4))
	sort.SliceStable(notes, func(i, j int) bool { return notes[i].Beat < notes[j].Beat })
	return MusicPlan{Mood: character.mood, Mode: character.mode, Timbre: character.timbre, TempoBPM: tempo, BeatsPerLoop: finalBar*4 + 4, Sections: sections, Notes: notes}
}
