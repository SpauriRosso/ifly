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
  IFLY plans a shared harmonic skeleton before writing both piano hands.
  This is a limited rule-based adaptation, not a reproduction of that system.
- [Diatony (IJCAI 2024)](https://ijcai.org/proceedings/2024/858)
  describes a layered tonal workflow: global form, harmonic motion, voice
  leading, and ornamentation. IFLY now keeps those steps separate. Its voice
  leading is far simpler than the paper's four-voice constraint model.
- [Baroque Classical Music Generator](https://weroarforareason.itch.io/baroque)
  is a useful piano-centered reference for distinct styles, two-hand writing,
  score visibility, and modelled piano sound. The creator's
  [development note](https://weroarforareason.itch.io/baroque/devlog/1768803/baroque-composer-is-out)
  says it procedurally adds features using recursion, without AI or samples.
  That note does not disclose its musical rules. IFLY independently adds
  bounded passing-note elaboration to selected passages.
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
key, and motifs. A measure plan first fixes the chord progression and nearest
inversions. The writing then adds a right-hand phrase, left-hand bass and
gesture, and short inner chord attacks. A second chord inside selected bars
prepares the following dominant; both hands follow that change. B, B-prime,
and the climax use different progressions before a shared dominant cadence.
The left hand uses phrase-dependent flowing, Alberti, pulsed, contrapuntal,
and reflective figures instead of one figure in every bar. Selected right-hand
intervals receive stepwise passing notes. The sampled piano's two hands now
have a closer audio balance and share room ambience. Score generation and
playback remain offline.

Automated tests check chord resolution, melodic intervals, melodic density,
pitch variety, question-and-answer cadences, harmonic development, left-hand
gesture changes, short inner chords, hand separation, bundled sample coverage,
and the player's event limit. They cannot
judge whether a piece resembles a convincing Mozart, Chopin, or Beethoven
idiom. Those names remain artistic references, not labels for validated style
models. The next evaluation should compare recordings of a quiet landscape,
an abstract image, and a tense scene, focusing on melodic contour, harmonic
direction, left-hand writing, dynamics, and the audible joins between notes.
