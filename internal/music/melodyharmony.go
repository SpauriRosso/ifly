package music

import "math"

// chordPitchClasses takes a chord root and third interval, then returns the pitch classes that form its triad.
func chordPitchClasses(root, third int) [3]int {
	return [3]int{root % 12, (root + third) % 12, (root + 7) % 12}
}

// pitchBelongsToChord takes a MIDI pitch and triad, then reports whether the pitch is a stable chord tone.
func pitchBelongsToChord(pitch int, chord [3]int) bool {
	for _, member := range chord {
		if pitch%12 == member {
			return true
		}
	}
	return false
}

// closestChordPitch takes a melodic target, prior pitch, and triad, then chooses a nearby chord tone in the lead register.
func closestChordPitch(target, previous, root, third int) int {
	best, bestCost := 0, math.MaxFloat64
	for octave := -2; octave <= 2; octave++ {
		for _, interval := range [3]int{0, third, 7} {
			candidate := root + 24 + interval + octave*12
			if candidate < 67 || candidate > 88 {
				continue
			}
			cost := 1.6 * math.Abs(float64(candidate-target))
			if previous != 0 {
				motion := math.Abs(float64(candidate - previous))
				cost += 0.7 * motion
				if motion > 12 {
					cost += 20
				}
				if candidate == previous {
					cost += 20
				}
			}
			if cost < bestCost {
				best, bestCost = candidate, cost
			}
		}
	}
	return best
}

// passingNeighbor takes two chord tones and a scale, then finds an optional weak-beat neighbor that resolves by step.
func passingNeighbor(target, previous, next, tonalRoot int, scale []int, chord [3]int) (int, bool) {
	best, bestCost := 0, math.MaxFloat64
	for octave := -1; octave <= 3; octave++ {
		for _, interval := range scale {
			candidate := tonalRoot + interval + octave*12
			if candidate < 67 || candidate > 88 || pitchBelongsToChord(candidate, chord) || candidate == previous {
				continue
			}
			if math.Abs(float64(candidate-previous)) > 4 || math.Abs(float64(candidate-next)) > 2 {
				continue
			}
			cost := math.Abs(float64(candidate-target)) + 0.35*math.Abs(float64(candidate-previous))
			if cost < bestCost {
				best, bestCost = candidate, cost
			}
		}
	}
	return best, best != 0
}

// harmonizeMelodyTargets takes one bar's motif targets and onsets, then returns chord-backed pitches with resolved weak-beat color notes.
func harmonizeMelodyTargets(targets []int, rhythm []float64, previous, tonalRoot, chordRoot, third int, scale []int) []int {
	pitches := make([]int, len(targets))
	for index, target := range targets {
		pitches[index] = closestChordPitch(target, previous, chordRoot, third)
		previous = pitches[index]
	}
	chord := chordPitchClasses(chordRoot, third)
	for index := 1; index+1 < len(pitches); index++ {
		if math.Mod(rhythm[index], 1) == 0 || (index+1 < len(pitches)-1 && math.Mod(rhythm[index+1], 1) != 0) {
			continue
		}
		if neighbor, ok := passingNeighbor(targets[index], pitches[index-1], pitches[index+1], tonalRoot, scale, chord); ok {
			pitches[index] = neighbor
		}
	}
	return pitches
}
