package ga

func (genA *GeneticAlgorithm) TournamentSelection(candidatePool Population, cities map[string]City) Population {
	offspring := make(Population, 0)

	for i := 0; i < len(candidatePool); i++ {
		parent1 := candidatePool[genA.RandomEngine.Intn(len(candidatePool))]
		parent2 := candidatePool[genA.RandomEngine.Intn(len(candidatePool))]

		if genA.Fitness(parent1, cities) > genA.Fitness(parent2, cities) {
			offspring = append(offspring, parent1)
		} else {
			offspring = append(offspring, parent2)
		}
	}

	return offspring
}

func (genA *GeneticAlgorithm) RouletteSelection(candidatePool Population, cities map[string]City) Population {
	offspring := make(Population, 0)
	for range candidatePool {
		weightSum := 0.0
		for _, val := range candidatePool {
			weightSum += genA.Fitness(val, cities)
		}
		choice := genA.RandomEngine.Float32() * float32(weightSum)
		for _, val := range candidatePool {
			choice -= float32(genA.Fitness(val, cities))
			if choice <= 0 {
				offspring = append(offspring, val.Copy())
				break
			}
		}
	}
	return offspring
}

// SetSelectionFunc changes the selection function to the function specified
func (genA *GeneticAlgorithm) Selection(candidatePool Population, cities map[string]City) Population {
	return genA.TournamentSelection(candidatePool, cities)
}
