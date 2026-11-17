package music

import "math"

// pianoOrnamentPitch takes two structural pitches and the current chord, then returns a scale or chord step that connects them.
func pianoOrnamentPitch(first, last, tonalRoot int, chord chordSpec, scale []int) (int, bool) {
	low, high := first, last
	if low > high {
		low, high = high, low
	}
	members := chordPitchClasses(tonalRoot+chord.rootOffset, chord.third)
	best, bestCost := 0, math.Inf(1)
	for pitch := low + 1; pitch < high; pitch++ {
		if math.Abs(float64(pitch-last)) > 2 || math.Abs(float64(pitch-first)) > 5 {
			continue
		}
		inScale := false
		for _, degree := range scale {
			if pitch%12 == (tonalRoot+degree)%12 {
				inScale = true
				break
			}
		}
		if !inScale && !pitchBelongsToChord(pitch, members) {
			continue
		}
		cost := math.Abs(float64(first+last)/2 - float64(pitch))
		if cost < bestCost {
			best, bestCost = pitch, cost
		}
	}
	return best, !math.IsInf(bestCost, 1)
}

// subdividePianoInterval takes a note pair and depth, then returns bounded passing notes that elaborate their interval recursively.
func subdividePianoInterval(first, last MusicNote, depth, tonalRoot int, measure pianoMeasure, scale []int) []MusicNote {
	gap := last.Beat - first.Beat
	if depth == 0 || gap < 0.7 {
		return nil
	}
	beat := first.Beat + gap/2
	chord := measure.chordAt(beat - float64(measure.bar*4))
	pitch, ok := pianoOrnamentPitch(first.MIDI, last.MIDI, tonalRoot, chord, scale)
	if !ok {
		return nil
	}
	between := MusicNote{Beat: beat, Duration: math.Min(0.5, gap/2-0.08), MIDI: pitch, Velocity: 0.76 * math.Min(first.Velocity, last.Velocity), Voice: "lead"}
	notes := subdividePianoInterval(first, between, depth-1, tonalRoot, measure, scale)
	notes = append(notes, between)
	notes = append(notes, subdividePianoInterval(between, last, depth-1, tonalRoot, measure, scale)...)
	return notes
}

// ornamentPianoBar takes a structural right-hand bar and its formal role, then elaborates one selected gesture without obscuring cadences.
func ornamentPianoBar(lead []MusicNote, measure pianoMeasure, complexity float64, tonalRoot int, scale []int) []MusicNote {
	if len(lead) < 2 || measure.role == 3 || measure.role == 5 {
		return lead
	}
	depth := 0
	if measure.role == 2 || measure.role == 4 || (measure.role == 1 && measure.bar%4 == 2) || (complexity > 0.7 && measure.bar%4 == 2) {
		depth = 1
	}
	if measure.role == 2 && measure.bar/4 == 9 && complexity > 0.4 {
		depth = 2
	}
	if depth == 0 {
		return lead
	}
	bestGap, bestIndex := 0.0, -1
	for index := 0; index+1 < len(lead); index++ {
		gap := lead[index+1].Beat - lead[index].Beat
		if gap > bestGap && gap >= 0.7 && math.Abs(float64(lead[index+1].MIDI-lead[index].MIDI)) >= 3 {
			bestGap, bestIndex = gap, index
		}
	}
	if bestIndex < 0 {
		return lead
	}
	ornaments := subdividePianoInterval(lead[bestIndex], lead[bestIndex+1], depth, tonalRoot, measure, scale)
	if len(ornaments) == 0 {
		return lead
	}
	lead[bestIndex].Duration = ornaments[0].Beat - lead[bestIndex].Beat - 0.06
	result := make([]MusicNote, 0, len(lead)+len(ornaments))
	result = append(result, lead[:bestIndex+1]...)
	result = append(result, ornaments...)
	result = append(result, lead[bestIndex+1:]...)
	return result
}
