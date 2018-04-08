package ga

import (
	"math"
)

// Given a pool of candidates, pick 2 parents at random, compare their fitness, and add the parent with the best fitness
// to the pool of offspring. Repeat N times to fill the offspring pool, such that
// 	N = len(candidatePool)
func (genA *GeneticAlgorithm) TournamentSelection(candidatePool Population, cities map[string]City) Population {
	offspring := make(Population, 0)

	for i := 0; i < len(candidatePool); i++ {
		parent1 := candidatePool[genA.RandomEngine.Intn(len(candidatePool))]
		parent2 := candidatePool[genA.RandomEngine.Intn(len(candidatePool))]

		// Important that value received is for the correct candidate, so single buffered channel for each
		fitness1 := make(chan float64, 1)
		fitness2 := make(chan float64, 1)

		// Fetch fitness of both values concurrently, negligible for small bitstrings, but helpful for large ones
		go func() {
			fitness1 <- genA.Fitness(parent1, cities)
		}()

		go func() {
			fitness2 <- genA.Fitness(parent2, cities)
		}()

		if <-fitness1 > <-fitness2 {
			offspring = append(offspring, parent1)
		} else {
			offspring = append(offspring, parent2)
		}
	}

	return offspring
}

// Multi-Threaded.
//
// Given a pool of candidates, sum the fitness of the candidates.
// Multiply the sum fitness by a random float, to pick a random point in the pool. Then iteratively
// reduce the sum by the fitness of each candidate in the pool in order, and the candidate that reduces the
// sum below zero is the chosen candidate.
//
// Higher fitnesses reduce the sum by more, therefore are more likely to be picked, giving this selection method its
// 'weighted', trait.
//
// Returns the index of the candidate in the pool that was selected
func (genA *GeneticAlgorithm) RouletteChoice(candidatePool Population, cities map[string]City) int {
	max := math.Abs(genA.MinFitness(candidatePool, cities)) + 1
	// Build weights
	weightSum := 0.0
	// Multi-thread to speed up
	ch := make(chan float64, len(candidatePool))
	for _, val := range candidatePool {
		go func() {
			val := max - math.Abs(genA.Fitness(val, cities))
			ch <- val
		}()
	}
	for range candidatePool {
		weightSum += <-ch
	}

	choice := 0
	chosenPoint := genA.RandomEngine.Float64() * weightSum

	// Reduce weights to find position chosen
	for index := range candidatePool {
		reduce := max - math.Abs(genA.Fitness(candidatePool[index], cities))
		chosenPoint -= reduce
		if chosenPoint <= 0 {
			choice = index
			break
		}
	}
	return choice
}

// Wrapper around RouletteChoice(), given a pool of candidates, call RouletteChoice() N times
// such that
// 	N = len(candidatePool)
// and return the new set of offspring.
func (genA *GeneticAlgorithm) RouletteSelection(candidatePool Population, cities map[string]City) Population {
	offspring := make(Population, 0)
	poolSize := len(candidatePool)

	ch := make(chan Genome, poolSize)
	// For as many children as we want
	for i := 0; i < poolSize; i++ {
		go func() {
			choice := genA.RouletteChoice(candidatePool, cities)
			ch <- candidatePool[choice].Copy()
		}()
	}
	for i := 0; i < poolSize; i++ {
		offspring = append(offspring, <-ch)
	}
	return offspring
}

// This function filters down the candidate pool so that the returned population
// contains only elements whose next chromosome to be picked is valid.
//
// e.g. with a partial solution
// of
// 	[ 0 1 ]
// the next position to be picked is 3, therefore only return elements whose N position such that
// 	elem = [ x x N x ... ]
// does not contain either 0 or 1
func (genA *GeneticAlgorithm) ACOFilter(genePartial Bitstring, candidatePool Population, cities map[string]City) Population {
	candidatePoolCopy := candidatePool.Copy()
	poolSize := len(candidatePoolCopy)
	filteredOffspring := make(Population, 0)
	numAlreadyVisited := len(genePartial)

	if numAlreadyVisited == 0 {
		return candidatePoolCopy
	}
	// Loop over candidate pool
	for i := 0; i < poolSize; i++ {
		// Only filter on elements in positions > those already picked as visited
		// Slice containing the next possible element to be picked
		remainingSlice := candidatePoolCopy[i].Sequence[numAlreadyVisited:]

		// Channel, to synchronise concurrent goroutines over
		ch := make(chan string, len(remainingSlice))
		for _, remainingCity := range remainingSlice {
			// Launch goroutine, speed up processing
			go func(ch chan string, remainingCity string) {
				for _, visitedCity := range genePartial {
					// Permutation based problem, so remove elements that would duplicate tours of a city if picked
					if visitedCity == remainingCity {
						// This city is not allowed, send key string over channel
						ch <- "nil"
						return
					}
				}
				// City is allowed, send anything over channel just to fill buffered spaces
				ch <- remainingCity
			}(ch, remainingCity)
		}

		// Flag, if set to false in the below detection loop, elem is to be filtered out
		allowed := true
		for x := 0; x < len(remainingSlice); x++ {
			elem := <-ch
			if elem == "nil" {
				allowed = false
			}
		}

		if allowed {
			filteredOffspring = append(filteredOffspring, candidatePoolCopy[i])
		}
	}

	return filteredOffspring
}

func (genA *GeneticAlgorithm) ACOPolyfill(genePartial Bitstring, cities map[string]City) Bitstring {
	output := genePartial.Copy()
	numCities := len(cities)
	lenPartial := len(output)
	if lenPartial > numCities {
		panic("Trying to polyfill string past tour length")
	} else if lenPartial == numCities {
		output = append(output, "-1")
	} else {
		count := make(map[string]int)
		for i := 1; i <= numCities; i++ {
			count[strconv.Itoa(i)] = 0
		}

		for _, val := range output {
			count[val]++
		}

		for key, val := range count {
			if val > 1 {
				panic("Duplicate key, should not exist")
			}
			if val == 0 {
				output = append(output, key)
				break
			}
		}
	}
	return output
}

func (genA *GeneticAlgorithm) ACOSelection(candidatePool Population, cities map[string]City) Population {
	numChildren := len(candidatePool)
	offspring := make(Population, 0)
	possibilities := candidatePool.Copy()

	// Assumes all bitstrings same length
	strLen := len(candidatePool[0].Sequence)

	// Communication channel
	ch := make(chan Bitstring, numChildren)

	// For as many children as we want
	for i := 0; i < numChildren; i++ {
		// Launch concurrent goroutines, each one generate a child
		go func() {
			child := Bitstring{}
			for x := 0; x < strLen; x++ {
				filteredPossibilities := genA.ACOFilter(child, possibilities, cities)
				choice := genA.RouletteChoice(filteredPossibilities, cities)
				child = append(child, filteredPossibilities[choice].Sequence[x])
			}
			ch <- child
		}()
	}

	// Receive children as they are created async
	for i := 0; i < numChildren; i++ {
		offspring = append(offspring, Genome{<-ch})
	}

	return offspring
}