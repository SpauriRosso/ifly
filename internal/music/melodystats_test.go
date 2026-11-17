package music

import (
	"math"
	"testing"

	"ifly/internal/analysis"
	"ifly/internal/semantic"
)

// upperMelody takes a polyphonic piano score and returns the highest note at each lead onset.
func upperMelody(plan MusicPlan) []MusicNote {
	notes := []MusicNote{}
	for _, note := range plan.Notes {
		if note.Voice != "lead" {
			continue
		}
		if len(notes) == 0 || math.Abs(note.Beat-notes[len(notes)-1].Beat) > 0.05 {
			notes = append(notes, note)
		} else if note.MIDI > notes[len(notes)-1].MIDI {
			notes[len(notes)-1] = note
		}
	}
	return notes
}

// TestPianoMelodyFormsQuestionAndAnswer checks that the opening has a moving lead and a half cadence answered by a tonic cadence.
func TestPianoMelodyFormsQuestionAndAnswer(t *testing.T) {
	for _, scene := range []string{"A serene pond with water lilies.", "A threatening storm.", "A lively dancing scene."} {
		for center := byte(0); center < 12; center++ {
			var digest [32]byte
			for index := range digest {
				digest[index] = byte(index*13 + 7)
			}
			digest[0] = center
			plan := composeMusic(analysis.StructuralAnalysis{Complexity: 0.5}, semantic.SemanticAnalysis{Description: scene}, digest)
			melody := upperMelody(plan)
			opening := []MusicNote{}
			for _, note := range melody {
				if note.Beat >= 32 {
					break
				}
				opening = append(opening, note)
			}
			if len(opening) < 30 {
				t.Fatalf("%q center %d has only %d lead attacks in eight bars", scene, center, len(opening))
			}
			pitchCounts := map[int]int{}
			mostRepeated := 0
			for _, note := range opening {
				pitchCounts[note.MIDI]++
				if pitchCounts[note.MIDI] > mostRepeated {
					mostRepeated = pitchCounts[note.MIDI]
				}
			}
			if len(pitchCounts) < 7 || mostRepeated*3 > len(opening) {
				t.Fatalf("%q center %d repeats one pitch too often: %d pitches, %d/%d on the most common note", scene, center, len(pitchCounts), mostRepeated, len(opening))
			}
			var question, firstAnswer, answer *MusicNote
			for index := range opening {
				note := &opening[index]
				if note.Beat < 16 {
					question = note
				} else {
					if firstAnswer == nil {
						firstAnswer = note
					}
					answer = note
				}
			}
			root := 48 + int(center)
			if question == nil || firstAnswer == nil || answer == nil || question.MIDI%12 != (root+11)%12 || firstAnswer.MIDI%12 != root%12 || answer.MIDI%12 != root%12 {
				t.Fatalf("%q center %d lacks a resolved question-and-answer period: %v / %v / %v", scene, center, question, firstAnswer, answer)
			}
		}
	}
}
