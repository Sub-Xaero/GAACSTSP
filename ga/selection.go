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

// SetSelectionFunc changes the selection function to the function specified
func (genA *GeneticAlgorithm) Selection(candidatePool Population, cities map[string]City) Population {
	return genA.TournamentSelection(candidatePool, cities)
}
