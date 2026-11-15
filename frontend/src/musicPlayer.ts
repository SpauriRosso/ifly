import type { MusicNote, MusicPlan } from '../bindings/changeme/models'

// MusicPlayer schedules a bounded Go score with disposable Web Audio nodes and a pauseable audio clock.
export class MusicPlayer {
  private context: AudioContext | null = null
  private master: GainNode | null = null
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
    if (notes.some((note, index) => !Number.isFinite(note.beat) || note.beat < 0 || note.beat >= plan.beatsPerLoop || !Number.isFinite(note.duration) || note.duration <= 0 || note.duration > 16 || (index > 0 && note.beat < notes[index - 1].beat))) {
      throw new Error('Les notes du plan musical sont invalides.')
    }
    const context = new AudioContext()
    const master = context.createGain()
    master.gain.value = 0.75
    master.connect(context.destination)
    this.context = context
    this.master = master
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

  // playNote turns one MIDI event into a short enveloped oscillator and disconnects it after release.
  private playNote(note: MusicNote, at: number): void {
    const context = this.context
    const master = this.master
    if (!context || !master || note.midi < 24 || note.midi > 100 || note.velocity <= 0 || note.velocity > 1) return
    const oscillator = context.createOscillator()
    const envelope = context.createGain()
    const duration = Math.max(0.1, Math.min(8, note.duration * this.secondsPerBeat))
    const pad = note.voice === 'pad'
    const attack = pad ? 0.09 : 0.015
    const release = Math.min(0.25, duration * 0.3)
    const level = note.velocity * (pad ? 0.12 : 0.18)
    oscillator.type = pad ? 'sine' : 'triangle'
    oscillator.frequency.value = 440 * 2 ** ((note.midi - 69) / 12)
    envelope.gain.setValueAtTime(0, at)
    envelope.gain.linearRampToValueAtTime(level, at + Math.min(attack, duration / 3))
    envelope.gain.setValueAtTime(level, at + duration - release)
    envelope.gain.linearRampToValueAtTime(0, at + duration)
    oscillator.connect(envelope).connect(master)
    oscillator.onended = () => {
      oscillator.disconnect()
      envelope.disconnect()
    }
    oscillator.start(at)
    oscillator.stop(at + duration + 0.01)
  }
}
