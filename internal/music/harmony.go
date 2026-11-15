package music

import (
	"math"
	"sort"
)

// voiceChord takes a triad and the prior three pitches, then returns a nearby inversion for smooth harmonic motion.
func voiceChord(root, third int, previous [3]int, hasPrevious bool) [3]int {
	base := [3]int{root + 12, root + 12 + third, root + 19}
	if !hasPrevious {
		return base
	}
	best := base
	bestCost := math.MaxInt
	for inversion := 0; inversion < 3; inversion++ {
		inverted := base
		for note := 0; note < inversion; note++ {
			inverted[note] += 12
		}
		sort.Ints(inverted[:])
		for octave := -12; octave <= 0; octave += 12 {
			candidate := [3]int{inverted[0] + octave, inverted[1] + octave, inverted[2] + octave}
			if candidate[0] < 48 || candidate[2] > 82 {
				continue
			}
			cost := 0
			for note := range candidate {
				cost += int(math.Abs(float64(candidate[note] - previous[note])))
			}
			if candidate[0] < 54 {
				cost += 54 - candidate[0]
			}
			if cost < bestCost {
				best = candidate
				bestCost = cost
			}
		}
	}
	return best
}
