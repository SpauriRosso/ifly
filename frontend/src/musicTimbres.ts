export type MusicVoice = 'pad' | 'bass' | 'lead' | 'arp' | 'pulse'

export interface VoiceTimbre {
  wave: OscillatorType
  gain: number
  attack: number
  release: number
  cutoff: number
}

export interface TimbreProfile {
  voices: Record<MusicVoice, VoiceTimbre>
  echoBeats: number
  echoMix: number
  echoFeedback: number
}

// TIMBRES holds fixed synthesis settings so score events and playback logic remain independent of sound design data.
export const TIMBRES: Record<string, TimbreProfile> = {
  warm: {
    voices: {
      pad: { wave: 'triangle', gain: 0.15, attack: 0.26, release: 0.55, cutoff: 1400 },
      bass: { wave: 'sine', gain: 0.22, attack: 0.012, release: 0.3, cutoff: 480 },
      lead: { wave: 'triangle', gain: 0.18, attack: 0.025, release: 0.24, cutoff: 2300 },
      arp: { wave: 'sine', gain: 0.11, attack: 0.008, release: 0.18, cutoff: 2600 },
      pulse: { wave: 'sine', gain: 0.11, attack: 0.005, release: 0.05, cutoff: 600 },
    },
    echoBeats: 0.75, echoMix: 0.12, echoFeedback: 0.18,
  },
  dark: {
    voices: {
      pad: { wave: 'sawtooth', gain: 0.09, attack: 0.18, release: 0.46, cutoff: 900 },
      bass: { wave: 'triangle', gain: 0.22, attack: 0.008, release: 0.22, cutoff: 520 },
      lead: { wave: 'triangle', gain: 0.17, attack: 0.012, release: 0.16, cutoff: 1700 },
      arp: { wave: 'triangle', gain: 0.1, attack: 0.006, release: 0.13, cutoff: 1700 },
      pulse: { wave: 'sine', gain: 0.14, attack: 0.004, release: 0.06, cutoff: 650 },
    },
    echoBeats: 0.5, echoMix: 0.09, echoFeedback: 0.15,
  },
  bright: {
    voices: {
      pad: { wave: 'triangle', gain: 0.13, attack: 0.12, release: 0.4, cutoff: 2900 },
      bass: { wave: 'triangle', gain: 0.2, attack: 0.007, release: 0.2, cutoff: 850 },
      lead: { wave: 'sawtooth', gain: 0.13, attack: 0.008, release: 0.16, cutoff: 3600 },
      arp: { wave: 'triangle', gain: 0.11, attack: 0.005, release: 0.12, cutoff: 4200 },
      pulse: { wave: 'sine', gain: 0.13, attack: 0.004, release: 0.06, cutoff: 950 },
    },
    echoBeats: 0.75, echoMix: 0.1, echoFeedback: 0.16,
  },
  soft: {
    voices: {
      pad: { wave: 'sine', gain: 0.17, attack: 0.22, release: 0.5, cutoff: 1600 },
      bass: { wave: 'sine', gain: 0.19, attack: 0.01, release: 0.27, cutoff: 500 },
      lead: { wave: 'triangle', gain: 0.16, attack: 0.018, release: 0.2, cutoff: 2100 },
      arp: { wave: 'sine', gain: 0.1, attack: 0.008, release: 0.17, cutoff: 2200 },
      pulse: { wave: 'sine', gain: 0.1, attack: 0.004, release: 0.05, cutoff: 600 },
    },
    echoBeats: 0.75, echoMix: 0.11, echoFeedback: 0.16,
  },
}
