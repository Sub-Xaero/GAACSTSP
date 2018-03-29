package ga

import (
	"math/rand"
	"time"
)
type GeneticAlgorithm struct {
	Candidates    Population
	BestCandidate Genome
	Generations   int

	IterationsSinceChange int

	Output       func(a ...interface{})
	RandomEngine *rand.Rand
}

func NewGeneticAlgorithm() GeneticAlgorithm {
	var geneticAlgorithm GeneticAlgorithm
	geneticAlgorithm.SetOutputFunc(PrintToConsole)
	geneticAlgorithm.SetSeed(time.Now().Unix())

	return geneticAlgorithm
}
