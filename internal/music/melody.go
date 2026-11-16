package music

import "math"

// makeMotif takes the image digest and scale size, then returns a rising and falling contour without wrapping its peak to the tonic.
func makeMotif(digest [32]byte, scaleSize int) [4]int {
	second := 2 + int(digest[1]%2)
	peak := 4 + int(digest[2]%2)
	if peak >= scaleSize {
		peak = scaleSize - 1
	}
	return [4]int{0, second, peak, second}
}

// makeBridgeMotif takes later digest bytes and the scale size, then returns a contrasting high-reaching four-note idea for the B section.
func makeBridgeMotif(digest [32]byte, scaleSize int) [4]int {
	start := 1 + int(digest[3])%(scaleSize-2)
	peak := scaleSize - 1 - int(digest[4]%2)
	if peak <= start {
		peak = scaleSize - 1
	}
	return [4]int{start, peak, peak - 2, int(digest[5] % 2)}
}

// turnMotif takes a four-note idea and rotation count, then returns a related starting angle for a later episode.
func turnMotif(motif [4]int, turn int) [4]int {
	var rotated [4]int
	for index := range rotated {
		rotated[index] = motif[(index+turn)%len(motif)]
	}
	return rotated
}

// melodyRhythmForBar takes the style, visual complexity, and phrase position, then returns a safe variation of that bar's rhythm.
func melodyRhythmForBar(bar, section int, complexity float64, character musicalCharacter) []float64 {
	rhythm := append([]float64(nil), character.melodyRhythm[bar%2]...)
	if complexity < 0.28 && (character.mood == "calm" || character.mood == "neutral") && len(rhythm) > 2 {
		rhythm = rhythm[:2]
	}
	if section == 3 {
		if bar%2 == 1 {
			return nil
		}
		return rhythm[:1]
	}
	if bar%4 == 3 || (section == 5 && bar%4 >= 2) {
		return rhythm[:1]
	}
	for index := 1; index < len(rhythm); index++ {
		rhythm[index] += phraseOffbeats[section][bar%4]
	}
	if section == 2 && complexity >= 0.45 && len(rhythm) > 1 && rhythm[1]-rhythm[0] >= 1 {
		passingBeat := (rhythm[0] + rhythm[1]) / 2
		rhythm = append(rhythm[:1], append([]float64{passingBeat}, rhythm[1:]...)...)
	}
	return rhythm
}

// melodyDegree takes a motif and phrase position, then returns a scale degree after the section's contour transformation.
func melodyDegree(motif [4]int, bar, section, noteIndex, scaleSize int) int {
	position := bar % 4
	motifIndex := (position + noteIndex) % len(motif)
	if section == 2 {
		motifIndex = (position + 1 + len(motif) - noteIndex%len(motif)) % len(motif)
	} else if (section == 1 && position == 1) || (section == 4 && position == 2) {
		motifIndex = (position + len(motif) - noteIndex%len(motif)) % len(motif)
	}
	degree := motif[motifIndex]
	if noteIndex > 0 {
		degree += phraseDegreeShifts[section][position]
	}
	return (degree%scaleSize + scaleSize) % scaleSize
}

// melodyNotes takes one bar's harmony, motif, phrase position, and prior pitch, then returns its shaped lead events and final pitch.
func melodyNotes(start float64, bar, section, finalBar, root int, chord chordSpec, motif [4]int, previousPitch int, energy, complexity float64, character musicalCharacter) ([]MusicNote, int) {
	rhythm := melodyRhythmForBar(bar, section, complexity, character)
	notes := make([]MusicNote, 0, len(rhythm))
	chordRoot := root + chord.rootOffset
	targets := make([]int, len(rhythm))
	for index := range rhythm {
		degree := melodyDegree(motif, bar, section, index, len(character.scale))
		target := root + 24 + character.scale[degree]
		if index == 0 && bar%2 == 0 {
			target = chordRoot + 24
			if bar%4 == 2 {
				target += chord.third
			}
		} else if index == len(rhythm)-1 && (bar%4 == 3 || section == 5) {
			interval := chord.third
			if bar%2 == 1 {
				interval = 7
			}
			target = chordRoot + 24 + interval
		}
		if bar == finalBar {
			target = root + 24
		}
		targets[index] = target
	}
	pitches := harmonizeMelodyTargets(targets, rhythm, previousPitch, root, chordRoot, chord.third, character.scale)
	for index, beat := range rhythm {
		duration := 1.2
		if index+1 < len(rhythm) {
			duration = math.Min(1.5, rhythm[index+1]-beat-0.12)
		} else if bar == finalBar {
			duration = 1.6
		} else if bar%4 == 3 {
			duration = 2.2
		}
		velocity := 0.57 * energy
		if index == 0 {
			velocity *= 1.18
		} else if index == len(rhythm)-1 {
			velocity *= 0.82
		}
		notes = append(notes, MusicNote{Beat: start + beat, Duration: duration, MIDI: pitches[index], Velocity: velocity, Voice: "lead"})
	}
	if len(pitches) > 0 {
		previousPitch = pitches[len(pitches)-1]
	}
	return notes, previousPitch
}
