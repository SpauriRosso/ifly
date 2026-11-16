# Music generation: research and current design

IFLY turns an image into a deterministic score. Its current composer is a rule-based
system, not a trained music model. The rules are useful for controlling harmony,
form, and the relationship to the image, but they do not guarantee an expressive
performance. Listening remains the decisive check.

## Research applied

- [MusicFrameworks (ISMIR 2021)](https://archives.ismir.net/ismir2021/paper/000017.pdf)
  models a song at several levels: sections, phrases, chords, melody, and rhythm.
  IFLY likewise plans a long form before assigning bar harmony and note events.
- [AccoMontage2 (ISMIR 2022)](https://archives.ismir.net/ismir2022/paper/000029.pdf)
  evaluates accompaniment against a lead melody at note, phrase, and piece levels.
  IFLY places accompaniment after its lead and leaves lead attacks unobscured.
  This is a limited rule-based adaptation, not a reproduction of that system.
- [PULSE-IT (ISMIR 2025)](https://ismir2025program.ismir.net/lbd_403.html)
  treats wind articulation, including slurs and vibrato, as essential to expressive
  synthesis. IFLY now separates flute breath groups, sustains notes through slurs,
  and delays vibrato on longer notes.
- [Flute performance study](https://dael.euracoustics.org/bin/EAA/aaua_dl?document_id=52021)
  links breathing to the musical task. IFLY reserves silence between two-bar
  wind phrases rather than placing a new attack on every beat.
- [Baroque Classical Music Generator](https://weroarforareason.itch.io/baroque)
  is a useful piano-centered reference for distinct styles, two-hand writing,
  score visibility, and modelled piano sound. Its page describes the product;
  it does not publish a composition algorithm for IFLY to reproduce.
- [Parameterized Melody Generation with Autoencoders and Temporally-Consistent Noise](https://nime.org/proceedings/2019/nime2019_paper035.pdf)
  uses Perlin noise to vary an existing melody smoothly in a learned latent
  space. IFLY uses a much smaller, deterministic one-dimensional gradient-noise
  signal only for gradual phrase dynamics. It does not claim that noise alone
  generates harmony or that this implements the paper's trained model.

## Implemented checks

The wind melody searches candidate pitches by dynamic programming. Strong
positions use chord tones; weak inner positions may use scale tones that resolve
by step. The cost balances motif shape, register, connected motion, and leaps.
The audio renderer uses local flute samples with a sustained envelope, slur
transitions, and delayed vibrato. The score and playback remain offline.
Quiet scenes now put sampled piano in the lead, bass, and lower broken-chord
responses. The strings sit behind the piano; tense and energetic scenes retain
different lead instruments. Piano samples retain their natural decay rather
than receiving an extra exponential fade or a repeated echo.

Automated tests check chord resolution, melodic intervals, breath gaps, slur
classification, score bounds, and the complete Wails build. They cannot judge
whether the result feels poetic or whether a particular flute sample sounds
convincing. The next evaluation should compare recordings of at least a quiet
landscape, an abstract image, and a tense scene, focusing on melodic contour,
harmonic direction, instrumental balance, and the audible joins between notes.
