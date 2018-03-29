package ga

// Genome represents a bitstring and associated fitness value
type Genome struct {
	Sequence Bitstring
}

type Population [] Genome
