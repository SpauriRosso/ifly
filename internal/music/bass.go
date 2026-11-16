package music

// bassNotes takes one bar's chord and form position, then returns a low line with section-dependent responses.
func bassNotes(start float64, bar, section, finalBar, chordRoot, chordThird int, energy float64, character musicalCharacter) []MusicNote {
	root := chordRoot - 12
	baseDuration := 1.6
	if section == 3 || bar == finalBar {
		baseDuration = 2.8
	}
	notes := []MusicNote{{Beat: start, Duration: baseDuration, MIDI: root, Velocity: 0.58 * energy, Voice: "bass"}}
	if section == 3 || bar == finalBar {
		return notes
	}
	responseBeat := 2.5
	if len(character.bassBeats) > 1 {
		responseBeat = character.bassBeats[1]
	} else {
		if (section != 2 && section != 4) || bar%4 == 0 || bar%4 == 3 {
			return notes
		}
	}
	interval := 7
	switch bar % 4 {
	case 1:
		interval = 12
	case 2:
		interval = chordThird
	}
	notes = append(notes, MusicNote{Beat: start + responseBeat, Duration: 1.15, MIDI: root + interval, Velocity: 0.38 * energy, Voice: "bass"})
	return notes
}
