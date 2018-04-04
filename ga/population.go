package ga

type Population [] Genome

func (pop Population) CopySize() Population {
	copyPop := make(Population, len(pop))
	return copyPop
}

func (pop Population) Copy() Population {
	copyPop := pop.CopySize()
	copy(copyPop, pop)
	return copyPop
}
