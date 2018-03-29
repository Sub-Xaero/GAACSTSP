package ga

// Genome represents a bitstring and associated fitness value
type Genome struct {
	Sequence Bitstring
}

type Population [] Genome

func (gene Genome) Copy() Genome {
	sequence := make(Bitstring, len(gene.Sequence))
	copy(sequence, gene.Sequence)
	return Genome{sequence}
}

