package music

import "math"

type pianoStrikeKey struct {
	onset int
	pitch int
}

// pianoVoicePriority takes a score voice and returns which piano attack should survive an exact unison.
func pianoVoicePriority(voice string) int {
	switch voice {
	case "lead":
		return 4
	case "bass":
		return 3
	case "arp":
		return 2
	default:
		return 1
	}
}

// mergePianoUnisons takes both hands' score events and returns one physical key strike for each pitch and onset.
func mergePianoUnisons(notes []MusicNote) []MusicNote {
	strikes := make(map[pianoStrikeKey]int, len(notes))
	merged := make([]MusicNote, 0, len(notes))
	for _, note := range notes {
		key := pianoStrikeKey{onset: int(math.Round(note.Beat * 1000)), pitch: note.MIDI}
		if index, exists := strikes[key]; exists {
			if pianoVoicePriority(note.Voice) > pianoVoicePriority(merged[index].Voice) {
				merged[index] = note
			}
			continue
		}
		strikes[key] = len(merged)
		merged = append(merged, note)
	}
	return merged
}
