import assert from 'node:assert/strict'
import { statSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import test from 'node:test'
import { nearestSample, sampleUrl } from '../src/musicInstruments.ts'

test('every piano pitch uses a packaged nearby Salamander recording', () => {
  const seen = new Set()
  for (let midi = 36; midi <= 90; midi++) {
    const sample = nearestSample('piano', midi)
    assert.ok(Math.abs(sample.midi - midi) <= 2, `large pitch shift at MIDI ${midi}`)
    const url = sampleUrl('piano', sample)
    assert.match(url, /^\/samples\/salamander\/\w+\.mp3$/)
    if (seen.has(url)) continue
    seen.add(url)
    const path = fileURLToPath(new URL(`../public${url}`, import.meta.url))
    assert.ok(statSync(path).size > 20000, `missing or truncated sample: ${url}`)
  }
  assert.equal(seen.size, 19)
})
