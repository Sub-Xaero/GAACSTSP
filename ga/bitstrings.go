package ga

import (
	"fmt"
	"errors"
)

// Defines the storage data type that chromosomes will be encoded as
type Bitstring []string

// GenerateCandidate returns an encoded string containing a permutation of the indices provided by the argument 'cities'
func (genA *GeneticAlgorithm) GenerateCandidate(length int, cities map[string]City) (Bitstring, error) {
	// +1 to account for return to first city
	reqLength := len(cities) + 1

	// Error conditions
	if length < reqLength {
		return nil, errors.New("string is not long enough to encode tour")
	}
	if length > reqLength {
		return nil, errors.New("string is longer than a full tour")
	}

	// Cities are a map, keys are the indices, so use those as the city identifier
	keys := make([]string, 0)
	for key := range cities {
		keys = append(keys, key)
	}

	// Construct bitstring
	sequence := make(Bitstring, 0)
	for i := 0; i < length-1; i++ {
		// Choose a random key
		choice := genA.RandomEngine.Intn(len(keys))
		// Add it the the sequence
		sequence = append(sequence, keys[choice])
		// Remove from possibilities so it does not occur twice
		keys = append(keys[:choice], keys[choice+1:]...)
	}

	// -1 means a return to the first city, add to the end of every candidate
	sequence = append(sequence, "-1")
	return sequence, nil
}

// Copy only the size of the Bitstring, not the contents
func (b Bitstring) CopySize() Bitstring {
	sequence := make(Bitstring, len(b))
	return sequence
}

// Copy the contents of a Bitstring
func (b Bitstring) Copy() Bitstring {
	sequence := b.CopySize()
	copy(sequence, b)
	return sequence
}

// Override of the standard String() method to be able to output the contents of a bitstring in a more custom and readable format
func (b Bitstring) String() string {
	output := ""
	for _, val := range b {
		output += fmt.Sprintf("%v", val) + " "
	}
	return "[" + output + "]"
}
