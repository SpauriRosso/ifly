package music

// pianoModeChordReplacements maps incompatible triads to consonant chords in each mode.
var pianoModeChordReplacements = map[string]map[chordSpec]chordSpec{
	"major": {
		{2, 4}: {2, 3},
	},
	"lydian": {
		{5, 4}: {2, 4},
		{2, 3}: {2, 4},
	},
	"mixolydian": {
		{7, 4}: {7, 3},
		{2, 4}: {2, 3},
		{0, 3}: {0, 4},
	},
	"minor": {
		{7, 4}: {7, 3},
		{2, 4}: {5, 3},
		{5, 4}: {5, 3},
	},
	"dorian": {
		{8, 4}: {10, 4},
		{7, 4}: {7, 3},
		{2, 4}: {2, 3},
	},
	"harmonic-minor": {
		{10, 4}: {8, 4},
		{3, 4}:  {5, 3},
		{2, 4}:  {5, 3},
		{0, 4}:  {0, 3},
	},
	"phrygian": {
		{10, 4}: {10, 3},
		{7, 4}:  {1, 4},
		{2, 4}:  {1, 4},
		{0, 4}:  {0, 3},
	},
}
