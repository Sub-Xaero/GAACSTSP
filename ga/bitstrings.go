package ga

import (
	"fmt"
)

type Bitstring []string

// GenerateCandidate returns an encoded string as set by calls SetGenerateBitString. Defaults to binary strings
func (genA *GeneticAlgorithm) GenerateCandidate(length int, cities map[string]City) (Bitstring, error) {
	keys := make([]string, 0)

	for key := range cities {
		keys = append(keys, key)
	}

	sequence := make(Bitstring, 0)
	for i := 0; i < length; i++ {
		choice := genA.RandomEngine.Int() % len(keys)
		sequence = append(sequence, keys[choice])
		keys = append(keys[:choice], keys[choice+1:]...)
	}

	sequence = append(sequence, "-1")
	return sequence, nil
}

// SetGenerateBitString sets the function that generates the Bitstring candidatePool
func (b Bitstring) String() string {
	output := ""
	for _, val := range b {
		output += fmt.Sprintf("%v", val) + " "
	}
	return "[" + output + "]"
}
