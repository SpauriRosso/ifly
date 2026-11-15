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
	scale        []int
	progression  [4]chordSpec
	bridge       [4]chordSpec
	melodyRhythm [2][]float64
	bassBeats    []float64
	pulseBeats   []float64
	arpBeats     []float64
}

// musicalStyles keeps musical choices as data so the score generator only arranges and transforms them.
var musicalStyles = map[string]musicalCharacter{
	"calm": {
		mood: "calm", mode: "pentatonic", timbre: "warm", baseTempo: 64,
		scale:        []int{0, 2, 4, 7, 9},
		progression:  [4]chordSpec{{0, 4}, {5, 4}, {7, 4}, {0, 4}},
		bridge:       [4]chordSpec{{9, 3}, {5, 4}, {7, 4}, {7, 4}},
		melodyRhythm: [2][]float64{{0, 1.5, 3}, {0, 2.5}},
		bassBeats:    []float64{0},
		arpBeats:     []float64{1, 3},
	},
	"tense": {
		mood: "tense", mode: "minor", timbre: "dark", baseTempo: 98,
		scale:        []int{0, 2, 3, 5, 7, 8, 10},
		progression:  [4]chordSpec{{0, 3}, {8, 4}, {10, 4}, {7, 4}},
		bridge:       [4]chordSpec{{5, 3}, {8, 4}, {10, 4}, {7, 4}},
		melodyRhythm: [2][]float64{{0, 0.75, 2, 2.75}, {0, 1.5, 2.25, 3.25}},
		bassBeats:    []float64{0, 2}, pulseBeats: []float64{1, 3},
		arpBeats: []float64{0.5, 1.5, 2.5, 3.5},
	},
	"energetic": {
		mood: "energetic", mode: "major", timbre: "bright", baseTempo: 106,
		scale:        []int{0, 2, 4, 5, 7, 9, 11},
		progression:  [4]chordSpec{{0, 4}, {5, 4}, {7, 4}, {9, 3}},
		bridge:       [4]chordSpec{{9, 3}, {5, 4}, {7, 4}, {7, 4}},
		melodyRhythm: [2][]float64{{0, 0.75, 1.75, 3}, {0, 1, 2.5, 3.25}},
		bassBeats:    []float64{0, 2}, pulseBeats: []float64{1, 3},
		arpBeats: []float64{0.5, 1.5, 2.5, 3.5},
	},
	"neutral": {
		mood: "neutral", mode: "dorian", timbre: "soft", baseTempo: 78,
		scale:        []int{0, 2, 3, 5, 7, 9, 10},
		progression:  [4]chordSpec{{0, 3}, {5, 4}, {10, 4}, {0, 3}},
		bridge:       [4]chordSpec{{10, 4}, {5, 4}, {7, 3}, {7, 3}},
		melodyRhythm: [2][]float64{{0, 1.5, 3}, {0, 2}},
		bassBeats:    []float64{0, 2},
		arpBeats:     []float64{1, 3},
	},
}
