package ga

func PMX(gene Genome, bitstring Bitstring) Genome {
	output := make(Bitstring, len(gene.Sequence))
	positions := make(map[string]int)
	swapNumbers := make(Bitstring, len(bitstring))
	copy(swapNumbers, bitstring)
	geneCopy := gene.Copy()
	for index, city := range geneCopy.Sequence {
		positions[city] = index
	}

	sectionEnd := len(swapNumbers)
	section := geneCopy.Sequence[:sectionEnd]
	for index := range section {
		swapFrom := section[index]
		swapWith := swapNumbers[index]
		posFrom, posWith := positions[swapFrom], positions[swapWith]
		positions[swapFrom] = posWith
		positions[swapWith] = posFrom
	}

	for city, index := range positions {
		output[index] = city
	}
	geneOutput := Genome{output}
	return geneOutput
}

func (genA *GeneticAlgorithm) SinglePointPMX(population Population) Population {
	offspring := make(Population, 0)

	for i := 0; i < len(population); i += 2 {
		parent1, parent2 := population[i].Copy(), population[i+1].Copy()

		min, max := 1, len(parent1.Sequence)-1
		var crossoverPoint = 0
		if max < 10 {
			crossoverPoint = genA.RandomEngine.Intn(max-min) + min
		} else {
			crossoverPoint = 10
		}
		offspring1, offspring2 := PMX(parent1, parent2.Sequence[:crossoverPoint]), PMX(parent2, parent1.Sequence[:crossoverPoint])
		offspring = append(offspring, offspring1, offspring2)
	}
	return offspring
}

// SetSelectionFunc changes the selection function to the function specified
func (genA *GeneticAlgorithm) Crossover(candidatePool Population) Population {
	return genA.SinglePointPMX(candidatePool)
}
