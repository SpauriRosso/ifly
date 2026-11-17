package music

// pianoInnerVoices takes a voiced measure and final-bar marker, then returns short inner chord attacks that support each phrase.
func pianoInnerVoices(measure pianoMeasure, finalBar int, character musicalCharacter) []MusicNote {
	position := measure.bar % 4
	count := pianoChordVoices[measure.role][position]
	if count == 0 {
		return nil
	}
	start := float64(measure.bar * 4)
	duration := pianoChordDurations[measure.role][position]
	if measure.bar == finalBar {
		duration = 2.2
	}
	notes := make([]MusicNote, 0, count+3)
	for voice := 0; voice < count; voice++ {
		pitch := measure.voicing[voice]
		if count == 1 {
			pitch = measure.voicing[1]
		}
		notes = append(notes, MusicNote{Beat: start, Duration: duration, MIDI: pitch, Velocity: 0.3 * measure.energy, Voice: "pad"})
	}
	if measure.lateChord != nil {
		for voice := 0; voice < 2; voice++ {
			notes = append(notes, MusicNote{Beat: start + 2, Duration: 1.1, MIDI: measure.lateVoicing[voice], Velocity: 0.28 * measure.energy, Voice: "pad"})
		}
	} else if (measure.role == 2 || measure.role == 4) && position == 1 {
		notes = append(notes, MusicNote{Beat: start + 2.5, Duration: 0.7, MIDI: measure.voicing[1], Velocity: 0.22 * measure.energy, Voice: "pad"})
	}
	if character.colorTone != 0 && position == 0 && (measure.role == 0 || measure.role == 4) {
		notes = append(notes, MusicNote{Beat: start + 0.25, Duration: 0.9, MIDI: measure.chordRoot + character.colorTone, Velocity: 0.17 * measure.energy, Voice: "pad"})
	}
	return notes
}
