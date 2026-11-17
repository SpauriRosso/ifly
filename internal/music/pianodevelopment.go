package music

// developedPianoDegree takes an original contour, its sampled position, an episode, and a bar within it, then returns a transformed scale degree.
func developedPianoDegree(contour [16]int, contourIndex, episodeIndex, position int) int {
	development, ok := pianoContourDevelopments[episodeIndex]
	if !ok {
		return contour[contourIndex]
	}
	if development.retrograde {
		contourIndex = (contourIndex/4)*4 + 3 - contourIndex%4
	}
	degree := contour[contourIndex]
	if development.invert {
		degree = 4 - degree
	}
	return degree + development.barSequence[position]
}
