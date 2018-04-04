package ga

// Function to actually perform the mapping operation described in SinglePointPMX()
func pmx(gene Genome, bitstring Bitstring) Genome {
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
func (genA *GeneticAlgorithm) SinglePointPMX(population Population) Population {
	offspring := make(Population, 0)

	// Identify breeding pairs, assumes an even number of genes in a candidate pool
	for i := 0; i < len(population); i += 2 {
		// Take copies of the 2 parents
		parent1, parent2 := population[i].Copy(), population[i+1].Copy()

		// Start city
		min, max := 0		, len(parent1.Sequence)-1
		var crossoverPoint = 0
		if max < 10 {
			crossoverPoint = genA.RandomEngine.Intn(max-min) + min
		} else {
			crossoverPoint = 10
		}
		offspring1, offspring2 := pmx(parent1, parent2.Sequence[:crossoverPoint]), pmx(parent2, parent1.Sequence[:crossoverPoint])
		offspring = append(offspring, offspring1, offspring2)
	}
	return offspring
}

// Calls the chosen Crossover operator
func (genA *GeneticAlgorithm) Crossover(candidatePool Population) Population {
	return genA.SinglePointPMX(candidatePool)
}
