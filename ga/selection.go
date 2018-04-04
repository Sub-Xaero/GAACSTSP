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

func (genA *GeneticAlgorithm) RouletteChoice(candidatePool Population, cities map[string]City) int {
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
	return choice
}

func (genA *GeneticAlgorithm) RouletteSelection(candidatePool Population, cities map[string]City) Population {
	offspring := make(Population, 0)

	// For as many children as we want
	for i := 0; i < len(candidatePool); i++ {
		choice := genA.RouletteChoice(candidatePool, cities)
		offspring = append(offspring, candidatePool[choice].Copy())
	}
	return offspring
}

// Permutation based problem, requires this function to filter down the candidate pool so that the returned population
// contains only elements whose next chromosome to be picked is valid. e.g. with a a partial solution
// of [ 0 1 ], only return elements whose N position such that elem = [ x x N x ... ], does not contain either 0 or 1
func (genA *GeneticAlgorithm) ACOFilter(genePartial Bitstring, candidatePool Population, cities map[string]City) Population {
	candidatePoolCopy := candidatePool.Copy()
	filteredOffspring := make(Population, 0)

	numAlreadyVisited := len(genePartial)

	// Loop over candidate pool
	for i := 0; i < len(candidatePoolCopy); i++ {

		// Only filter on elements in positions > those already picked as visited
		// Slice containing the next possible element to be picked
		remainingSlice := candidatePoolCopy[i].Sequence[numAlreadyVisited : numAlreadyVisited+1]

		// Flag, if set to false in the below detection loop, elem is to be filtered out
		allowed := true

	toOuterLoop:
		for _, remainingCity := range remainingSlice {
			for _, visitedCity := range genePartial {
				// Permutation based problem, so remove elements that would duplicate tours of a city if picked
				if visitedCity == remainingCity {
					allowed = false
					break toOuterLoop
				}
			}
		}
		if allowed {
			filteredOffspring = append(filteredOffspring, candidatePoolCopy[i])
		}
	}

	return filteredOffspring
}

// SetSelectionFunc changes the selection function to the function specified
func (genA *GeneticAlgorithm) Selection(candidatePool Population, cities map[string]City) Population {
	return genA.TournamentSelection(candidatePool, cities)
}
