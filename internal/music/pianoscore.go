package music

type pianoMeasure struct {
	bar       int
	role      int
	chord     chordSpec
	chordRoot int
	voicing   [3]int
	energy    float64
}

// planPianoMeasures takes the form, tonal root, and scene character, then returns a voice-led harmonic skeleton for every bar.
func planPianoMeasures(root, finalBar int, sections []MusicSection, digest [32]byte, character musicalCharacter) []pianoMeasure {
	measures := make([]pianoMeasure, finalBar+1)
	var previous [3]int
	for bar := range measures {
		role := formEpisodes[bar/4].role
		chord := harmonyForBar(character, bar, role, finalBar)
		chordRoot := root + chord.rootOffset
		voicing := voiceChord(chordRoot, chord.third, previous, bar > 0)
		measures[bar] = pianoMeasure{
			bar: bar, role: role, chord: chord, chordRoot: chordRoot,
			voicing: voicing,
			energy:  phraseEnergy(sections[bar/4].Energy*phraseDynamics[role][bar%4], bar, digest),
		}
		previous = voicing
	}
	return measures
}

// composePianoScore takes the harmonic skeleton and form, then writes both piano hands and their inner chord voices as one score.
func composePianoScore(root, finalBar int, complexity float64, sections []MusicSection, digest [32]byte, character musicalCharacter) []MusicNote {
	measures := planPianoMeasures(root, finalBar, sections, digest, character)
	melody := planPianoMelody(root, finalBar, complexity, sections, digest, character, measures)
	notes := make([]MusicNote, 0, 700)
	for _, measure := range measures {
		start := float64(measure.bar * 4)
		padLength := padSustainBeats[measure.bar%4]
		if measure.role == 3 {
			padLength = 3.4
		}
		if measure.bar == finalBar {
			padLength = 2
		}
		for voiceIndex, pitch := range measure.voicing {
			if (measure.role == 0 || measure.role == 3) && voiceIndex == 1 {
				continue
			}
			notes = append(notes, MusicNote{start, padLength, pitch, 0.24 * measure.energy, "pad"})
		}
		if character.colorTone != 0 && measure.bar%2 == 0 && measure.role != 3 && measure.bar != finalBar {
			notes = append(notes, MusicNote{start + 0.3, 2.1, measure.chordRoot + character.colorTone, 0.16 * measure.energy, "pad"})
		}
		if measure.bar != finalBar {
			notes = append(notes, pianoCadentialSeventh(start, measure.bar, measure.role, measure.chordRoot, measure.chord, measure.energy)...)
		}
		notes = append(notes, bassNotes(start, measure.bar, measure.role, finalBar, measure.chordRoot, measure.chord.third, measure.energy, character)...)
		lead := melody[measure.bar]
		notes = append(notes, lead...)
		notes = append(notes, pianoRightHandDoubling(lead, measure.chordRoot, measure.chord.third, measure.bar, measure.role)...)
		notes = append(notes, pianoLeftHandNotes(start, measure.bar, measure.role, finalBar, measure.voicing, measure.energy, character)...)
	}
	return notes
}
