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
	genomeSequence := gene.Copy().Sequence

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
