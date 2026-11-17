import type { InstrumentId } from './musicInstruments'

export type MusicVoice = 'pad' | 'bass' | 'lead' | 'arp' | 'hum' | 'pulse'
type TonalVoice = Exclude<MusicVoice, 'pulse'>

export interface VoiceTimbre {
  instrument: InstrumentId
  gain: number
  attack: number
  release: number
  cutoff: number
}

export interface TimbreProfile {
  voices: Record<TonalVoice, VoiceTimbre> & { pulse: PulseTimbre }
  echoBeats: number
  echoMix: number
  echoFeedback: number
  roomMix: number
  roomSeconds: number
}

export interface PulseTimbre {
  instrument?: InstrumentId
  gain: number
  attack: number
  cutoff: number
}

// FLUTE_EXPRESSION keeps breath, slur, and vibrato settings separate from the audio scheduling logic.
export const FLUTE_EXPRESSION = {
  breathGroupBeats: 8,
  slurGapBeats: 0.12,
  sampleSkipSeconds: 0.09,
  attackSeconds: 0.11,
  slurAttackSeconds: 0.045,
  releaseSeconds: 0.16,
  vibratoDelaySeconds: 0.35,
  vibratoFrequencyHz: 5.2,
  vibratoDepthCents: 5,
} as const

// PIANO_VOICES separates the two piano hands and quiet inner chord notes in the sampled solo instrument.
const PIANO_VOICES: TimbreProfile['voices'] = {
  pad: { instrument: 'piano', gain: 0.4, attack: 0.012, release: 0.42, cutoff: 4900 },
  bass: { instrument: 'piano', gain: 0.5, attack: 0.006, release: 0.32, cutoff: 4800 },
  lead: { instrument: 'piano', gain: 0.82, attack: 0.006, release: 0.38, cutoff: 7200 },
  arp: { instrument: 'piano', gain: 0.6, attack: 0.006, release: 0.24, cutoff: 5900 },
  hum: { instrument: 'piano', gain: 0, attack: 0.01, release: 0.2, cutoff: 4000 },
  pulse: { gain: 0, attack: 0.003, cutoff: 2000 },
}

// TIMBRES holds scene-specific room balances while every active score voice uses the same grand piano.
export const TIMBRES: Record<string, TimbreProfile> = {
  warm: { voices: PIANO_VOICES, echoBeats: 0.75, echoMix: 0, echoFeedback: 0, roomMix: 0.15, roomSeconds: 1.6 },
  dark: { voices: PIANO_VOICES, echoBeats: 0.5, echoMix: 0, echoFeedback: 0, roomMix: 0.11, roomSeconds: 1.35 },
  bright: { voices: PIANO_VOICES, echoBeats: 0.75, echoMix: 0, echoFeedback: 0, roomMix: 0.08, roomSeconds: 1.1 },
  soft: { voices: PIANO_VOICES, echoBeats: 0.75, echoMix: 0, echoFeedback: 0, roomMix: 0.13, roomSeconds: 1.5 },
}
