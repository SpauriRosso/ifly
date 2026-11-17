package music

type pianoMeasure struct {
	bar         int
	role        int
	chord       chordSpec
	chordRoot   int
	voicing     [3]int
	lateChord   *chordSpec
	lateVoicing [3]int
	energy      float64
}

// chordAt takes a beat within a measure and returns the active chord, including any mid-bar dominant preparation.
func (measure pianoMeasure) chordAt(beat float64) chordSpec {
	if measure.lateChord != nil && beat >= 2 {
		return *measure.lateChord
	}
	return measure.chord
}

// voicingAt takes a beat within a measure and returns the active left-hand inversion.
func (measure pianoMeasure) voicingAt(beat float64) [3]int {
	if measure.lateChord != nil && beat >= 2 {
		return measure.lateVoicing
	}
	return measure.voicing
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
		measure := pianoMeasure{
			bar: bar, role: role, chord: chord, chordRoot: chordRoot,
			voicing: voicing,
			energy:  phraseEnergy(sections[bar/4].Energy*phraseDynamics[role][bar%4], bar, digest),
		}
		if bar%4 == 2 && bar < finalBar && role != 3 && role != 5 {
			next := harmonyForBar(character, bar+1, role, finalBar)
			if next.rootOffset == 7 && chord.rootOffset != 7 {
				approach := chordSpec{rootOffset: 2, third: 4}
				measure.lateChord = &approach
				measure.lateVoicing = voiceChord(root+approach.rootOffset, approach.third, voicing, true)
			}
		}
		measures[bar] = measure
		previous = measure.voicingAt(3)
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
		if measure.lateChord != nil {
			padLength = 1.7
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
		bass := bassNotes(start, measure.bar, measure.role, finalBar, measure.chordRoot, measure.chord.third, measure.energy, character)
		if measure.lateChord != nil {
			bass = []MusicNote{bass[0], {Beat: start + 2, Duration: 1.55, MIDI: root + measure.lateChord.rootOffset - 12, Velocity: 0.48 * measure.energy, Voice: "bass"}}
			notes = append(notes, MusicNote{Beat: start + 2, Duration: 1.5, MIDI: measure.lateVoicing[0], Velocity: 0.27 * measure.energy, Voice: "pad"})
		}
		notes = append(notes, bass...)
		lead := melody[measure.bar]
		notes = append(notes, lead...)
		notes = append(notes, pianoRightHandDoubling(lead, measure.chordRoot, measure.chord.third, measure.bar, measure.role)...)
		notes = append(notes, pianoLeftHandGesture(measure, lead, finalBar, character)...)
	}
	return notes
}
