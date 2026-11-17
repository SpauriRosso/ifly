package music

// Piano phrase data defines complete four-bar rhythmic gestures and two contrasting melodic ideas.
var pianoRhythms = map[string][4][]float64{
	"calm": {
		{0, 1, 1.75, 2.5, 3.25}, {0, 1.5, 2.25, 3.25},
		{0, 0.75, 1.5, 2.5, 3.25}, {0, 1, 2, 3},
	},
	"tense": {
		{0, 0.75, 1.5, 2.25, 3.25}, {0, 1, 1.75, 2.5, 3.25},
		{0, 0.75, 1.5, 2.25, 3}, {0, 1, 1.75, 2.75},
	},
	"energetic": {
		{0, 0.5, 1, 1.75, 2.5, 3.25}, {0, 0.75, 1.5, 2, 3},
		{0, 0.5, 1.25, 2, 2.75, 3.5}, {0, 1, 1.75, 2.75},
	},
	"neutral": {
		{0, 1, 2, 3}, {0, 1.25, 2.25, 3.25},
		{0, 0.75, 1.75, 2.5, 3.25}, {0, 1, 2, 3},
	},
}

var pianoInterludeRhythms = [4][]float64{{0, 2}, nil, {0, 2.5}, nil}
var pianoEndingRhythms = [4][]float64{{0, 1.5, 3}, {0, 2}, {0, 2.5}, {0}}
var pianoVariationRhythms = [4][]float64{nil, {0, 1, 2.5, 3.25}, {0, 0.5, 1.5, 2.5, 3.25}, nil}

// pianoEpisodeRhythms gives later returns their own four-bar gestures while preserving phrase arrivals.
var pianoEpisodeRhythms = map[int][4][]float64{
	5:  {{0, 0.75, 1.5, 2.25, 3.25}, {0, 1, 1.75, 2.75}, {0, 0.75, 1.25, 2, 3}, {0, 1, 2.25, 3}},
	6:  {{0, 1.25, 2, 2.75, 3.5}, {0, 0.75, 1.75, 2.5, 3.25}, {0, 1, 1.5, 2.25, 3}, {0, 1.5, 2.25, 3}},
	9:  {{0, 0.5, 1.25, 2, 2.75, 3.5}, {0, 0.75, 1.25, 1.75, 2.5, 3.25}, {0, 0.5, 1, 1.75, 2.25, 3}, {0, 0.75, 1.5, 2.25, 3}},
	12: {{0, 1.5, 2.5}, {0, 1.25, 3}, {0, 1.5, 2.25, 3.25}, {0, 1.5, 3}},
}

var pianoQuestionContour = [16]int{0, 1, 3, 2, 1, 2, 0, 1, 2, 4, 5, 3, 2, 1, 0, 1}
var pianoVariationContour = [16]int{0, 2, 3, 1, 2, 4, 3, 1, 3, 2, 5, 4, 3, 1, 0, 1}
var pianoAnswerContour = [16]int{3, 4, 2, 5, 4, 2, 1, 3, 5, 4, 2, 3, 2, 1, 0, 1}

type pianoContourDevelopment struct {
	invert      bool
	retrograde  bool
	barSequence [4]int
}

// pianoContourDevelopments transforms the existing motives at major formal returns instead of introducing unrelated note patterns.
var pianoContourDevelopments = map[int]pianoContourDevelopment{
	5:  {invert: true, barSequence: [4]int{0, 0, 1, 0}},
	6:  {retrograde: true, barSequence: [4]int{0, 1, 0, 0}},
	9:  {barSequence: [4]int{0, 1, 1, 2}},
	12: {invert: true, barSequence: [4]int{0, -1, -1, 0}},
}

var pianoMelodyCosts = struct {
	target   float64
	motion   float64
	leap     float64
	repeated float64
	passing  float64
}{target: 1.7, motion: 0.42, leap: 2.2, repeated: 14, passing: 1.0}
