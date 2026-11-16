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
		voiceIndex := (bar + index + section) % len(voicing)
		if section == 2 && bar%2 == 1 {
			voiceIndex = (bar + len(voicing) - index%len(voicing)) % len(voicing)
		}
		pitch := voicing[voiceIndex]
		velocity := 0.34 * energy
		if section == 2 {
			velocity = 0.4 * energy
		}
		duration := 0.52
		if section == 2 && index%2 == 1 {
			duration = 0.38
		}
		notes = append(notes, MusicNote{Beat: start + beat, Duration: duration, MIDI: pitch, Velocity: velocity, Voice: "arp"})
	}
	return notes
}
