import type { MusicSection } from '../bindings/ifly/internal/music/models'

export interface TimedSection extends MusicSection {
  startSeconds: number
}

export interface TempoMap {
  sections: TimedSection[]
  loopSeconds: number
}

// createTempoMap takes contiguous scored sections and their loop length, then returns validated time anchors for playback.
export function createTempoMap(sections: MusicSection[] | null, beatsPerLoop: number): TempoMap {
  if (!sections?.length || sections.length > 8) throw new Error('Le parcours de tempo est invalide.')
  const timed: TimedSection[] = []
  let expectedBeat = 0
  let elapsed = 0
  for (const section of sections) {
    if (!Number.isFinite(section.startBeat) || !Number.isFinite(section.endBeat) || section.startBeat !== expectedBeat || section.endBeat <= section.startBeat || !Number.isInteger(section.tempoBpm) || section.tempoBpm < 40 || section.tempoBpm > 180) {
      throw new Error('Le parcours de tempo est invalide.')
    }
    timed.push({ ...section, startSeconds: elapsed })
    elapsed += (section.endBeat - section.startBeat) * 60 / section.tempoBpm
    expectedBeat = section.endBeat
  }
  if (expectedBeat !== beatsPerLoop) throw new Error('Le parcours de tempo est incomplet.')
  return { sections: timed, loopSeconds: elapsed }
}

// secondsAtBeat takes a validated tempo map and beat position, then returns its exact elapsed playback time.
export function secondsAtBeat(map: TempoMap, beat: number): number {
  if (beat <= 0) return 0
  for (const section of map.sections) {
    if (beat < section.endBeat) return section.startSeconds + (beat - section.startBeat) * 60 / section.tempoBpm
  }
  return map.loopSeconds
}
