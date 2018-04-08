package ga

import (
	"math"
)

// Function to actually perform the mapping operation described in SinglePointPMX()
func pmxInternal(gene Genome, bitstring Bitstring) Genome {
	// Copy, so as not to corrupt existing gene
	geneCopy := gene.Copy()

	// The bitstring to contain the end result, initialised to empty size of string
	output := gene.Sequence.CopySize()

	// The numbers to find the positions of and swap with
	swapNumbers := bitstring.Copy()

	// Map positions of cities in the string
	positions := make(map[string]int)
	for index, city := range geneCopy.Sequence {
		positions[city] = index
	}

	sectionEnd := len(swapNumbers)
	section := geneCopy.Sequence[:sectionEnd]
	for index := range section {
		// Store city numbers to swap
		swapFrom := section[index]
		swapWith := swapNumbers[index]
		// Store index of each city, ready to swap
		posFrom, posWith := positions[swapFrom], positions[swapWith]
		// Remap the positions of the numbers
		positions[swapFrom] = posWith
		positions[swapWith] = posFrom
	}

	// Reconstruct string based on map of changed positions
	for city, index := range positions {
		output[index] = city
	}

	// Construct new gene and return
	geneOutput := Genome{output}
	return geneOutput
}

// Permutation based problems cannot have multiple occurrences of a tour, crossover would
// break constraints of problem, Partially-Mapped Crossover (PMX) instead acts as a mutation operator of sorts.
//
// The crossover section, instead of swapping with a breeding parents respective crossover section, acts as a directive -
// identifying indices to swap with within the current parent, as opposed to absorbing those elements in those
// positions into the current parent.
// i.e.
// 		Parent 1 - [123456789]
// 		Parent 2 - [987654321]
// 		Crossover point = 3
// 		Crossover Section of Parent1 = [123]
// 		Crossover Section of Parent2 = [987]
//
// 		Parent1 PMX Crossover -
// 			[123456789]
// 			Swap [123] with positions of [987] in parent 1,
// 			Offspring = [987456321]
//
// 		Parent2 PMX Crossover -
// 			[987654321]
// 			Swap [987] with positions of [123] in parent 2,
// 			Offspring = [123654789]
func (genA *GeneticAlgorithm) PartiallyMappedCrossover(population Population) Population {
	offspring := make(Population, 0)

	// Identify breeding pairs, assumes an even number of genes in a candidate pool
	for i := 0; i < len(population); i += 2 {
		// Take copies of the 2 parents
		parent1, parent2 := population[i].Copy(), population[i+1].Copy()

		// Start city
		min, max := 0, len(parent1.Sequence)-1
		var crossoverPoint = 0
		crossoverPoint = genA.RandomEngine.Intn(max-min) + min
		offspring1, offspring2 := pmxInternal(parent1, parent2.Sequence[:crossoverPoint]), pmxInternal(parent2, parent1.Sequence[:crossoverPoint])
		offspring = append(offspring, offspring1, offspring2)
	}
	return offspring
}

func (genA *GeneticAlgorithm) MarkOut(bitstring Bitstring, markoutNumbers Bitstring) Bitstring {
	offspring := bitstring.Copy()
	for index, val := range offspring {
		for _, mark := range markoutNumbers {
			if val == mark {
				offspring[index] = "x"
				break
			}
		}
	}
	return offspring
}

func (genA *GeneticAlgorithm) SortMarksToSection(bitstring Bitstring, start, end int) Bitstring {
	offspring := bitstring.Copy()
	changed := true
	for changed {
		changed = false
		for i := 0; i < len(offspring); i++ {
			if offspring[i] == "x" && i < start {
				changed = true
				if offspring[i+1] == "x" {
					for j := i; j < end; j++ {
						if offspring[j] != "x" {
							offspring[i], offspring[j] = offspring[j], offspring[i]
							break
						}
					}
				} else {
					offspring[i], offspring[i+1] = offspring[i+1], offspring[i]
				}
			} else if offspring[i] == "x" && i >= end {
				changed = true
				if offspring[i-1] == "x" {
					for j := i; j >= start; j-- {
						if offspring[j] != "x" {
							offspring[i], offspring[j] = offspring[j], offspring[i]
							break
						}
					}
				} else {
					offspring[i], offspring[i-1] = offspring[i-1], offspring[i]
				}
			}
		}
	}
	return offspring
}

func (genA *GeneticAlgorithm) OrderedCrossover(population Population) Population {
	offspring := make(Population, 0)

	poolSize := len(population)
	ch := make(chan Genome, poolSize)

	// Identify breeding pairs, assumes an even number of genes in a candidate pool
	for i := 0; i < len(population); i += 2 {
		go func(index int) {
			// Take copies of the 2 parents
			parent1, parent2 := population[index].Copy(), population[index+1].Copy()

			// Start city
			min, max := 0, len(parent1.Sequence)-1
			// Choose 2 indexes
			choice1 := genA.RandomEngine.Intn(max-min) + min
			choice2 := genA.RandomEngine.Intn(max-min) + min

			start := int(math.Min(float64(choice1), float64(choice2)))
			end := int(math.Max(float64(choice1), float64(choice2)))

			go func() {
				markOut := parent2.Sequence[start:end]
				marked := genA.MarkOut(parent1.Sequence, markOut)
				sorted := genA.SortMarksToSection(marked, start, end)
				copy(sorted[start:end], markOut)
				ch <- Genome{sorted.Copy()}
			}()

			go func() {
				markOut := parent1.Sequence[start:end]
				marked := genA.MarkOut(parent2.Sequence, markOut)
				sorted := genA.SortMarksToSection(marked, start, end)
				copy(sorted[start:end], markOut)
				ch <- Genome{sorted.Copy()}
			}()

		}(i)
	}

	for i := 0; i < poolSize; i++ {
		offspring = append(offspring, <-ch)
	}
	return offspring
}
