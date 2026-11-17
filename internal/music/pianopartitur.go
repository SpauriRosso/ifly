package music

import "math"

// pianoScoredPitch takes a chord tone, alteration, previous key, and intended motion, then returns a playable right-hand key.
func pianoScoredPitch(chordRoot, third, tone, alteration, previous, direction int) int {
	intervals := [3]int{0, third, 7}
	pitchClass := (chordRoot + intervals[tone] + alteration) % 12
	reference := previous
	if reference == 0 {
		reference = 76
	}
	best := 0
	bestCost := math.Inf(1)
	for pitch := 67; pitch <= 88; pitch++ {
		if pitch%12 != pitchClass {
			continue
		}
		cost := math.Abs(float64(pitch - reference))
		if cost > 7 {
			cost += (cost - 7) * 3
		}
		if previous != 0 && direction > 0 && pitch <= previous {
			cost += 5
		}
		if previous != 0 && direction < 0 && pitch >= previous {
			cost += 5
		}
		if previous != 0 && direction == 0 && pitch == previous {
			cost += 4
		}
		if cost < bestCost {
			best = pitch
			bestCost = cost
		}
	}
	return best
}

// pianoLowerScaleOffset takes a chord tone and the artwork's mode, then returns the nearest lower scale step.
func pianoLowerScaleOffset(chordRoot, third, tone, tonalRoot int, scale []int) int {
	intervals := [3]int{0, third, 7}
	base := (chordRoot + intervals[tone]) % 12
	for distance := 1; distance <= 2; distance++ {
		neighbor := (base - distance + 12) % 12
		for _, degree := range scale {
			if neighbor == (tonalRoot+degree)%12 {
				return -distance
			}
		}
	}
	return -1
}

// planPianoScoreLead takes harmony, form, visual complexity, and an artwork seed, then writes and ornaments playable right-hand figures.
func planPianoScoreLead(root, finalBar int, complexity float64, sections []MusicSection, digest [32]byte, character musicalCharacter, measures []pianoMeasure) [][]MusicNote {
	bars := make([][]MusicNote, finalBar+1)
	previous := 0
	lastBar := -1
	for episodeIndex := range sections {
		kind := pianoEpisodeFigures[episodeIndex]
		phrase := pianoPhraseFigures[kind]
		for position, figure := range phrase {
			if position < 2 && (kind == pianoTheme || kind == pianoVariation) {
				variation := 0
				if kind == pianoVariation {
					variation = 1
				}
				figure = pianoOpeningIdeas[digest[1]%3][variation][position]
			}
			bar := episodeIndex*4 + position
			measure := measures[bar]
			for index, written := range figure {
				step := written
				chord := measure.chordAt(step.beat)
				alteration := step.alteration
				if alteration < 0 {
					alteration = pianoLowerScaleOffset(root+chord.rootOffset, chord.third, step.tone, root, character.scale)
				}
				pitch := pianoScoredPitch(root+chord.rootOffset, chord.third, step.tone, alteration, previous, step.direction)
				end := 4.0
				if index+1 < len(figure) {
					end = figure[index+1].beat
				}
				duration := end - step.beat + 0.04
				if measure.lateChord != nil && step.beat < 2 && end > 2 {
					duration = 2 - step.beat - 0.06
				}
				if bar == finalBar {
					duration = 3.4
				}
				if pitch == previous && lastBar >= 0 && bar != finalBar {
					last := &bars[lastBar][len(bars[lastBar])-1]
					last.Duration = float64(bar*4) + end - last.Beat + 0.04
					continue
				}
				velocity := 0.53 * measure.energy
				if step.beat == 0 || step.beat == 2 {
					velocity *= 1.12
				}
				if step.alteration != 0 {
					velocity *= 0.78
				}
				bars[bar] = append(bars[bar], MusicNote{Beat: float64(bar*4) + step.beat, Duration: duration, MIDI: pitch, Velocity: velocity, Voice: "lead"})
				previous = pitch
				lastBar = bar
			}
		}
	}
	for bar := range bars {
		bars[bar] = ornamentPianoBar(bars[bar], measures[bar], complexity, root, character.scale)
	}
	return bars
}
