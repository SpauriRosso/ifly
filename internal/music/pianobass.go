package music

// pianoBassRegister takes a chord tone and returns its nearest lower-piano octave between MIDI 36 and 50.
func pianoBassRegister(pitch int) int {
	for pitch > 50 {
		pitch -= 12
	}
	for pitch < 36 {
		pitch += 12
	}
	return pitch
}

// pianoBassLine takes a harmonic measure and returns a low root with an optional chord-fifth or secondary-dominant response.
func pianoBassLine(measure pianoMeasure, finalBar int, character musicalCharacter) []MusicNote {
	start := float64(measure.bar * 4)
	root := pianoBassRegister(measure.chordRoot - 12)
	duration := 1.55
	if measure.role == 3 || measure.bar == finalBar {
		duration = 2.8
	}
	notes := []MusicNote{{Beat: start, Duration: duration, MIDI: root, Velocity: 0.52 * measure.energy, Voice: "bass"}}
	if measure.lateChord != nil {
		approach := pianoBassRegister(measure.chordRoot - measure.chord.rootOffset + measure.lateChord.rootOffset - 12)
		return append(notes, MusicNote{Beat: start + 2, Duration: 1.55, MIDI: approach, Velocity: 0.44 * measure.energy, Voice: "bass"})
	}
	if measure.bar == finalBar || measure.bar%4 == 3 {
		return notes
	}
	beat := pianoBassResponses[character.mood][measure.role]
	if beat == 0 {
		return notes
	}
	fifth := pianoBassRegister(measure.chordRoot - 12 + 7)
	return append(notes, MusicNote{Beat: start + beat, Duration: 0.95, MIDI: fifth, Velocity: 0.34 * measure.energy, Voice: "bass"})
}
