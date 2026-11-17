package music

// pianoBassRegister takes a chord tone and returns its low piano octave between MIDI 36 and 50.
func pianoBassRegister(pitch int) int {
	for pitch > 50 {
		pitch -= 12
	}
	for pitch < 36 {
		pitch += 12
	}
	return pitch
}

// pianoSoloBass takes one harmonic measure and the final bar, then returns only the low piano attacks needed at its chord changes.
func pianoSoloBass(measure pianoMeasure, finalBar int) []MusicNote {
	start := float64(measure.bar * 4)
	root := pianoBassRegister(measure.chordRoot - 12)
	duration := 3.8
	if measure.lateChord != nil {
		duration = 1.9
	}
	notes := []MusicNote{{Beat: start, Duration: duration, MIDI: root, Velocity: 0.47 * measure.energy, Voice: "bass"}}
	if measure.lateChord != nil {
		approach := pianoBassRegister(measure.chordRoot - measure.chord.rootOffset + measure.lateChord.rootOffset - 12)
		notes = append(notes, MusicNote{Beat: start + 2, Duration: 1.85, MIDI: approach, Velocity: 0.43 * measure.energy, Voice: "bass"})
	}
	if measure.bar == finalBar {
		notes[0].Duration = 3.5
	}
	return notes
}
