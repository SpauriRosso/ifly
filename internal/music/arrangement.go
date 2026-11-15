package music

// arpeggioNotes takes a bar's harmony and phrase position and returns a light counterline only where the arrangement grows.
func arpeggioNotes(start float64, section, bar int, voicing [3]int, energy float64, character musicalCharacter) []MusicNote {
	if section == 0 || section == 3 || (section == 1 && bar%4 < 2) || (section == 5 && bar%4 >= 2) {
		return nil
	}
	rhythm := character.arpBeats
	if section != 2 && len(rhythm) > 2 {
		rhythm = rhythm[:2]
	}
	notes := make([]MusicNote, 0, len(rhythm))
	for index, beat := range rhythm {
		pitch := voicing[(index+bar)%3] + 12
		velocity := 0.34 * energy
		if section == 2 {
			velocity = 0.4 * energy
		}
		notes = append(notes, MusicNote{Beat: start + beat, Duration: 0.44, MIDI: pitch, Velocity: velocity, Voice: "arp"})
	}
	return notes
}
