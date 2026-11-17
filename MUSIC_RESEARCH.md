# Music generation: research and current design

IFLY turns an image into a deterministic score. Its current composer is a rule-based
system, not a trained music model. The rules are useful for controlling harmony,
form, and the relationship to the image, but they do not guarantee an expressive
performance. Listening remains the decisive check.

## Research applied

- [MusicFrameworks (ISMIR 2021)](https://archives.ismir.net/ismir2021/paper/000017.pdf)
  models a song at several levels: sections, phrases, chords, melody, and rhythm.
  IFLY likewise plans a long form before assigning bar harmony and note events.
- [MeloForm (ISMIR 2022)](https://archives.ismir.net/ismir2022/paper/000068.pdf)
  builds motifs into phrases and sections with repetition and variation, then
  refines the rule-generated melody with a learned model. IFLY now plans each
  four-bar piano phrase before emitting notes, with a question ending on the
  dominant and an answer resolving to the tonic. It does not implement the
  learned refinement stage, so musical richness remains limited.
- [AccoMontage2 (ISMIR 2022)](https://archives.ismir.net/ismir2022/paper/000029.pdf)
  evaluates accompaniment against a lead melody at note, phrase, and piece levels.
  IFLY plans harmony and melody before writing the lower piano part.
  This is a limited rule-based adaptation, not a reproduction of that system.
- [Baroque Classical Music Generator](https://weroarforareason.itch.io/baroque)
  is a useful piano-centered reference for distinct styles, two-hand writing,
  score visibility, and modelled piano sound. Its page describes the product;
  it does not publish a composition algorithm for IFLY to reproduce.
- [Salamander Grand Piano V3](https://github.com/sfzinstruments/SalamanderGrandPiano)
  documents a multi-layer grand-piano recording. IFLY now bundles a compact
  19-anchor MP3 subset for offline playback, with its license in
  `frontend/public/samples/ATTRIBUTION.md`. This subset has fewer velocity
  layers and no sampled pedal or release noises, so its sound is less detailed
  than the original SFZ instrument or Baroque's advertised modelled piano.
- [Parameterized Melody Generation with Autoencoders and Temporally-Consistent Noise](https://nime.org/proceedings/2019/nime2019_paper035.pdf)
  uses Perlin noise to vary an existing melody smoothly in a learned latent
  space. IFLY uses a much smaller, deterministic one-dimensional gradient-noise
  signal only for gradual phrase dynamics. It does not claim that noise alone
  generates harmony or that this implements the paper's trained model.

## Implemented checks

Every active scene now produces solo piano. The image controls its mood, tempo,
key, and motifs. The score uses a chord plan and whole-phrase pitch search rather
than independent bar-level note choices. Strong melodic onsets land on chord
tones; passing notes resolve by step; recurring themes change rhythm and contour.
The score uses a right-hand melody with occasional
consonant dyads, a low bass line, and a more continuous broken-chord left hand.
The left-hand figure changes with the scene and becomes sparse in reflective
sections. Secondary dominants and dominant sevenths prepare selected cadences.
The player leaves a short release tail on piano notes and uses room ambience
without a repeated echo. Score generation and playback remain offline.

Automated tests check chord resolution, melodic intervals, melodic density,
pitch variety, question-and-answer cadences, hand separation, bundled sample
coverage, and the player's event limit. They cannot
judge whether a piece resembles a convincing Mozart, Chopin, or Beethoven
idiom. Those names remain artistic references, not labels for validated style
models. The next evaluation should compare recordings of a quiet landscape,
an abstract image, and a tense scene, focusing on melodic contour, harmonic
direction, left-hand writing, dynamics, and the audible joins between notes.
