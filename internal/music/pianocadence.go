package music

// pianoCadentialSeventh takes a phrase-ending dominant chord and returns its seventh as a quiet inner piano voice before the tonic.
func pianoCadentialSeventh(start float64, bar, section int, chordRoot int, chord chordSpec, energy float64) []MusicNote {
	if bar%4 != 3 || section == 3 || chord.rootOffset != 7 || chord.third != 4 {
		return nil
	}
	return []MusicNote{{Beat: start + 2, Duration: 1.6, MIDI: chordRoot + 10, Velocity: 0.32 * energy, Voice: "pad"}}
}
