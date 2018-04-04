package ga

import (
	"math"
)

func EuclideanDistance(c1, c2 City) float64 {
	var difX = math.Pow(c1.x-c2.x, 2)
	var difY = math.Pow(c1.y-c2.y, 2)
	var sqSum = math.Sqrt(difX + difY)
	return sqSum
}

func (genA *GeneticAlgorithm) Fitness(gene Genome, cities map[string]City) float64 {
	genomeSequence := gene.Sequence.Copy()

	totalDistance := 0.0
	for i := 0; i < len(genomeSequence)-2; i++ {
		city1, city2 := genomeSequence[i], genomeSequence[i+1]
		if city2 == "-1" {
			city2 = genomeSequence[0]
		}
		totalDistance += EuclideanDistance(cities[city1], cities[city2])
	}
	return 0.0 - totalDistance
}

// AverageFitness returns the average fitness of a [] Genome candidatePool
func (genA *GeneticAlgorithm) AverageFitness(candidatePool Population, cities map[string]City) float64 {
	var average = 0.0
	for _, i := range candidatePool {
		average += genA.Fitness(i, cities)
	}
	return average / float64(len(candidatePool))
}

// MaxFitness returns the highest fitness found in a [] Genome candidatePool
func (genA *GeneticAlgorithm) MaxFitnessCandidate(candidatePool Population, cities map[string]City) Genome {
	var (
		set     = false
		max     = 0.0
		maxGene Genome
	)
	for _, i := range candidatePool {
		fitness := genA.Fitness(i, cities)
		if !set {
			set = true
			max = fitness
			maxGene = i.Copy()
		}
		if fitness > max {
			max = fitness
			maxGene = i.Copy()
		}
	}
	return maxGene
}

// MaxFitness returns the highest fitness found in a [] Genome candidatePool
func (genA *GeneticAlgorithm) MaxFitness(candidatePool Population, cities map[string]City) float64 {
	return genA.Fitness(genA.MaxFitnessCandidate(candidatePool, cities), cities)
}
