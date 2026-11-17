package music

type pianoScoreStep struct {
	beat       float64
	tone       int
	alteration int
	direction  int
}

type pianoPhraseKind int

const (
	pianoQuestion pianoPhraseKind = iota
	pianoConsequence
	pianoBridge
	pianoTransition
	pianoDevelopment
	pianoBridgeReturn
	pianoClimax
	pianoEcho
	pianoCoda
)

// Chord tones are root, third, and fifth; alteration introduces a neighbor that resolves at the next scored attack.
var pianoQuestionFigure = [4][]pianoScoreStep{
	{{0, 2, 0, 0}, {0.5, 2, -1, -1}, {1, 2, 0, 1}, {1.5, 1, 0, -1}, {2.5, 0, 0, -1}, {3.25, 1, 0, 1}},
	{{0, 1, 0, 0}, {0.75, 2, 0, 1}, {1.5, 1, 0, -1}, {2.25, 0, 0, -1}, {3.25, 2, 0, 1}},
	{{0, 2, 0, 0}, {0.75, 1, 0, -1}, {1.5, 0, 0, -1}, {2.25, 1, 0, 1}, {3.25, 2, 0, 1}},
	{{0, 2, 0, 0}, {1, 1, 0, -1}, {2, 2, 0, 1}, {3, 1, 0, -1}},
}

var pianoConsequenceFigure = [4][]pianoScoreStep{
	{{0, 0, 0, 0}, {0.5, 2, 0, 1}, {1, 2, -1, -1}, {1.5, 2, 0, 1}, {2.5, 1, 0, -1}, {3.25, 0, 0, -1}},
	{{0, 1, 0, 0}, {0.75, 2, 0, 1}, {1.5, 1, 0, -1}, {2.5, 0, 0, -1}, {3.25, 1, 0, 1}},
	{{0, 2, 0, 0}, {0.5, 1, 0, -1}, {1.25, 0, 0, -1}, {2, 1, 0, 1}, {3, 2, 0, 1}},
	{{0, 1, 0, 0}, {1, 2, 0, 1}, {2, 1, 0, -1}, {3, 0, 0, -1}},
}

// pianoOpeningIdeas are three scored thematic identities; the consequent answers each with a tonic arrival.
var pianoOpeningIdeas = [3][2][2][]pianoScoreStep{
	{
		{pianoQuestionFigure[0], pianoQuestionFigure[1]},
		{pianoConsequenceFigure[0], pianoConsequenceFigure[1]},
	},
	{
		{
			{{0, 1, 0, 0}, {0.5, 2, 0, 1}, {1, 1, 0, -1}, {1.75, 0, 0, -1}, {2.5, 1, 0, 1}, {3.25, 2, 0, 1}},
			{{0, 2, 0, 0}, {0.75, 1, 0, -1}, {1.5, 2, 0, 1}, {2.25, 0, 0, -1}, {3.25, 1, 0, 1}},
		},
		{
			{{0, 0, 0, 0}, {0.5, 1, 0, 1}, {1, 2, 0, 1}, {1.75, 1, 0, -1}, {2.5, 0, 0, -1}, {3.25, 1, 0, 1}},
			{{0, 2, 0, 0}, {0.75, 1, 0, -1}, {1.5, 2, 0, 1}, {2.5, 1, 0, -1}, {3.25, 0, 0, -1}},
		},
	},
	{
		{
			{{0, 0, 0, 0}, {0.75, 1, 0, 1}, {1.5, 2, 0, 1}, {2.25, 1, 0, -1}, {3.25, 0, 0, -1}},
			{{0, 1, 0, 0}, {0.5, 2, 0, 1}, {1.25, 1, 0, -1}, {2, 0, 0, -1}, {2.75, 1, 0, 1}, {3.5, 2, 0, 1}},
		},
		{
			{{0, 0, 0, 0}, {0.75, 2, 0, 1}, {1.5, 1, 0, -1}, {2.25, 2, 0, 1}, {3.25, 0, 0, -1}},
			{{0, 1, 0, 0}, {0.5, 2, 0, 1}, {1.25, 1, 0, -1}, {2, 0, 0, -1}, {2.75, 1, 0, 1}, {3.5, 2, 0, 1}},
		},
	},
}

var pianoBridgeFigure = [4][]pianoScoreStep{
	{{0, 1, 0, 0}, {0.5, 0, 0, -1}, {1, 1, 0, 1}, {1.75, 2, 0, 1}, {2.5, 1, 0, -1}, {3.25, 0, 0, -1}},
	{{0, 2, 0, 0}, {0.75, 1, 0, -1}, {1.5, 2, 0, 1}, {2.25, 0, 0, -1}, {3, 1, 0, 1}},
	{{0, 0, 0, 0}, {0.5, 1, 0, 1}, {1, 2, 0, 1}, {1.5, 1, 0, -1}, {2.5, 0, 0, -1}, {3.25, 1, 0, 1}},
	{{0, 2, 0, 0}, {1, 1, 0, -1}, {2, 2, 0, 1}, {3, 1, 0, -1}},
}

var pianoTransitionFigure = [4][]pianoScoreStep{
	{{0, 1, 0, 0}, {1.25, 2, 0, 1}, {2.5, 1, 0, -1}, {3.25, 0, 0, -1}},
	{{0, 2, 0, 0}, {1.5, 1, 0, -1}, {2.75, 0, 0, -1}},
	{{0, 0, 0, 0}, {1.25, 1, 0, 1}, {2.25, 2, 0, 1}, {3.25, 1, 0, -1}},
	{{0, 2, 0, 0}, {1.5, 1, 0, -1}, {2.75, 0, 0, -1}},
}

var pianoDevelopmentFigure = [4][]pianoScoreStep{
	{{0, 2, 0, 0}, {0.75, 1, 0, -1}, {1.5, 0, 0, -1}, {2, 1, 0, 1}, {3, 2, 0, 1}},
	{{0, 2, 0, 0}, {0.5, 2, -1, -1}, {1, 2, 0, 1}, {1.75, 1, 0, -1}, {2.5, 0, 0, -1}, {3.25, 1, 0, 1}},
	{{0, 0, 0, 0}, {0.75, 1, 0, 1}, {1.5, 2, 0, 1}, {2.25, 1, 0, -1}, {3, 0, 0, -1}},
	{{0, 2, 0, 0}, {1, 1, 0, -1}, {2, 2, 0, 1}, {3, 0, 0, -1}},
}

var pianoBridgeReturnFigure = [4][]pianoScoreStep{
	{{0, 2, 0, 0}, {1.25, 1, 0, -1}, {2, 0, 0, -1}, {2.75, 1, 0, 1}, {3.5, 2, 0, 1}},
	{{0, 1, 0, 0}, {0.75, 2, 0, 1}, {1.75, 1, 0, -1}, {2.5, 0, 0, -1}, {3.25, 1, 0, 1}},
	{{0, 0, 0, 0}, {0.5, 1, 0, 1}, {1.5, 2, 0, 1}, {2.25, 1, 0, -1}, {3, 0, 0, -1}},
	{{0, 2, 0, 0}, {1.5, 1, 0, -1}, {2.25, 2, 0, 1}, {3, 1, 0, -1}},
}

var pianoClimaxFigure = [4][]pianoScoreStep{
	{{0, 0, 0, 0}, {0.5, 1, 0, 1}, {1, 2, 0, 1}, {1.5, 1, 0, -1}, {2, 0, 0, -1}, {2.75, 1, 0, 1}, {3.5, 2, 0, 1}},
	{{0, 2, 0, 0}, {0.5, 2, -1, -1}, {1, 2, 0, 1}, {1.5, 1, 0, -1}, {2, 0, 0, -1}, {2.75, 1, 0, 1}, {3.5, 2, 0, 1}},
	{{0, 0, 0, 0}, {0.5, 1, 0, 1}, {1, 2, 0, 1}, {1.5, 1, 0, -1}, {2, 0, 0, -1}, {2.75, 1, 0, 1}, {3.5, 2, 0, 1}},
	{{0, 2, 0, 0}, {0.75, 1, 0, -1}, {1.5, 2, 0, 1}, {2.25, 0, 0, -1}, {3, 1, 0, 1}},
}

var pianoEchoFigure = [4][]pianoScoreStep{
	{{0, 2, 0, 0}, {1.5, 1, 0, -1}, {2.75, 0, 0, -1}},
	{{0, 1, 0, 0}, {1.25, 2, 0, 1}, {2.75, 1, 0, -1}},
	{{0, 2, 0, 0}, {1, 1, 0, -1}, {2.25, 0, 0, -1}, {3.25, 1, 0, 1}},
	{{0, 2, 0, 0}, {1.5, 1, 0, -1}, {3, 1, 0, 0}},
}

var pianoCodaFigure = [4][]pianoScoreStep{
	{{0, 1, 0, 0}, {1.5, 2, 0, 1}, {3, 1, 0, -1}},
	{{0, 2, 0, 0}, {1.5, 1, 0, -1}, {3, 0, 0, -1}},
	{{0, 1, 0, 0}, {1.5, 2, 0, 1}, {3, 1, 0, -1}},
	{{0, 0, 0, 0}},
}

// pianoPhraseFigures assigns authored four-bar ideas to formal roles; harmony and register are resolved when the score is written.
var pianoPhraseFigures = map[pianoPhraseKind][4][]pianoScoreStep{
	pianoQuestion:     pianoQuestionFigure,
	pianoConsequence:  pianoConsequenceFigure,
	pianoBridge:       pianoBridgeFigure,
	pianoTransition:   pianoTransitionFigure,
	pianoDevelopment:  pianoDevelopmentFigure,
	pianoBridgeReturn: pianoBridgeReturnFigure,
	pianoClimax:       pianoClimaxFigure,
	pianoEcho:         pianoEchoFigure,
	pianoCoda:         pianoCodaFigure,
}

// pianoEpisodeFigures forms an antecedent and consequent, a contrasting middle, developed returns, and a cadential ending.
var pianoEpisodeFigures = [16]pianoPhraseKind{
	pianoQuestion, pianoConsequence, pianoBridge, pianoTransition,
	pianoQuestion, pianoDevelopment, pianoBridgeReturn, pianoTransition,
	pianoQuestion, pianoClimax, pianoTransition, pianoQuestion,
	pianoEcho, pianoTransition, pianoCoda, pianoCoda,
}
