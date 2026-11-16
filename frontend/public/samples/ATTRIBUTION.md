# Instrument samples

The audio files in `fluidr3/` are selected General MIDI samples from the
FluidR3 SoundFont by Frank Wen and contributors, rendered and hosted by
Benjamin Gleitzman in [midi-js-soundfonts](https://github.com/gleitz/midi-js-soundfonts/tree/gh-pages/FluidR3_GM).
The source repository lists these samples as released
under the [Creative Commons Attribution 3.0 license](https://creativecommons.org/licenses/by/3.0/).

IFLY bundles a small subset for offline playback. The sample file names and
instrument names are preserved. On Windows, run
`powershell -ExecutionPolicy Bypass -File scripts/fetch-instruments.ps1` to
retrieve the same selection from the source repository.

The piano files in `salamander/` are selected MP3 recordings from
[Salamander Grand Piano V3](https://github.com/sfzinstruments/SalamanderGrandPiano)
by Alexander Holm, distributed in MP3 form by the
[Tone.js audio repository](https://github.com/Tonejs/audio/tree/master/salamander).
Salamander Grand Piano V3 is licensed under
[Creative Commons Attribution 3.0](https://creativecommons.org/licenses/by/3.0/).
IFLY bundles 19 pitch anchors for offline use. On Windows, run
`powershell -ExecutionPolicy Bypass -File scripts/fetch-salamander.ps1` to
retrieve this selection.
