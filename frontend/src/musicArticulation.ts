import type { MusicNote } from '../bindings/ifly/internal/music/models'

export interface WindPhraseSettings {
  breathGroupBeats: number
  slurGapBeats: number
}

// findFluteSlurs takes ordered score events and wind settings, then returns lead notes joined within one breath group.
export function findFluteSlurs(notes: MusicNote[], settings: WindPhraseSettings): WeakSet<MusicNote> {
  const slurred = new WeakSet<MusicNote>()
  let previous: MusicNote | null = null
  for (const note of notes) {
    if (note.voice !== 'lead') continue
    if (previous) {
      const gap = note.beat - previous.beat - previous.duration
      const sameBreath = Math.floor(note.beat / settings.breathGroupBeats) === Math.floor(previous.beat / settings.breathGroupBeats)
      if (sameBreath && note.midi !== previous.midi && gap >= -0.18 && gap <= settings.slurGapBeats) slurred.add(note)
    }
    previous = note
  }
  return slurred
}
