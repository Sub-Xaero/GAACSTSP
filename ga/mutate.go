package ga

import (
	"math"
)

// SwapMutate() picks two random indices along a bitstring, excluding the special "-1" element at the end, and swaps the
// values contained in those positions.
// Returns the resulting offspring Genome
func (genA *GeneticAlgorithm) SwapMutate(gene Genome) Genome {
	geneCopy := gene.Copy()
	newSequence := gene.Copy()
	min := 0
	max := len(geneCopy.Sequence) - 1

	// Choose 2 indexes
	choice1 := genA.RandomEngine.Intn(max-min) + min
	choice2 := genA.RandomEngine.Intn(max-min) + min

	// Swap elements
	newSequence.Sequence[choice1], newSequence.Sequence[choice2] = geneCopy.Sequence[choice2], geneCopy.Sequence[choice1]
	return newSequence.Copy()
}

// InversionMutate() picks two random indices along a bitstring, excluding the special "-1" element at the end, and reverses the
// order of the values between those positions.
// Returns the resulting offspring Genome
func (genA *GeneticAlgorithm) InversionMutate(gene Genome) Genome {
	newSequence := gene.Copy()
	min := 0
	max := len(newSequence.Sequence) - 1

	// Choose 2 indexes
	choice1 := genA.RandomEngine.Intn(max-min) + min
	choice2 := genA.RandomEngine.Intn(max-min) + min

	start := int(math.Min(float64(choice1), float64(choice2)))
	end := int(math.Max(float64(choice1), float64(choice2)))
	section := newSequence.Sequence[start:end]
	for i := len(section)/2 - 1; i >= 0; i-- {
		opp := len(section) - 1 - i
		section[i], section[opp] = section[opp], section[i]
	}
	return newSequence.Copy()
}
