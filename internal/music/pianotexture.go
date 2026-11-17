package music

// pianoRightHandDoubling takes a melodic arrival and its chord, then returns an occasional consonant lower note in the right hand.
func pianoRightHandDoubling(lead []MusicNote, chordRoot, third, bar, section int) []MusicNote {
	if len(lead) == 0 || section == 3 || section == 5 || (bar%4 != 0 && (section != 2 || bar%4 != 2)) {
		return nil
	}
	arrival := lead[0]
	chord := chordPitchClasses(chordRoot, third)
	for distance := 3; distance <= 9; distance++ {
		pitch := arrival.MIDI - distance
		if pitch >= 67 && pitchBelongsToChord(pitch, chord) {
			return []MusicNote{{Beat: arrival.Beat, Duration: arrival.Duration * 0.85, MIDI: pitch, Velocity: arrival.Velocity * 0.52, Voice: "lead"}}
		}
	}
	return nil
}
