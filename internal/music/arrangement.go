package music

import "math"

// avoidsLeadOnsets takes a candidate beat and lead events, then reports whether the counterline leaves the melody's attacks clear.
func avoidsLeadOnsets(beat float64, lead []MusicNote) bool {
	for _, note := range lead {
		if math.Abs(note.Beat-beat) < 0.36 {
			return false
		}
	}
	return true
}

// arpeggioNotes takes a bar's harmony, phrase position, and lead notes, then places light counterline answers away from lead attacks.
func arpeggioNotes(start float64, section, bar int, voicing [3]int, energy float64, character musicalCharacter, lead []MusicNote) []MusicNote {
	if section == 0 || section == 3 || (section == 1 && bar%4 < 2) || (section == 5 && bar%4 >= 2) {
		return nil
	}
	rhythm := character.arpBeats
	if section != 2 && len(rhythm) > 2 {
		rhythm = rhythm[:2]
	}
	notes := make([]MusicNote, 0, len(rhythm))
	for index, beat := range rhythm {
		onset := start + beat
		if !avoidsLeadOnsets(onset, lead) {
			onset += 0.5
			if onset > start+3.5 || !avoidsLeadOnsets(onset, lead) {
				continue
			}
		}
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
		duration = math.Min(duration, start+3.9-onset)
		notes = append(notes, MusicNote{Beat: onset, Duration: duration, MIDI: pitch, Velocity: velocity, Voice: "arp"})
	}
	return notes
}
