package music

// pianoLeftHandGesture takes a voiced measure and its right-hand phrase, then returns its coordinated lower-hand figure.
func pianoLeftHandGesture(measure pianoMeasure, lead []MusicNote, finalBar int, character musicalCharacter) []MusicNote {
	if measure.bar == finalBar || (measure.role == 5 && measure.bar%4 >= 2) {
		return nil
	}
	texture := pianoTextures[character.mood][measure.role]
	gesture := pianoGestures[texture][measure.bar%4]
	ascendingLead := len(lead) >= 2 && lead[len(lead)-1].MIDI > lead[0].MIDI
	notes := make([]MusicNote, 0, len(gesture.beats))
	for index, beat := range gesture.beats {
		tone := gesture.tones[index]
		if texture == pianoCounterpoint && !ascendingLead {
			tone = 2 - tone
		}
		pitch := measure.voicingAt(beat)[tone] - 12
		for pitch > 60 {
			pitch -= 12
		}
		for pitch < 48 {
			pitch += 12
		}
		duration := 0.48
		if texture == pianoPulse {
			duration = 0.75
		} else if texture == pianoStill {
			duration = 1.1
		}
		velocity := 0.34 * measure.energy
		if index == 0 || (texture == pianoPulse && index == 2) {
			velocity *= 1.16
		}
		notes = append(notes, MusicNote{Beat: float64(measure.bar*4) + beat, Duration: duration, MIDI: pitch, Velocity: velocity, Voice: "arp"})
	}
	return notes
}
