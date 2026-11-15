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

// TIMBRES holds fixed synthesis settings so score events and playback logic remain independent of sound design data.
export const TIMBRES: Record<string, TimbreProfile> = {
  warm: {
    voices: {
      pad: { instrument: 'violin', gain: 0.42, attack: 0.12, release: 0.45, cutoff: 5200 },
      bass: { instrument: 'bass', gain: 0.65, attack: 0.012, release: 0.28, cutoff: 1800 },
      lead: { instrument: 'flute', gain: 0.55, attack: 0.035, release: 0.24, cutoff: 5800 },
      arp: { instrument: 'guitar', gain: 0.36, attack: 0.008, release: 0.18, cutoff: 5600 },
      hum: { instrument: 'voice', gain: 0.47, attack: 0.12, release: 0.42, cutoff: 4200 },
      pulse: { gain: 0.2, attack: 0.003, cutoff: 2200 },
    },
    echoBeats: 0.75, echoMix: 0.12, echoFeedback: 0.18, roomMix: 0.12, roomSeconds: 1.5,
  },
  dark: {
    voices: {
      pad: { instrument: 'cello', gain: 0.47, attack: 0.11, release: 0.42, cutoff: 3600 },
      bass: { instrument: 'bass', gain: 0.7, attack: 0.01, release: 0.22, cutoff: 1700 },
      lead: { instrument: 'violin', gain: 0.52, attack: 0.035, release: 0.18, cutoff: 4500 },
      arp: { instrument: 'piano', gain: 0.36, attack: 0.008, release: 0.14, cutoff: 4800 },
      hum: { instrument: 'voice', gain: 0.39, attack: 0.09, release: 0.4, cutoff: 3600 },
      pulse: { instrument: 'timpani', gain: 0.22, attack: 0.002, cutoff: 1600 },
    },
    echoBeats: 0.5, echoMix: 0.09, echoFeedback: 0.15, roomMix: 0.09, roomSeconds: 1.25,
  },
  bright: {
    voices: {
      pad: { instrument: 'piano', gain: 0.46, attack: 0.014, release: 0.38, cutoff: 6500 },
      bass: { instrument: 'bass', gain: 0.64, attack: 0.008, release: 0.2, cutoff: 2200 },
      lead: { instrument: 'guitar', gain: 0.6, attack: 0.008, release: 0.18, cutoff: 6200 },
      arp: { instrument: 'piano', gain: 0.37, attack: 0.006, release: 0.12, cutoff: 6500 },
      hum: { instrument: 'voice', gain: 0.38, attack: 0.08, release: 0.34, cutoff: 4800 },
      pulse: { instrument: 'woodblock', gain: 0.24, attack: 0.002, cutoff: 4200 },
    },
    echoBeats: 0.75, echoMix: 0.1, echoFeedback: 0.16, roomMix: 0.08, roomSeconds: 1.1,
  },
  soft: {
    voices: {
      pad: { instrument: 'piano', gain: 0.43, attack: 0.025, release: 0.5, cutoff: 4800 },
      bass: { instrument: 'bass', gain: 0.59, attack: 0.01, release: 0.26, cutoff: 1600 },
      lead: { instrument: 'flute', gain: 0.5, attack: 0.04, release: 0.23, cutoff: 4800 },
      arp: { instrument: 'guitar', gain: 0.32, attack: 0.008, release: 0.17, cutoff: 5000 },
      hum: { instrument: 'voice', gain: 0.44, attack: 0.13, release: 0.45, cutoff: 4000 },
      pulse: { gain: 0.17, attack: 0.003, cutoff: 2000 },
    },
    echoBeats: 0.75, echoMix: 0.11, echoFeedback: 0.16, roomMix: 0.13, roomSeconds: 1.7,
  },
}
