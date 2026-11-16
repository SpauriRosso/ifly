export type InstrumentId = 'piano' | 'guitar' | 'bass' | 'violin' | 'cello' | 'flute' | 'voice' | 'woodblock' | 'timpani'

// INSTRUMENT_LABELS gives the automatic score palette short, user-facing names.
export const INSTRUMENT_LABELS: Record<InstrumentId, string> = {
  piano: 'Piano', guitar: 'Guitare', bass: 'Contrebasse', violin: 'Violon',
  cello: 'Violoncelle', flute: 'Flûte', woodblock: 'Woodblock', timpani: 'Timbales',
  voice: 'Voix ooh',
}

export interface InstrumentSample {
  midi: number
  file: string
}

export interface InstrumentDefinition {
  path: string
  sustain: boolean
  samples: InstrumentSample[]
}

// INSTRUMENTS lists the local General MIDI sample anchors and whether an instrument can sustain a note.
export const INSTRUMENTS: Record<InstrumentId, InstrumentDefinition> = {
  piano: { path: 'salamander', sustain: false, samples: [
    { midi: 36, file: 'C2' }, { midi: 39, file: 'Ds2' }, { midi: 42, file: 'Fs2' }, { midi: 45, file: 'A2' },
    { midi: 48, file: 'C3' }, { midi: 51, file: 'Ds3' }, { midi: 54, file: 'Fs3' }, { midi: 57, file: 'A3' },
    { midi: 60, file: 'C4' }, { midi: 63, file: 'Ds4' }, { midi: 66, file: 'Fs4' }, { midi: 69, file: 'A4' },
    { midi: 72, file: 'C5' }, { midi: 75, file: 'Ds5' }, { midi: 78, file: 'Fs5' }, { midi: 81, file: 'A5' },
    { midi: 84, file: 'C6' }, { midi: 87, file: 'Ds6' }, { midi: 90, file: 'Fs6' },
  ] },
  guitar: { path: 'fluidr3/acoustic_guitar_nylon', sustain: false, samples: [
    { midi: 48, file: 'C3' }, { midi: 55, file: 'G3' }, { midi: 62, file: 'D4' },
    { midi: 69, file: 'A4' }, { midi: 76, file: 'E5' }, { midi: 83, file: 'B5' },
  ] },
  bass: { path: 'fluidr3/acoustic_bass', sustain: false, samples: [
    { midi: 36, file: 'C2' }, { midi: 43, file: 'G2' }, { midi: 50, file: 'D3' },
    { midi: 57, file: 'A3' }, { midi: 64, file: 'E4' },
  ] },
  violin: { path: 'fluidr3/violin', sustain: true, samples: [
    { midi: 60, file: 'C4' }, { midi: 67, file: 'G4' }, { midi: 74, file: 'D5' },
    { midi: 81, file: 'A5' }, { midi: 88, file: 'E6' },
  ] },
  cello: { path: 'fluidr3/cello', sustain: true, samples: [
    { midi: 48, file: 'C3' }, { midi: 55, file: 'G3' }, { midi: 62, file: 'D4' },
    { midi: 69, file: 'A4' }, { midi: 76, file: 'E5' },
  ] },
  flute: { path: 'fluidr3/flute', sustain: true, samples: [
    { midi: 60, file: 'C4' }, { midi: 67, file: 'G4' }, { midi: 74, file: 'D5' },
    { midi: 81, file: 'A5' }, { midi: 88, file: 'E6' },
  ] },
  voice: { path: 'fluidr3/voice_oohs', sustain: true, samples: [
    { midi: 60, file: 'C4' }, { midi: 67, file: 'G4' }, { midi: 74, file: 'D5' },
    { midi: 81, file: 'A5' }, { midi: 88, file: 'E6' },
  ] },
  woodblock: { path: 'fluidr3/woodblock', sustain: false, samples: [{ midi: 60, file: 'C4' }] },
  timpani: { path: 'fluidr3/timpani', sustain: false, samples: [{ midi: 48, file: 'C3' }] },
}

// nearestSample takes an instrument and MIDI pitch, then returns the closest local sample anchor for transposition.
export function nearestSample(instrument: InstrumentId, midi: number): InstrumentSample {
  const samples = INSTRUMENTS[instrument].samples
  let closest = samples[0]
  for (const sample of samples.slice(1)) {
    if (Math.abs(sample.midi - midi) < Math.abs(closest.midi - midi)) closest = sample
  }
  return closest
}

// sampleKey takes an instrument and anchor and returns the stable key used by the decoded audio cache.
export function sampleKey(instrument: InstrumentId, sample: InstrumentSample): string {
  return `${instrument}/${sample.file}`
}

// sampleUrl takes an instrument and anchor and returns its packaged offline MP3 URL.
export function sampleUrl(instrument: InstrumentId, sample: InstrumentSample): string {
  return `/samples/${INSTRUMENTS[instrument].path}/${sample.file}.mp3`
}
