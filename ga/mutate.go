package ga

func (genA *GeneticAlgorithm) Mutate(gene Genome) Genome {
	geneCopy := gene.Copy()
	newSequence := gene.Copy()
	min := 1
	max := len(geneCopy.Sequence) - 1

	// Choose 2 indexes
	choice1 := genA.RandomEngine.Intn(max-min) + min
	choice2 := genA.RandomEngine.Intn(max-min) + min

	// Swap elements
	newSequence.Sequence[choice1], newSequence.Sequence[choice2] = geneCopy.Sequence[choice2], geneCopy.Sequence[choice1]
	return newSequence.Copy()
}
