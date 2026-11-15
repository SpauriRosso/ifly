package music

import "math"

// expressNotes takes score events and an image digest, then adds bounded timing and intensity variation without moving chord or cadence anchors.
func expressNotes(notes []MusicNote, digest [32]byte) {
	for index := range notes {
		note := &notes[index]
		if note.Voice == "pad" {
			continue
		}
		intensity := float64(digest[(index*5+7)%len(digest)])/255 - 0.5
		note.Velocity = math.Min(1, math.Max(0.01, note.Velocity*(1+0.16*intensity)))
		if (note.Voice == "lead" || note.Voice == "arp") && note.Beat < 60 {
			timing := float64(digest[(index*11+13)%len(digest)])/255 - 0.5
			barStart := math.Floor(note.Beat/4) * 4
			note.Beat = math.Min(barStart+3.98, math.Max(barStart, note.Beat+0.04*timing))
		}
	}
}
