package music

type pianoTexture int

const (
	pianoFlow pianoTexture = iota
	pianoAlberti
	pianoPulse
	pianoCounterpoint
	pianoStill
)

type pianoGesture struct {
	beats []float64
	tones []int
}

// pianoTextures chooses a coherent left-hand idiom for each scene and formal role.
var pianoTextures = map[string][6]pianoTexture{
	"calm":      {pianoFlow, pianoFlow, pianoCounterpoint, pianoStill, pianoFlow, pianoStill},
	"tense":     {pianoPulse, pianoAlberti, pianoCounterpoint, pianoStill, pianoPulse, pianoStill},
	"energetic": {pianoAlberti, pianoFlow, pianoCounterpoint, pianoStill, pianoAlberti, pianoStill},
	"neutral":   {pianoFlow, pianoAlberti, pianoCounterpoint, pianoStill, pianoFlow, pianoStill},
}

// pianoGestures defines bar-specific piano figures; the score writer maps tone indices onto voice-led chord notes.
var pianoGestures = map[pianoTexture][4]pianoGesture{
	pianoFlow: {
		{[]float64{0.5, 1.25, 2, 2.75, 3.5}, []int{0, 1, 2, 1, 0}},
		{[]float64{0.5, 1.5, 2.25, 3.25}, []int{2, 1, 0, 1}},
		{[]float64{0.5, 1.25, 1.75, 2.5, 3.25}, []int{0, 2, 1, 2, 1}},
		{[]float64{1, 2.5}, []int{1, 0}},
	},
	pianoAlberti: {
		{[]float64{0.5, 1, 1.5, 2, 2.5, 3, 3.5}, []int{0, 2, 1, 2, 0, 2, 1}},
		{[]float64{0.5, 1, 1.5, 2.25, 2.75, 3.5}, []int{1, 2, 0, 2, 1, 2}},
		{[]float64{0.5, 1, 1.5, 2, 2.5, 3, 3.5}, []int{0, 2, 1, 2, 1, 2, 0}},
		{[]float64{1, 2, 3}, []int{2, 1, 0}},
	},
	pianoPulse: {
		{[]float64{1, 1, 2.5, 2.5}, []int{0, 1, 1, 2}},
		{[]float64{0.75, 0.75, 2.75, 2.75}, []int{1, 2, 0, 1}},
		{[]float64{1, 1, 2, 2, 3.25}, []int{0, 2, 1, 2, 0}},
		{[]float64{1.5, 1.5, 3}, []int{1, 2, 0}},
	},
	pianoCounterpoint: {
		{[]float64{0.5, 1.25, 2, 2.75, 3.5}, []int{2, 1, 0, 1, 2}},
		{[]float64{0.25, 1, 1.75, 2.5, 3.25}, []int{2, 1, 0, 1, 2}},
		{[]float64{0.5, 1.25, 2, 2.75, 3.5}, []int{0, 1, 2, 1, 0}},
		{[]float64{0.5, 1.5, 2.5, 3.25}, []int{2, 1, 0, 1}},
	},
	pianoStill: {
		{[]float64{1, 3}, []int{0, 2}},
		{nil, nil},
		{[]float64{1.5, 3}, []int{1, 2}},
		{nil, nil},
	},
}

// pianoChordDurations and pianoChordVoices pace inner chord attacks instead of sustaining every inversion across the bar.
var pianoChordDurations = [6][4]float64{
	{1.5, 0.75, 1.2, 0.65},
	{1.4, 1.0, 1.3, 0.9},
	{1.15, 0.9, 1.2, 0.85},
	{1.4, 0, 0.9, 0},
	{1.55, 1.05, 1.35, 0.9},
	{1.5, 1.1, 0.8, 0.9},
}

var pianoChordVoices = [6][4]int{
	{2, 1, 1, 2},
	{2, 1, 2, 2},
	{2, 2, 2, 2},
	{1, 0, 1, 0},
	{2, 2, 2, 2},
	{2, 1, 1, 3},
}

// pianoBassResponses places optional fifths in the bass according to the scene and phrase role; zero leaves the bar open.
var pianoBassResponses = map[string][6]float64{
	"calm":      {0, 0, 2.5, 0, 2.5, 0},
	"tense":     {2, 2.5, 2.5, 0, 2.5, 0},
	"energetic": {2, 2.5, 2.5, 0, 2.5, 0},
	"neutral":   {0, 2.5, 2.5, 0, 2.5, 0},
}
