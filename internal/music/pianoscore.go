package music

type pianoMeasure struct {
	bar       int
	role      int
	chord     chordSpec
	chordRoot int
	lateChord *chordSpec
	energy    float64
}

// chordAt takes a beat within a measure and returns the active chord, including any mid-bar dominant preparation.
func (measure pianoMeasure) chordAt(beat float64) chordSpec {
	if measure.lateChord != nil && beat >= 2 {
		return *measure.lateChord
	}
	return measure.chord
}

// planPianoMeasures takes the form, tonal root, and scene character, then returns the harmonic changes for every bar.
func planPianoMeasures(root, finalBar int, sections []MusicSection, digest [32]byte, character musicalCharacter) []pianoMeasure {
	measures := make([]pianoMeasure, finalBar+1)
	for bar := range measures {
		role := formEpisodes[bar/4].role
		chord := harmonyForBar(character, bar, role, finalBar)
		chordRoot := root + chord.rootOffset
		measure := pianoMeasure{
			bar: bar, role: role, chord: chord, chordRoot: chordRoot,
			energy: phraseEnergy(sections[bar/4].Energy*phraseDynamics[role][bar%4], bar, digest),
		}
		if bar%4 < 3 && bar < finalBar && pianoAppliedApproaches[bar/4][bar%4] {
			next := harmonyForBar(character, bar+1, role, finalBar)
			approachRoot := (next.rootOffset + 7) % 12
			if approachRoot != chord.rootOffset {
				approach := chordSpec{rootOffset: approachRoot, third: 4}
				measure.lateChord = &approach
			}
		}
		measures[bar] = measure
	}
	return measures
}

// composePianoScore takes the harmonic skeleton and form, then returns one piano part with shaped upper notes and low harmonic arrivals.
func composePianoScore(root, finalBar int, complexity float64, sections []MusicSection, digest [32]byte, character musicalCharacter) []MusicNote {
	measures := planPianoMeasures(root, finalBar, sections, digest, character)
	melody := planPianoScoreLead(root, finalBar, complexity, sections, digest, character, measures)
	notes := make([]MusicNote, 0, 450)
	for _, measure := range measures {
		notes = append(notes, pianoSoloBass(measure, finalBar)...)
		lead := melody[measure.bar]
		notes = append(notes, lead...)
		notes = append(notes, pianoRightHandDoubling(lead, measure.chordRoot, measure.chord.third, measure.bar, measure.role)...)
	}
	return mergePianoUnisons(notes)
}
