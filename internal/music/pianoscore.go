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
		notes = append(notes, pianoInnerVoices(measure, finalBar, character)...)
		if measure.bar != finalBar {
			notes = append(notes, pianoCadentialSeventh(start, measure.bar, measure.role, measure.chordRoot, measure.chord, measure.energy)...)
		}
		notes = append(notes, pianoBassLine(measure, finalBar, character)...)
		lead := melody[measure.bar]
		notes = append(notes, lead...)
		notes = append(notes, pianoRightHandDoubling(lead, measure.chordRoot, measure.chord.third, measure.bar, measure.role)...)
		notes = append(notes, pianoLeftHandGesture(measure, lead, finalBar, character)...)
	}
	return mergePianoUnisons(notes)
}
