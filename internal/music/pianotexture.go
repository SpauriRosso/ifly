package music

// pianoLeftHandNotes takes a bar's chord voicing and formal role, then returns a lower-register broken-chord figure for the piano.
func pianoLeftHandNotes(start float64, bar, section, finalBar int, voicing [3]int, energy float64, character musicalCharacter) []MusicNote {
	if bar == finalBar || (section == 5 && bar%4 >= 2) {
		return nil
	}
	beats := pianoLeftHandBeats
	if section == 2 || (section == 4 && bar%4 < 3) {
		beats = pianoDevelopedBeats
	}
	if section == 3 || (section == 5 && bar%4 == 1) {
		beats = pianoReflectiveBeats
	}
	if bar%4 == 3 && section != 2 && section != 3 {
		beats = pianoReflectiveBeats
	}
	notes := make([]MusicNote, 0, len(beats))
	for index, beat := range beats {
		voiceIndex := character.pianoFigure[index%len(character.pianoFigure)]
		if section == 2 && bar%2 == 1 {
			voiceIndex = character.pianoFigure[(len(character.pianoFigure)-1-index)%len(character.pianoFigure)]
		}
		duration := 0.4
		if section == 3 || (section == 5 && bar%4 == 1) {
			duration = 0.8
		}
		velocity := 0.27 * energy
		if index%4 == 0 {
			velocity *= 1.13
		}
		pitch := voicing[voiceIndex] - 12
		for pitch > 60 {
			pitch -= 12
		}
		for pitch < 48 {
			pitch += 12
		}
		notes = append(notes, MusicNote{Beat: start + beat, Duration: duration, MIDI: pitch, Velocity: velocity, Voice: "arp"})
	}
	return notes
}
