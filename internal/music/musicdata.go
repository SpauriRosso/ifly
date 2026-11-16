package music

type chordSpec struct {
	rootOffset int
	third      int
}

type musicalCharacter struct {
	mood         string
	mode         string
	timbre       string
	baseTempo    int
	tempoArc     [6]int
	scale        []int
	progression  [4]chordSpec
	bridge       [4]chordSpec
	melodyRhythm [2][]float64
	windLead     bool
	pianoLead    bool
	bassBeats    []float64
	pulseBeats   []float64
	arpBeats     []float64
	colorTone    int
}

// windPhraseRhythm holds two-bar breath-group gestures for a flute lead instead of the keyboard-like repeating cells.
var windPhraseRhythm = [4][]float64{
	{0, 1.75},
	{0, 1.15, 2.25},
	{0.25, 2.0},
	{0},
}

// windMelodyCosts weights contour fidelity against stepwise motion and resolved passing tones in the wind-specific search.
var windMelodyCosts = struct {
	target   float64
	motion   float64
	leap     float64
	repeated float64
	passing  float64
	motif    float64
}{target: 2.0, motion: 0.65, leap: 3.0, repeated: 50.0, passing: 1.2, motif: 9.0}

type formEpisode struct {
	label      string
	role       int
	energy     float64
	tempoShift int
	motifTurn  int
}

// formEpisodes defines a longer one-way arc with recognizable returns and new motif angles instead of copying the first pass.
var formEpisodes = [...]formEpisode{
	{"A", 0, 0.72, 0, 0},
	{"A'", 1, 0.88, 0, 0},
	{"B", 2, 1.00, 0, 0},
	{"Interlude", 3, 0.42, 0, 0},
	{"Retour", 4, 0.92, 0, 0},
	{"Développement", 1, 0.86, 1, 1},
	{"B'", 2, 0.96, -2, 1},
	{"Respiration", 3, 0.38, -4, 1},
	{"Reprise", 4, 0.90, 0, 1},
	{"Climax", 2, 1.00, 2, 2},
	{"Silence", 3, 0.34, -3, 2},
	{"Dernier retour", 4, 0.84, -1, 0},
	{"Écho", 0, 0.66, -2, 2},
	{"Suspension", 3, 0.28, -5, 0},
	{"Dénouement", 5, 0.58, -3, 1},
	{"Conclusion", 5, 0.43, -5, 0},
}

// phraseOffbeats holds small phrase-level timing changes while keeping the first theme and its return recognizable.
var phraseOffbeats = [6][4]float64{
	{0, 0, 0, 0},
	{-0.25, -0.25, 0.25, 0},
	{0.25, 0, -0.25, 0},
	{0, 0, 0, 0},
	{0, 0.25, -0.25, 0},
	{-0.25, 0, 0, 0},
}

// phraseDegreeShifts keeps each section's sequence distinct while reserving the unchanged first bar for the returning theme.
var phraseDegreeShifts = [6][4]int{
	{0, 0, 0, 0},
	{1, 1, 2, 0},
	{0, 2, 1, 0},
	{0, 0, 0, 0},
	{0, 1, 2, 0},
	{0, -1, -2, 0},
}

// phraseDynamics shapes intensity within each four-bar phrase so each section rises and releases before the next.
var phraseDynamics = [6][4]float64{
	{0.82, 0.96, 1.08, 0.78},
	{0.90, 1.04, 1.12, 0.82},
	{0.94, 1.08, 1.16, 0.86},
	{0.80, 0.64, 0.72, 0.58},
	{0.90, 1.04, 1.14, 0.82},
	{0.86, 0.72, 0.54, 0.36},
}

// padSustainBeats leaves audible space before phrase answers and cadences instead of holding every chord for the full bar.
var padSustainBeats = [4]float64{3.4, 2.6, 3.1, 2.1}

// musicalStyles keeps musical choices as data so the score generator only arranges and transforms them.
var musicalStyles = map[string]musicalCharacter{
	"calm": {
		mood: "calm", mode: "major", timbre: "warm", baseTempo: 64,
		pianoLead:    true,
		tempoArc:     [6]int{0, 3, 6, -5, 4, -3},
		scale:        []int{0, 2, 4, 5, 7, 9, 11},
		progression:  [4]chordSpec{{0, 4}, {9, 3}, {5, 4}, {7, 4}},
		bridge:       [4]chordSpec{{9, 3}, {5, 4}, {2, 3}, {7, 4}},
		melodyRhythm: [2][]float64{{0, 1.5, 3}, {0, 2.5}},
		bassBeats:    []float64{0},
		arpBeats:     []float64{0.75, 1.75, 2.75, 3.25},
		colorTone:    14,
	},
	"tense": {
		mood: "tense", mode: "minor", timbre: "dark", baseTempo: 98,
		pianoLead:    true,
		tempoArc:     [6]int{0, 4, 9, -7, 6, -4},
		scale:        []int{0, 2, 3, 5, 7, 8, 10},
		progression:  [4]chordSpec{{0, 3}, {8, 4}, {10, 4}, {7, 4}},
		bridge:       [4]chordSpec{{5, 3}, {8, 4}, {10, 4}, {7, 4}},
		melodyRhythm: [2][]float64{{0, 0.75, 2, 2.75}, {0, 1.5, 2.25, 3.25}},
		bassBeats:    []float64{0, 2}, pulseBeats: []float64{1, 3},
		arpBeats: []float64{0.5, 1.5, 2.5, 3.5},
	},
	"energetic": {
		mood: "energetic", mode: "major", timbre: "bright", baseTempo: 106,
		pianoLead:    true,
		tempoArc:     [6]int{0, 5, 10, -6, 7, -2},
		scale:        []int{0, 2, 4, 5, 7, 9, 11},
		progression:  [4]chordSpec{{0, 4}, {5, 4}, {7, 4}, {9, 3}},
		bridge:       [4]chordSpec{{9, 3}, {5, 4}, {7, 4}, {7, 4}},
		melodyRhythm: [2][]float64{{0, 0.75, 1.75, 3}, {0, 1, 2.5, 3.25}},
		bassBeats:    []float64{0, 2}, pulseBeats: []float64{1, 3},
		arpBeats: []float64{0.5, 1.5, 2.5, 3.5},
	},
	"neutral": {
		mood: "neutral", mode: "dorian", timbre: "soft", baseTempo: 78,
		pianoLead:    true,
		tempoArc:     [6]int{0, 3, 7, -5, 5, -3},
		scale:        []int{0, 2, 3, 5, 7, 9, 10},
		progression:  [4]chordSpec{{0, 3}, {5, 4}, {10, 4}, {0, 3}},
		bridge:       [4]chordSpec{{10, 4}, {5, 4}, {7, 3}, {7, 3}},
		melodyRhythm: [2][]float64{{0, 1.5, 3}, {0, 2}},
		bassBeats:    []float64{0, 2},
		arpBeats:     []float64{0.75, 1.75, 2.75, 3.25},
	},
}
