package ga

// Wrapper type, to make for easier manipulation and passing around of collections of Genomes
type Population [] Genome

// Copy size of Population, not contents
func (pop Population) CopySize() Population {
	copyPop := make(Population, len(pop))
	return copyPop
}

// Copy contents of Population
func (pop Population) Copy() Population {
	copyPop := pop.CopySize()
	copy(copyPop, pop)
	return copyPop
}
