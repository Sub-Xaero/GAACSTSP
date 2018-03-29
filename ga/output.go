package ga

import "fmt"

var PrintToConsole = func(a ...interface{}) {
	fmt.Println(a...)
}

func (genA *GeneticAlgorithm) SetOutputFunc(f func(a ...interface{})) {
	genA.Output = f
}
