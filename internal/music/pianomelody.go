package music

import "math"

type pianoMelodyEvent struct {
	bar       int
	beat      float64
	chordRoot int
	third     int
	target    int
	strong    bool
	cadence   int
}

type pianoPitchChoice struct {
	pitch     int
	chordTone bool
}

// pianoRhythmForBar takes a scene mood, phrase role, and visual complexity, then returns one bar's melodic onsets.
func pianoRhythmForBar(mood string, role, position int, complexity float64) []float64 {
	if role == 3 {
		return pianoInterludeRhythms[position]
	}
	if role == 5 {
		return pianoEndingRhythms[position]
	}
	rhythm := pianoRhythms[mood][position]
	if role == 1 && pianoVariationRhythms[position] != nil {
		rhythm = pianoVariationRhythms[position]
	}
	if role == 2 && position < 3 {
		rhythm = append([]float64{rhythm[0], 0.25}, rhythm[1:]...)
	}
	if complexity >= 0.45 && (mood == "calm" || mood == "neutral") && position%2 == 0 && role != 5 && role != 3 {
		for index := 1; index < len(rhythm); index++ {
			if rhythm[index-1] < 0.5 && rhythm[index] > 0.5 {
				return append(append([]float64{}, rhythm[:index]...), append([]float64{0.5}, rhythm[index:]...)...)
			}
		}
	}
	return rhythm
}

// pianoDiatonicPitch takes a tonal base, scale, and degree and returns a pitch that can extend into a neighboring octave.
func pianoDiatonicPitch(base int, scale []int, degree int) int {
	octave := 0
	for degree < 0 {
		degree += len(scale)
		octave--
	}
	for degree >= len(scale) {
		degree -= len(scale)
		octave++
	}
	return base + scale[degree] + octave*12
}

// pianoPhraseEvents takes one episode, its harmony, visual complexity, and image seed, then returns ordered melody targets for the whole phrase.
func pianoPhraseEvents(episodeIndex, root, finalBar int, complexity float64, digest [32]byte, character musicalCharacter, measures []pianoMeasure) []pianoMelodyEvent {
	episode := formEpisodes[episodeIndex]
	events := []pianoMelodyEvent{}
	for position := 0; position < 4; position++ {
		bar := episodeIndex*4 + position
		for _, beat := range pianoRhythmForBar(character.mood, episode.role, position, complexity) {
			chord := measures[bar].chordAt(beat)
			events = append(events, pianoMelodyEvent{
				bar: bar, beat: beat, chordRoot: root + chord.rootOffset, third: chord.third,
				strong: math.Mod(beat, 1) == 0,
			})
		}
	}
	for index := range events {
		if index+1 == len(events) || events[index+1].bar != events[index].bar {
			events[index].strong = true
		}
	}
	base := root + 24
	if base > 76 {
		base -= 12
	}
	contour := pianoQuestionContour
	if episode.role == 1 {
		contour = pianoVariationContour
	} else if episode.role == 2 {
		contour = pianoAnswerContour
	}
	for index := range events {
		contourIndex := index * len(contour) / len(events)
		degree := int(digest[1]%3) + contour[contourIndex] + episode.motifTurn
		if contourIndex%4 == 2 || contourIndex%4 == 3 {
			degree += int(digest[2+contourIndex/4]%3) - 1
		}
		if episode.role == 3 || episode.role == 5 {
			degree--
		}
		events[index].target = pianoDiatonicPitch(base, character.scale, degree)
	}
	if len(events) > 0 {
		last := &events[len(events)-1]
		if last.bar == finalBar || (episode.role == 1 && last.chordRoot == root) {
			last.cadence = 1
		} else if episode.role != 3 && last.chordRoot-root == 7 && last.third == 4 {
			last.cadence = 2
		}
	}
	return events
}

// pianoPitchCandidates takes an onset and scale, then returns chord anchors or optional scale neighbors in the lead register.
func pianoPitchCandidates(event pianoMelodyEvent, tonalRoot int, scale []int) []pianoPitchChoice {
	chord := chordPitchClasses(event.chordRoot, event.third)
	choices := make([]pianoPitchChoice, 0, 22)
	for pitch := 67; pitch <= 88; pitch++ {
		chordTone := pitchBelongsToChord(pitch, chord)
		if event.cadence == 1 && pitch%12 != tonalRoot%12 {
			continue
		}
		if event.cadence == 2 && pitch%12 != (tonalRoot+11)%12 {
			continue
		}
		if event.strong && !chordTone {
			continue
		}
		if !chordTone {
			inScale := false
			for _, interval := range scale {
				if pitch%12 == (tonalRoot+interval)%12 {
					inScale = true
					break
				}
			}
			if !inScale {
				continue
			}
		}
		choices = append(choices, pianoPitchChoice{pitch: pitch, chordTone: chordTone})
	}
	return choices
}

// pianoTransitionCost takes adjacent pitch choices and returns a penalty for large leaps, repeated notes, or unresolved weak tones.
func pianoTransitionCost(previous, current pianoPitchChoice) float64 {
	motion := math.Abs(float64(current.pitch - previous.pitch))
	if !previous.chordTone && (!current.chordTone || motion > 2) {
		return math.Inf(1)
	}
	if !current.chordTone && motion > 4 {
		return math.Inf(1)
	}
	cost := pianoMelodyCosts.motion * motion
	if motion > 7 {
		cost += pianoMelodyCosts.leap * (motion - 7) * (motion - 7)
	}
	if motion == 0 {
		cost += pianoMelodyCosts.repeated
	}
	return cost
}

// harmonizePianoPhrase takes all phrase targets and a preceding pitch, then returns a globally connected chord-aware melody.
func harmonizePianoPhrase(events []pianoMelodyEvent, previousPitch, tonalRoot int, scale []int) []int {
	if len(events) == 0 {
		return nil
	}
	choices := make([][]pianoPitchChoice, len(events))
	costs := make([][]float64, len(events))
	parents := make([][]int, len(events))
	for index, event := range events {
		choices[index] = pianoPitchCandidates(event, tonalRoot, scale)
		costs[index] = make([]float64, len(choices[index]))
		parents[index] = make([]int, len(choices[index]))
		for choice, candidate := range choices[index] {
			base := pianoMelodyCosts.target * math.Abs(float64(candidate.pitch-event.target))
			if !candidate.chordTone {
				base += pianoMelodyCosts.passing
			}
			costs[index][choice] = math.Inf(1)
			if index == 0 {
				if previousPitch != 0 {
					base += pianoTransitionCost(pianoPitchChoice{pitch: previousPitch, chordTone: true}, candidate)
				}
				costs[index][choice] = base
				continue
			}
			for prior, preceding := range choices[index-1] {
				cost := costs[index-1][prior] + pianoTransitionCost(preceding, candidate) + base
				if cost < costs[index][choice] {
					costs[index][choice] = cost
					parents[index][choice] = prior
				}
			}
		}
	}
	last := len(events) - 1
	best := 0
	for index := range costs[last] {
		if costs[last][index] < costs[last][best] {
			best = index
		}
	}
	pitches := make([]int, len(events))
	for index := last; index >= 0; index-- {
		pitches[index] = choices[index][best].pitch
		best = parents[index][best]
	}
	return pitches
}

// planPianoMelody takes a harmonic skeleton, form, and image seed, then returns complete right-hand phrases indexed by bar.
func planPianoMelody(root, finalBar int, complexity float64, sections []MusicSection, digest [32]byte, character musicalCharacter, measures []pianoMeasure) [][]MusicNote {
	bars := make([][]MusicNote, finalBar+1)
	previous := 0
	for episodeIndex := range sections {
		events := pianoPhraseEvents(episodeIndex, root, finalBar, complexity, digest, character, measures)
		pitches := harmonizePianoPhrase(events, previous, root, character.scale)
		for index, event := range events {
			duration := math.Min(1.7, 4-event.beat-0.18)
			if index+1 < len(events) && events[index+1].bar == event.bar {
				duration = math.Min(1.7, events[index+1].beat-event.beat-0.1)
			}
			if event.bar == finalBar {
				duration = 2.1
			}
			energy := phraseEnergy(sections[episodeIndex].Energy*phraseDynamics[formEpisodes[episodeIndex].role][event.bar%4], event.bar, digest)
			velocity := 0.56 * energy
			if event.strong {
				velocity *= 1.1
			}
			bars[event.bar] = append(bars[event.bar], MusicNote{Beat: float64(event.bar*4) + event.beat, Duration: duration, MIDI: pitches[index], Velocity: velocity, Voice: "lead"})
		}
		if len(pitches) > 0 {
			previous = pitches[len(pitches)-1]
		}
	}
	return bars
}
