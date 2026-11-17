package music

// pianoAppliedApproaches marks phrase positions where an applied major dominant may prepare the next chord.
var pianoAppliedApproaches = map[int][4]bool{
	0: {true, false, true, false},
	1: {false, true, false, false},
	2: {false, true, true, false},
	5: {true, false, true, false},
	6: {false, true, true, false},
	9: {true, false, true, false},
}
