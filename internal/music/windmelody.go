package music

import "math"

type windPitchCandidate struct {
	pitch     int
	chordTone bool
}

// windCandidates takes a chord, scale, and onset role, then returns playable pitches with passing tones only between harmonic anchors.
func windCandidates(tonalRoot, chordRoot, third int, scale []int, innerWeakBeat bool) []windPitchCandidate {
	chord := chordPitchClasses(chordRoot, third)
	candidates := make([]windPitchCandidate, 0, 22)
	for pitch := 67; pitch <= 88; pitch++ {
		if pitchBelongsToChord(pitch, chord) {
			candidates = append(candidates, windPitchCandidate{pitch: pitch, chordTone: true})
			continue
		}
		if !innerWeakBeat {
			continue
		}
		for _, degree := range scale {
			if (pitch-tonalRoot+120)%12 == degree {
				candidates = append(candidates, windPitchCandidate{pitch: pitch})
				break
			}
		}
	}
	return candidates
}

// windPitchCost takes a candidate and motif target, then scores melodic shape, range, and the tension of an optional passing tone.
func windPitchCost(candidate windPitchCandidate, target int) float64 {
	cost := windMelodyCosts.target * math.Abs(float64(candidate.pitch-target))
	if candidate.pitch%12 == target%12 {
		cost -= windMelodyCosts.motif
	}
	if candidate.pitch < 72 {
		cost += 0.5 * float64(72-candidate.pitch)
	}
	if candidate.pitch > 84 {
		cost += 0.5 * float64(candidate.pitch-84)
	}
	if !candidate.chordTone {
		cost += windMelodyCosts.passing
	}
	return cost
}

// windTransitionCost takes adjacent candidates and returns a finite cost only for stepwise, resolving wind-instrument movement.
func windTransitionCost(previous, current windPitchCandidate) float64 {
	motion := math.Abs(float64(current.pitch - previous.pitch))
	if !previous.chordTone && (!current.chordTone || motion > 2) {
		return math.Inf(1)
	}
	if !current.chordTone && motion > 4 {
		return math.Inf(1)
	}
	cost := windMelodyCosts.motion * motion
	if motion > 5 {
		cost += windMelodyCosts.leap * (motion - 5)
	}
	if motion == 0 {
		cost += windMelodyCosts.repeated
	}
	return cost
}

// harmonizeWindTargets takes a bar's melodic targets and returns the lowest-cost playable path through chord anchors and passing tones.
func harmonizeWindTargets(targets []int, rhythm []float64, previous, tonalRoot, chordRoot, third int, scale []int) []int {
	if len(targets) == 0 {
		return nil
	}
	candidates := make([][]windPitchCandidate, len(targets))
	costs := make([][]float64, len(targets))
	parents := make([][]int, len(targets))
	for index, target := range targets {
		weak := index > 0 && index+1 < len(targets) && math.Mod(rhythm[index], 1) != 0
		candidates[index] = windCandidates(tonalRoot, chordRoot, third, scale, weak)
		costs[index] = make([]float64, len(candidates[index]))
		parents[index] = make([]int, len(candidates[index]))
		for choice, candidate := range candidates[index] {
			costs[index][choice] = math.Inf(1)
			base := windPitchCost(candidate, target)
			if index == 0 {
				if previous != 0 {
					base += windTransitionCost(windPitchCandidate{pitch: previous, chordTone: true}, candidate)
				}
				costs[index][choice] = base
				continue
			}
			for prior, preceding := range candidates[index-1] {
				cost := costs[index-1][prior] + windTransitionCost(preceding, candidate) + base
				if cost < costs[index][choice] {
					costs[index][choice] = cost
					parents[index][choice] = prior
				}
			}
		}
	}
	last := len(targets) - 1
	choice := 0
	for index := range costs[last] {
		if costs[last][index] < costs[last][choice] {
			choice = index
		}
	}
	pitches := make([]int, len(targets))
	for index := last; index >= 0; index-- {
		pitches[index] = candidates[index][choice].pitch
		choice = parents[index][choice]
	}
	return pitches
}
