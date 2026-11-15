export type InstrumentId = 'piano' | 'guitar' | 'bass' | 'violin' | 'cello' | 'flute' | 'woodblock' | 'timpani'

// INSTRUMENT_LABELS gives the automatic score palette short, user-facing names.
export const INSTRUMENT_LABELS: Record<InstrumentId, string> = {
  piano: 'Piano', guitar: 'Guitare', bass: 'Contrebasse', violin: 'Violon',
  cello: 'Violoncelle', flute: 'Flûte', woodblock: 'Woodblock', timpani: 'Timbales',
}

export interface InstrumentSample {
  midi: number
  file: string
}

export interface InstrumentDefinition {
  folder: string
  sustain: boolean
  samples: InstrumentSample[]
}

// INSTRUMENTS lists the local General MIDI sample anchors and whether an instrument can sustain a note.
export const INSTRUMENTS: Record<InstrumentId, InstrumentDefinition> = {
  piano: { folder: 'acoustic_grand_piano', sustain: false, samples: [
    { midi: 48, file: 'C3' }, { midi: 55, file: 'G3' }, { midi: 62, file: 'D4' },
    { midi: 69, file: 'A4' }, { midi: 76, file: 'E5' }, { midi: 83, file: 'B5' }, { midi: 91, file: 'G6' },
  ] },
  guitar: { folder: 'acoustic_guitar_nylon', sustain: false, samples: [
    { midi: 48, file: 'C3' }, { midi: 55, file: 'G3' }, { midi: 62, file: 'D4' },
    { midi: 69, file: 'A4' }, { midi: 76, file: 'E5' }, { midi: 83, file: 'B5' },
  ] },
  bass: { folder: 'acoustic_bass', sustain: false, samples: [
    { midi: 36, file: 'C2' }, { midi: 43, file: 'G2' }, { midi: 50, file: 'D3' },
    { midi: 57, file: 'A3' }, { midi: 64, file: 'E4' },
  ] },
  violin: { folder: 'violin', sustain: true, samples: [
    { midi: 60, file: 'C4' }, { midi: 67, file: 'G4' }, { midi: 74, file: 'D5' },
    { midi: 81, file: 'A5' }, { midi: 88, file: 'E6' },
  ] },
  cello: { folder: 'cello', sustain: true, samples: [
    { midi: 48, file: 'C3' }, { midi: 55, file: 'G3' }, { midi: 62, file: 'D4' },
    { midi: 69, file: 'A4' }, { midi: 76, file: 'E5' },
  ] },
  flute: { folder: 'flute', sustain: true, samples: [
    { midi: 60, file: 'C4' }, { midi: 67, file: 'G4' }, { midi: 74, file: 'D5' },
    { midi: 81, file: 'A5' }, { midi: 88, file: 'E6' },
  ] },
  woodblock: { folder: 'woodblock', sustain: false, samples: [{ midi: 60, file: 'C4' }] },
  timpani: { folder: 'timpani', sustain: false, samples: [{ midi: 48, file: 'C3' }] },
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
  return `/samples/fluidr3/${INSTRUMENTS[instrument].folder}/${sample.file}.mp3`
}
