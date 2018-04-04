package ga

import "fmt"

// Default output function, prints whatever is provided to it to the standard TTY console.
func PrintToConsole (a ...interface{}) {
	fmt.Println(a...)
}

// Allows customisation of the way in which the algorithm outputs logging information
//
// i.e. Log to file
func (genA *GeneticAlgorithm) SetOutputFunc(f func(a ...interface{})) {
	genA.Output = f
}
