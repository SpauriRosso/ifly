package music

import (
	"math"
	"sort"
)

// harmonyForBar takes a style and bar position, then chooses its chord, prepares the opening and reprise cadences, and resolves the final bar.
func harmonyForBar(character musicalCharacter, bar, section, finalBar int) chordSpec {
	position := bar % 4
	chord := character.progression[position]
	if section == 2 || (section == 1 && position == 3) || (section == 4 && position == 2) {
		chord = character.bridge[position]
	}
	if section == 3 {
		chord = character.bridge[(position+2)%4]
	}
	if (section == 0 || section == 4) && position == 3 {
		chord = character.bridge[3]
	}
	if bar == finalBar {
		third := 4
		if character.mode == "minor" || character.mode == "dorian" {
			third = 3
		}
		chord = chordSpec{rootOffset: 0, third: third}
	}
	return chord
}

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
