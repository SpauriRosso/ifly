import type { MusicNote } from '../bindings/ifly/internal/music/models'
import { INSTRUMENTS, nearestSample, sampleKey, sampleUrl, type InstrumentId, type InstrumentSample } from './musicInstruments'
import type { MusicVoice, TimbreProfile } from './musicTimbres'

export interface LoadedSample {
  buffer: AudioBuffer
  loopStart: number
  loopEnd: number
}

// prepareSustainLoop takes a decoded string or wind sample and returns a click-reduced loop for notes longer than the recording.
function prepareSustainLoop(context: AudioContext, original: AudioBuffer): LoadedSample {
  const rate = original.sampleRate
  const start = Math.floor(rate * 0.18)
  const fade = Math.floor(rate * 0.08)
  const end = Math.min(Math.floor(original.length * 0.78), original.length - fade)
  if (end <= start + 2 * fade) return { buffer: original, loopStart: 0, loopEnd: 0 }
  const buffer = context.createBuffer(original.numberOfChannels, end, rate)
  for (let channel = 0; channel < original.numberOfChannels; channel++) {
    const input = original.getChannelData(channel)
    const output = buffer.getChannelData(channel)
    output.set(input.subarray(0, end))
    for (let index = 0; index < fade; index++) {
      const blend = (index + 1) / fade
      output[end - fade + index] = input[end - fade + index] * (1 - blend) + input[start + index] * blend
    }
  }
  return { buffer, loopStart: (start + fade) / rate, loopEnd: end / rate }
}

// loadInstrumentSamples takes a bounded score and profile, fetches only required local MP3 anchors, and returns decoded audio.
export async function loadInstrumentSamples(context: AudioContext, notes: MusicNote[], profile: TimbreProfile, signal: AbortSignal): Promise<Map<string, LoadedSample>> {
  const requests = new Map<string, { instrument: InstrumentId; sample: InstrumentSample }>()
  for (const note of notes) {
    const instrument = profile.voices[note.voice as MusicVoice].instrument
    if (!instrument) continue
    const sample = nearestSample(instrument, note.midi)
    requests.set(sampleKey(instrument, sample), { instrument, sample })
  }
  const loaded = await Promise.all([...requests].map(async ([key, { instrument, sample }]) => {
    const response = await fetch(sampleUrl(instrument, sample), { signal })
    if (!response.ok) throw new Error(`Échantillon instrumental indisponible : ${instrument}/${sample.file}.`)
    const decoded = await context.decodeAudioData(await response.arrayBuffer())
    const prepared = INSTRUMENTS[instrument].sustain
      ? prepareSustainLoop(context, decoded)
      : { buffer: decoded, loopStart: 0, loopEnd: 0 }
    return [key, prepared] as const
  }))
  return new Map(loaded)
}
