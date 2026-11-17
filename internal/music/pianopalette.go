package music

// choosePianoPalette takes a scene character and image digest, then returns a copy with the seed's scale and chord substitutions.
func choosePianoPalette(character musicalCharacter, digest [32]byte) musicalCharacter {
	choice := pianoPaletteChoices[character.mood][digest[3]%3]
	if choice.mode == "" {
		return character
	}
	character.mode = choice.mode
	character.scale = choice.scale
	for _, change := range choice.changes {
		switch change.phrase {
		case 0:
			character.progression[change.position] = change.chord
		case 1:
			character.variation[change.position] = change.chord
		case 2:
			character.bridge[change.position] = change.chord
		case 3:
			character.bridgeReturn[change.position] = change.chord
		case 4:
			character.climax[change.position] = change.chord
		}
	}
	return character
}
