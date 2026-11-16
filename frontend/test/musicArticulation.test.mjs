import assert from 'node:assert/strict'
import test from 'node:test'
import { findFluteSlurs } from '../src/musicArticulation.ts'

test('flute slurs join notes within one breath and rearticulate after a breath or repeated pitch', () => {
  const notes = [
    { voice: 'lead', beat: 0, duration: 1.8, midi: 72 },
    { voice: 'pad', beat: 0.25, duration: 3, midi: 60 },
    { voice: 'lead', beat: 1.75, duration: 2.3, midi: 74 },
    { voice: 'lead', beat: 4, duration: 3.1, midi: 76 },
    { voice: 'lead', beat: 8.25, duration: 1.31, midi: 79 },
    { voice: 'lead', beat: 9.5, duration: 1.3, midi: 81 },
    { voice: 'lead', beat: 10.75, duration: 1.2, midi: 81 },
  ]
  const slurs = findFluteSlurs(notes, { breathGroupBeats: 8, slurGapBeats: 0.12 })
  assert.equal(slurs.has(notes[0]), false)
  assert.equal(slurs.has(notes[2]), true)
  assert.equal(slurs.has(notes[3]), true)
  assert.equal(slurs.has(notes[4]), false)
  assert.equal(slurs.has(notes[5]), true)
  assert.equal(slurs.has(notes[6]), false)
})
