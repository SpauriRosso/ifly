package music

type pianoChordSubstitution struct {
	phrase   int
	position int
	chord    chordSpec
}

type pianoPaletteChoice struct {
	mode    string
	scale   []int
	changes []pianoChordSubstitution
}

// pianoPaletteChoices supplies three coherent modal and harmonic colors for each scene theme; choice zero keeps the base style.
var pianoPaletteChoices = map[string][3]pianoPaletteChoice{
	"calm": {
		{},
		{mode: "lydian", scale: []int{0, 2, 4, 6, 7, 9, 11}, changes: []pianoChordSubstitution{{0, 2, chordSpec{2, 4}}, {1, 1, chordSpec{2, 4}}}},
		{mode: "mixolydian", scale: []int{0, 2, 4, 5, 7, 9, 10}, changes: []pianoChordSubstitution{{0, 1, chordSpec{10, 4}}, {2, 0, chordSpec{10, 4}}, {4, 1, chordSpec{10, 4}}}},
	},
	"energetic": {
		{},
		{mode: "lydian", scale: []int{0, 2, 4, 6, 7, 9, 11}, changes: []pianoChordSubstitution{{0, 1, chordSpec{2, 4}}, {1, 1, chordSpec{2, 4}}}},
		{mode: "mixolydian", scale: []int{0, 2, 4, 5, 7, 9, 10}, changes: []pianoChordSubstitution{{0, 2, chordSpec{10, 4}}, {2, 0, chordSpec{10, 4}}, {4, 1, chordSpec{10, 4}}}},
	},
	"tense": {
		{},
		{mode: "harmonic-minor", scale: []int{0, 2, 3, 5, 7, 8, 11}, changes: []pianoChordSubstitution{{0, 1, chordSpec{5, 3}}, {1, 1, chordSpec{8, 4}}}},
		{mode: "phrygian", scale: []int{0, 1, 3, 5, 7, 8, 10}, changes: []pianoChordSubstitution{{0, 1, chordSpec{1, 4}}, {2, 0, chordSpec{1, 4}}, {3, 0, chordSpec{1, 4}}}},
	},
	"melancholy": {
		{},
		{mode: "dorian", scale: []int{0, 2, 3, 5, 7, 9, 10}, changes: []pianoChordSubstitution{{0, 2, chordSpec{5, 4}}, {1, 1, chordSpec{5, 4}}, {3, 1, chordSpec{5, 4}}}},
		{mode: "harmonic-minor", scale: []int{0, 2, 3, 5, 7, 8, 11}, changes: []pianoChordSubstitution{{0, 1, chordSpec{5, 3}}, {2, 0, chordSpec{3, 4}}}},
	},
	"neutral": {
		{},
		{mode: "minor", scale: []int{0, 2, 3, 5, 7, 8, 10}, changes: []pianoChordSubstitution{{0, 1, chordSpec{8, 4}}, {0, 2, chordSpec{5, 3}}, {1, 1, chordSpec{5, 3}}}},
		{mode: "mixolydian", scale: []int{0, 2, 4, 5, 7, 9, 10}, changes: []pianoChordSubstitution{{0, 0, chordSpec{0, 4}}, {0, 2, chordSpec{10, 4}}, {1, 0, chordSpec{0, 4}}}},
	},
}
