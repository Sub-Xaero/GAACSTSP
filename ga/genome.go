package ga

import "fmt"

// Wrapper type for Bitstring
type Genome struct {
	Sequence Bitstring
}

// Copy the size of Genome, not the contents
func (gene Genome) CopySize() Genome {
	sequence := make(Bitstring, len(gene.Sequence))
	return Genome{sequence}
}

// Copy the contents of Genome
func (gene Genome) Copy() Genome {
	sequence := gene.Sequence.CopySize()
	copy(sequence, gene.Sequence)
	return Genome{sequence}
}

// Override of standard String() method, for custom output format when printing Genome structures
func (gene Genome) String() string {
	return fmt.Sprintf("{%v}", gene.Sequence)
}
