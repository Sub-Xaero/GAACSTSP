package ga

import (
	"fmt"
)

type Bitstring []string

// SetGenerateBitString sets the function that generates the Bitstring candidatePool
func (b Bitstring) String() string {
	output := ""
	for _, val := range b {
		output += fmt.Sprintf("%v", val) + " "
	}
	return "[" + output + "]"
}
