package music

import "math"

// noiseGradient takes an image digest and integer lattice point, then returns a stable signed gradient.
func noiseGradient(digest [32]byte, lattice int) float64 {
	hash := uint32(lattice)*0x9e3779b9 ^ uint32(digest[lattice&31])<<16 ^ uint32(digest[(lattice+11)&31])
	hash ^= hash >> 16
	hash *= 0x7feb352d
	hash ^= hash >> 15
	return float64(hash&0xffff)/32767.5 - 1
}

// smoothGradientNoise takes a continuous position and image digest, then returns deterministic one-dimensional gradient noise.
func smoothGradientNoise(position float64, digest [32]byte) float64 {
	lattice := int(math.Floor(position))
	fraction := position - float64(lattice)
	fade := fraction * fraction * fraction * (fraction*(fraction*6-15) + 10)
	left := noiseGradient(digest, lattice) * fraction
	right := noiseGradient(digest, lattice+1) * (fraction - 1)
	return left*(1-fade) + right*fade
}

// phraseEnergy takes a bar's planned energy and image digest, then returns a gently evolving dynamic level across neighboring bars.
func phraseEnergy(base float64, bar int, digest [32]byte) float64 {
	position := float64(bar)*0.18 + float64(digest[6])*0.019
	return base * (1 + 0.3*smoothGradientNoise(position, digest))
}
