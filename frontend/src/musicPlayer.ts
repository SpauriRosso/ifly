import type { MusicNote, MusicPlan } from '../bindings/ifly/internal/music/models'
import { TIMBRES, type MusicVoice, type TimbreProfile } from './musicTimbres'

// seededNoise takes an image seed and returns a reproducible stream of audio samples between minus one and one.
function seededNoise(seed: string): () => number {
  let state = 2166136261
  for (const character of seed) state = Math.imul(state ^ character.charCodeAt(0), 16777619)
  return () => {
    state ^= state << 13
    state ^= state >>> 17
    state ^= state << 5
    return (state >>> 0) / 2147483648 - 1
  }
}

// createPulseBuffer takes an audio context and image seed, then returns a reproducible quarter-second noise source.
function createPulseBuffer(context: AudioContext, seed: string): AudioBuffer {
  const nextSample = seededNoise(seed)
  const buffer = context.createBuffer(1, Math.floor(context.sampleRate * 0.25), context.sampleRate)
  const samples = buffer.getChannelData(0)
  for (let index = 0; index < samples.length; index++) {
    samples[index] = nextSample()
  }
  return buffer
}

// createRoomImpulse takes an audio context, image seed, and decay length, then returns a stereo impulse for local ambience.
function createRoomImpulse(context: AudioContext, seed: string, seconds: number): AudioBuffer {
  const nextSample = seededNoise(seed + 'room')
  const buffer = context.createBuffer(2, Math.floor(context.sampleRate * seconds), context.sampleRate)
  for (let channel = 0; channel < 2; channel++) {
    const samples = buffer.getChannelData(channel)
    for (let index = 0; index < samples.length; index++) {
      const fadeIn = Math.min(1, index / (context.sampleRate * 0.008))
      const decay = (1 - index / samples.length) ** 2.5
      samples[index] = nextSample() * fadeIn * decay
    }
  }
  return buffer
}

// MusicPlayer schedules a bounded Go score with disposable Web Audio nodes and a pauseable audio clock.
export class MusicPlayer {
  private context: AudioContext | null = null
  private master: GainNode | null = null
  private echo: DelayNode | null = null
  private room: ConvolverNode | null = null
  private profile: TimbreProfile | null = null
  private pulseBuffer: AudioBuffer | null = null
  private timer: number | null = null
  private notes: MusicNote[] = []
  private secondsPerBeat = 0
  private loopStart = 0
  private loopDuration = 0
  private nextNote = 0

  // start takes a validated score, opens audio during the user's click, and begins its repeating phrase.
  async start(plan: MusicPlan): Promise<void> {
    if (this.context) throw new Error('Une lecture est déjà active.')
    if (!plan.notes?.length || plan.notes.length > 256 || plan.tempoBpm < 40 || plan.tempoBpm > 180 || plan.beatsPerLoop < 4 || plan.beatsPerLoop > 128) {
      throw new Error('Le plan musical est invalide.')
    }
    const notes = plan.notes
    const profile = TIMBRES[plan.timbre]
    if (!profile || notes.some((note, index) => !Number.isFinite(note.beat) || note.beat < 0 || note.beat >= plan.beatsPerLoop || !Number.isFinite(note.duration) || note.duration <= 0 || note.duration > 16 || !Number.isFinite(note.velocity) || note.velocity <= 0 || note.velocity > 1 || !Number.isInteger(note.midi) || note.midi < 24 || note.midi > 100 || !Object.prototype.hasOwnProperty.call(profile.voices, note.voice) || (index > 0 && note.beat < notes[index - 1].beat))) {
      throw new Error('Les notes du plan musical sont invalides.')
    }
    const context = new AudioContext()
    const master = context.createGain()
    const compressor = context.createDynamicsCompressor()
    const echo = context.createDelay(2)
    const feedback = context.createGain()
    const wet = context.createGain()
    const room = context.createConvolver()
    const roomWet = context.createGain()
    master.gain.value = 0.75
    compressor.threshold.value = -18
    compressor.ratio.value = 3
    echo.delayTime.value = profile.echoBeats * 60 / plan.tempoBpm
    feedback.gain.value = profile.echoFeedback
    wet.gain.value = profile.echoMix
    room.buffer = createRoomImpulse(context, plan.seed, profile.roomSeconds)
    roomWet.gain.value = profile.roomMix
    master.connect(compressor).connect(context.destination)
    echo.connect(feedback).connect(echo)
    echo.connect(wet).connect(master)
    room.connect(roomWet).connect(master)
    this.context = context
    this.master = master
    this.echo = echo
    this.room = room
    this.profile = profile
    this.pulseBuffer = createPulseBuffer(context, plan.seed)
    this.notes = notes
    this.secondsPerBeat = 60 / plan.tempoBpm
    this.loopDuration = plan.beatsPerLoop * this.secondsPerBeat
    this.loopStart = context.currentTime + 0.08
    this.nextNote = 0
    try {
      await context.resume()
      if (this.context !== context) return
      this.schedule()
      this.timer = window.setInterval(() => this.schedule(), 75)
    } catch (cause) {
      if (this.context === context) this.stop()
      throw cause
    }
  }

  // pause suspends the audio clock and returns after all scheduled voices have stopped advancing.
  async pause(): Promise<void> {
    if (this.context?.state === 'running') await this.context.suspend()
  }

  // resume restarts the suspended audio clock at its previous position and returns when scheduling resumes.
  async resume(): Promise<void> {
    if (this.context?.state === 'suspended') {
      await this.context.resume()
      this.schedule()
    }
  }

  // stop clears scheduling, fades the output, and releases the audio context for a new artwork.
  stop(): void {
    if (this.timer !== null) window.clearInterval(this.timer)
    this.timer = null
    const context = this.context
    const master = this.master
    this.context = null
    this.master = null
    this.echo = null
    this.room = null
    this.profile = null
    this.pulseBuffer = null
    this.notes = []
    if (!context) return
    if (context.state === 'running' && master) {
      master.gain.setTargetAtTime(0, context.currentTime, 0.015)
      window.setTimeout(() => void context.close().catch(() => undefined), 80)
    } else {
      void context.close().catch(() => undefined)
    }
  }

  // schedule places notes just ahead of the audio clock and returns without waiting for playback.
  private schedule(): void {
    const context = this.context
    if (!context || context.state !== 'running' || !this.notes.length) return
    const horizon = context.currentTime + 0.25
    while (this.loopStart + this.notes[this.nextNote].beat * this.secondsPerBeat <= horizon) {
      const note = this.notes[this.nextNote]
      const at = this.loopStart + note.beat * this.secondsPerBeat
      if (at >= context.currentTime - 0.02) this.playNote(note, at)
      this.nextNote++
      if (this.nextNote === this.notes.length) {
        this.nextNote = 0
        this.loopStart += this.loopDuration
      }
    }
  }

  // playNote takes one scheduled event, shapes its voice and echo, and disconnects the finished audio nodes.
  private playNote(note: MusicNote, at: number): void {
    const context = this.context
    const master = this.master
    const profile = this.profile
    if (!context || !master || !profile) return
    const voice = note.voice as MusicVoice
    if (voice === 'pulse') {
      this.playPulse(note, at)
      return
    }
    const timbre = profile.voices[voice]
    if (!timbre) return
    const oscillator = context.createOscillator()
    const envelope = context.createGain()
    const filter = context.createBiquadFilter()
    const duration = Math.max(0.1, Math.min(8, note.duration * this.secondsPerBeat))
    const attack = Math.min(timbre.attack, duration * 0.3)
    const release = Math.min(timbre.release, duration * 0.4)
    const level = note.velocity * timbre.gain
    const frequency = 440 * 2 ** ((note.midi - 69) / 12)
    oscillator.type = timbre.wave
    oscillator.frequency.setValueAtTime(frequency, at)
    filter.type = 'lowpass'
    filter.frequency.setValueAtTime(timbre.cutoff, at)
    if (voice === 'lead' || voice === 'bass' || voice === 'arp') filter.frequency.exponentialRampToValueAtTime(Math.max(180, timbre.cutoff * 0.65), at + duration)
    envelope.gain.setValueAtTime(0, at)
    envelope.gain.linearRampToValueAtTime(level, at + attack)
    if (voice === 'pad') {
      envelope.gain.setValueAtTime(level, at + duration - release)
    } else {
      envelope.gain.exponentialRampToValueAtTime(Math.max(0.0001, level * 0.38), at + duration - release)
    }
    envelope.gain.linearRampToValueAtTime(0, at + duration)
    oscillator.connect(envelope).connect(filter)
    filter.connect(master)
    if (this.echo && (voice === 'pad' || voice === 'lead' || voice === 'arp')) filter.connect(this.echo)
    if (this.room && (voice === 'pad' || voice === 'lead' || voice === 'arp')) filter.connect(this.room)
    let secondary: OscillatorNode | null = null
    let harmonicGain: GainNode | null = null
    if (voice === 'pad') {
      secondary = context.createOscillator()
      secondary.type = timbre.wave
      secondary.frequency.value = frequency
      secondary.detune.value = 5
      secondary.connect(envelope)
      secondary.start(at)
      secondary.stop(at + duration + 0.01)
    } else if (voice === 'lead') {
      secondary = context.createOscillator()
      harmonicGain = context.createGain()
      secondary.type = 'sine'
      secondary.frequency.value = frequency * 2
      harmonicGain.gain.value = 0.18
      secondary.connect(harmonicGain).connect(envelope)
      secondary.start(at)
      secondary.stop(at + duration + 0.01)
    }
    oscillator.onended = () => {
      oscillator.disconnect()
      secondary?.disconnect()
      harmonicGain?.disconnect()
      envelope.disconnect()
      filter.disconnect()
    }
    oscillator.start(at)
    oscillator.stop(at + duration + 0.01)
  }

  // playPulse takes one rhythmic event and plays a short filtered noise hit through the main output.
  private playPulse(note: MusicNote, at: number): void {
    const context = this.context
    const master = this.master
    const timbre = this.profile?.voices.pulse
    const buffer = this.pulseBuffer
    if (!context || !master || !timbre || !buffer) return
    const source = context.createBufferSource()
    const filter = context.createBiquadFilter()
    const envelope = context.createGain()
    const duration = Math.max(0.04, Math.min(0.16, note.duration * this.secondsPerBeat))
    source.buffer = buffer
    filter.type = 'highpass'
    filter.frequency.value = timbre.cutoff
    filter.Q.value = 0.7
    envelope.gain.setValueAtTime(0, at)
    envelope.gain.linearRampToValueAtTime(note.velocity * timbre.gain, at + Math.min(timbre.attack, duration * 0.2))
    envelope.gain.exponentialRampToValueAtTime(0.0001, at + duration)
    source.connect(filter).connect(envelope).connect(master)
    source.onended = () => {
      source.disconnect()
      filter.disconnect()
      envelope.disconnect()
    }
    source.start(at)
    source.stop(at + duration)
  }
}
