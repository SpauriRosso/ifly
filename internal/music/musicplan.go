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
	words := strings.FieldsFunc(strings.ToLower(interpretation.Description+" "+interpretation.Atmosphere+" "+strings.Join(interpretation.Subjects, " ")), func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	scores := [3]int{}
	for _, word := range words {
		switch word {
		case "violent", "violence", "threat", "threatening", "torture", "suffering", "blood", "fear", "terrifying", "ominous", "tense", "danger", "storm":
			scores[0] += 2
		case "serene", "peaceful", "calm", "tranquil", "gentle", "quiet", "relaxed", "dreamy", "reflective", "soft", "delicate":
			scores[1] += 2
		case "sunset", "lullaby", "lily", "lilies", "pond", "garden", "blossoms", "flowers", "reflection":
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

// composeMusic takes measured structure, tentative semantics, and an image digest and returns a long-form score with related episodes.
func composeMusic(structure analysis.StructuralAnalysis, interpretation semantic.SemanticAnalysis, digest [32]byte) MusicPlan {
	character := selectMusicalCharacter(interpretation)
	useHum := shouldHum(structure, interpretation, character)
	complexity := math.Max(0, math.Min(1, structure.Complexity))
	tempo := character.baseTempo + int(math.Round(10*complexity))
	root := 48 + int(digest[0]%12)
	motif := makeMotif(digest, len(character.scale))
	bridgeMotif := makeBridgeMotif(digest, len(character.scale))
	sections := make([]MusicSection, 0, len(formEpisodes))
	for index, episode := range formEpisodes {
		sections = append(sections, MusicSection{
			Label: episode.label, StartBeat: float64(index * 16), EndBeat: float64((index + 1) * 16),
			Energy: episode.energy, TempoBPM: tempo + character.tempoArc[episode.role] + episode.tempoShift,
		})
	}
	notes := make([]MusicNote, 0, 700)
	var previousVoicing [3]int
	previousLeadPitch := 0
	finalBar := len(sections)*4 - 1
	for bar := 0; bar <= finalBar; bar++ {
		start := float64(bar * 4)
		episodeIndex := bar / 4
		episode := formEpisodes[episodeIndex]
		section := episode.role
		chord := harmonyForBar(character, bar, section, finalBar)
		chordRoot := root + chord.rootOffset
		energy := phraseEnergy(sections[episodeIndex].Energy*phraseDynamics[section][bar%4], bar, digest)
		padLength := padSustainBeats[bar%4]
		if section == 3 {
			padLength = 3.4
		}
		if bar == finalBar {
			padLength = 2.0
		}
		voicing := voiceChord(chordRoot, chord.third, previousVoicing, bar > 0)
		previousVoicing = voicing
		for voiceIndex, pitch := range voicing {
			if (section == 0 || section == 3) && voiceIndex == 1 {
				continue
			}
			notes = append(notes, MusicNote{start, padLength, pitch, 0.24 * energy, "pad"})
		}
		if character.colorTone != 0 && bar%2 == 0 && section != 3 && bar != finalBar {
			notes = append(notes, MusicNote{start + 0.3, 2.1, chordRoot + character.colorTone, 0.16 * energy, "pad"})
		}
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
		currentMotif = turnMotif(currentMotif, episode.motifTurn)
		lead, finalPitch := melodyNotes(start, bar, section, finalBar, root, chord, currentMotif, previousLeadPitch, energy, complexity, character)
		notes = append(notes, lead...)
		notes = append(notes, arpeggioNotes(start, section, bar, voicing, energy, character, lead)...)
		previousLeadPitch = finalPitch
		if useHum {
			notes = append(notes, hummingNotes(start, section, bar, chordRoot, chord.third, energy)...)
		}
	}
	expressNotes(notes, digest, float64(finalBar*4))
	sort.SliceStable(notes, func(i, j int) bool { return notes[i].Beat < notes[j].Beat })
	return MusicPlan{Mood: character.mood, Mode: character.mode, Timbre: character.timbre, TempoBPM: tempo, BeatsPerLoop: finalBar*4 + 4, Sections: sections, Notes: notes}
}
