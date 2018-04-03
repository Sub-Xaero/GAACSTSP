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

	// For as many children as we want
	for i := 0; i < len(candidatePool); i++ {
		// Build weights
		weightSum := 0.0
		// Multi-thread to speed up
		ch := make(chan float64, len(candidatePool))
		for _, val := range candidatePool {
			go func() {
				ch <- genA.Fitness(val, cities)
			}()
		}
		for range candidatePool {
			weightSum += <-ch
		}

		choice := 0
		chosenPoint := genA.RandomEngine.Float64() * weightSum

		// Reduce weights to find position chosen
		for index := range candidatePool {
			chosenPoint -= genA.Fitness(candidatePool[index], cities)
			if chosenPoint <= 0 {
				choice = index
				break
			}
		}

		offspring = append(offspring, candidatePool[choice].Copy())
	}
	return offspring
}

// SetSelectionFunc changes the selection function to the function specified
func (genA *GeneticAlgorithm) Selection(candidatePool Population, cities map[string]City) Population {
	return genA.TournamentSelection(candidatePool, cities)
}
