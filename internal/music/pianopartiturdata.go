package music

type pianoScoreStep struct {
	beat       float64
	tone       int
	alteration int
	direction  int
}

type pianoPhraseKind int

const (
	pianoTheme pianoPhraseKind = iota
	pianoVariation
	pianoBridge
	pianoTransition
	pianoDevelopment
	pianoBridgeReturn
	pianoClimax
	pianoEcho
	pianoCoda
)

// Chord tones are root, third, and fifth; a negative alteration requests the lower modal neighbor before resolution.
var pianoThemeFigure = [4][]pianoScoreStep{
	{{0, 2, 0, 0}, {0.5, 2, -1, -1}, {1, 2, 0, 1}, {1.5, 1, 0, -1}, {2.5, 0, 0, -1}, {3.25, 1, 0, 1}},
	{{0, 1, 0, 0}, {0.75, 2, 0, 1}, {1.5, 1, 0, -1}, {2.25, 0, 0, -1}, {3.25, 2, 0, 1}},
	{{0, 2, 0, 0}, {0.75, 1, 0, -1}, {1.5, 0, 0, -1}, {2.25, 1, 0, 1}, {3.25, 2, 0, 1}},
	{{0, 2, 0, 0}, {1, 1, 0, -1}, {2, 2, 0, 1}, {3, 2, 0, 0}},
}

var pianoVariationFigure = [4][]pianoScoreStep{
	{{0, 2, 0, 0}, {0.5, 1, 0, -1}, {1, 2, -1, 1}, {1.5, 2, 0, 1}, {2.5, 1, 0, -1}, {3.25, 0, 0, -1}},
	{{0, 1, 0, 0}, {0.75, 2, 0, 1}, {1.5, 1, 0, -1}, {2.5, 0, 0, -1}, {3.25, 1, 0, 1}},
	{{0, 2, 0, 0}, {0.5, 1, 0, -1}, {1.25, 0, 0, -1}, {2, 1, 0, 1}, {3, 2, 0, 1}},
	{{0, 1, 0, 0}, {1, 2, 0, 1}, {2, 1, 0, -1}, {3, 0, 0, -1}},
}

// pianoOpeningIdeas are three scored thematic identities, each with a related but independent variation.
var pianoOpeningIdeas = [3][2][2][]pianoScoreStep{
	{
		{pianoThemeFigure[0], pianoThemeFigure[1]},
		{pianoVariationFigure[0], pianoVariationFigure[1]},
	},
	{
		{
			{{0, 1, 0, 0}, {0.5, 2, 0, 1}, {1, 1, 0, -1}, {1.75, 0, 0, -1}, {2.5, 1, 0, 1}, {3.25, 2, 0, 1}},
			{{0, 2, 0, 0}, {0.75, 1, 0, -1}, {1.5, 2, 0, 1}, {2.25, 0, 0, -1}, {3.25, 1, 0, 1}},
		},
		{
			{{0, 1, 0, 0}, {0.5, 0, 0, -1}, {1, 2, 0, 1}, {1.75, 1, 0, -1}, {2.5, 0, 0, -1}, {3.25, 1, 0, 1}},
			{{0, 2, 0, 0}, {0.75, 1, 0, -1}, {1.5, 2, 0, 1}, {2.5, 1, 0, -1}, {3.25, 0, 0, -1}},
		},
	},
	{
		{
			{{0, 0, 0, 0}, {0.75, 1, 0, 1}, {1.5, 2, 0, 1}, {2.25, 1, 0, -1}, {3.25, 0, 0, -1}},
			{{0, 1, 0, 0}, {0.5, 2, 0, 1}, {1.25, 1, 0, -1}, {2, 0, 0, -1}, {2.75, 1, 0, 1}, {3.5, 2, 0, 1}},
		},
		{
			{{0, 2, 0, 0}, {0.75, 1, 0, -1}, {1.5, 2, 0, 1}, {2.25, 1, 0, -1}, {3.25, 0, 0, -1}},
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
	{{0, 1, 0, 0}, {0.75, 2, 0, 1}, {1.5, 1, 0, -1}, {2.5, 0, 0, -1}, {3.25, 1, 0, 1}},
	{{0, 2, 0, 0}, {1, 1, 0, -1}, {1.75, 0, 0, -1}, {2.5, 1, 0, 1}, {3.25, 2, 0, 1}},
	{{0, 0, 0, 0}, {0.75, 1, 0, 1}, {1.5, 2, 0, 1}, {2.25, 1, 0, -1}, {3.25, 0, 0, -1}},
	{{0, 2, 0, 0}, {0.75, 1, 0, -1}, {1.75, 0, 0, -1}, {2.5, 1, 0, 1}, {3.25, 2, 0, 1}},
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
	pianoTheme:        pianoThemeFigure,
	pianoVariation:    pianoVariationFigure,
	pianoBridge:       pianoBridgeFigure,
	pianoTransition:   pianoTransitionFigure,
	pianoDevelopment:  pianoDevelopmentFigure,
	pianoBridgeReturn: pianoBridgeReturnFigure,
	pianoClimax:       pianoClimaxFigure,
	pianoEcho:         pianoEchoFigure,
	pianoCoda:         pianoCodaFigure,
}

// pianoEpisodeFigures places an opening theme, variations, contrasting material, developed returns, and a cadential ending.
var pianoEpisodeFigures = [16]pianoPhraseKind{
	pianoTheme, pianoVariation, pianoBridge, pianoTransition,
	pianoTheme, pianoDevelopment, pianoBridgeReturn, pianoTransition,
	pianoTheme, pianoClimax, pianoTransition, pianoTheme,
	pianoEcho, pianoTransition, pianoCoda, pianoCoda,
}
