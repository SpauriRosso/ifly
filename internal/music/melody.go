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
	if section == 1 && len(rhythm) > 1 && rhythm[len(rhythm)-1]-rhythm[len(rhythm)-2] > 0.6 {
		rhythm[len(rhythm)-1] -= 0.25
	}
	if section == 2 && complexity >= 0.45 && len(rhythm) > 1 && rhythm[1]-rhythm[0] >= 1 {
		passingBeat := (rhythm[0] + rhythm[1]) / 2
		rhythm = append(rhythm[:1], append([]float64{passingBeat}, rhythm[1:]...)...)
	}
	return rhythm
}

// nearestMelodyRegister takes a pitch and prior lead note, then returns a nearby upper-register octave distinct from the accompaniment.
func nearestMelodyRegister(pitch, previous int) int {
	if previous == 0 {
		previous = 76
	}
	best := 0
	bestDistance := math.MaxInt
	for octave := -24; octave <= 24; octave += 12 {
		candidate := pitch + octave
		if candidate < 67 || candidate > 88 {
			continue
		}
		distance := candidate - previous
		if distance < 0 {
			distance = -distance
		}
		if distance < bestDistance {
			best, bestDistance = candidate, distance
		}
	}
	if best == 0 {
		return pitch
	}
	return best
}

// melodyNotes takes one bar's harmony, motif, phrase position, and prior pitch, then returns its shaped lead events and final pitch.
func melodyNotes(start float64, bar, section, finalBar, root int, chord chordSpec, motif [4]int, previousPitch int, energy, complexity float64, character musicalCharacter) ([]MusicNote, int) {
	rhythm := melodyRhythmForBar(bar, section, complexity, character)
	notes := make([]MusicNote, 0, len(rhythm))
	chordRoot := root + chord.rootOffset
	for index, beat := range rhythm {
		degree := motif[(bar%4+index)%len(motif)]
		if section == 1 && index > 0 {
			degree = (degree + 1) % len(character.scale)
		}
		if section == 4 && index > 0 {
			degree = (degree + 2) % len(character.scale)
		}
		if section == 5 && index > 0 {
			degree = (degree + len(character.scale) - 1) % len(character.scale)
		}
		pitch := root + 24 + character.scale[degree]
		if index == 0 && bar%2 == 0 {
			pitch = chordRoot + 24
		} else if index == len(rhythm)-1 && (bar%4 == 3 || section == 5) {
			interval := chord.third
			if bar%2 == 1 {
				interval = 7
			}
			pitch = chordRoot + 24 + interval
		}
		if bar == finalBar {
			pitch = root + 24
		}
		pitch = nearestMelodyRegister(pitch, previousPitch)
		if pitch == previousPitch && bar != finalBar {
			if index == 0 {
				pitch = nearestMelodyRegister(chordRoot+24+chord.third, previousPitch)
			} else {
				pitch = nearestMelodyRegister(root+24+character.scale[(degree+1)%len(character.scale)], previousPitch)
			}
		}
		previousPitch = pitch
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
		notes = append(notes, MusicNote{Beat: start + beat, Duration: duration, MIDI: pitch, Velocity: velocity, Voice: "lead"})
	}
	return notes, previousPitch
}
