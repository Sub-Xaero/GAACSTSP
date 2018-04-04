package ga

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
	"math"
)

type GeneticAlgorithm struct {
	Candidates    Population
	BestCandidate Genome
	Generations   int

	IterationsSinceChange int

	Output       func(a ...interface{})
	RandomEngine *rand.Rand
}

func NewGeneticAlgorithm() GeneticAlgorithm {
	var geneticAlgorithm GeneticAlgorithm
	geneticAlgorithm.SetOutputFunc(PrintToConsole)
	geneticAlgorithm.SetSeed(time.Now().Unix())

	return geneticAlgorithm
}

func (genA *GeneticAlgorithm) SetSeed(seed int64) {
	genA.RandomEngine = rand.New(rand.NewSource(seed))
}

func (genA *GeneticAlgorithm) UpdateBestCandidate(bestOfGeneration Genome, cities map[string]City) {
	if genA.Fitness(bestOfGeneration, cities) > genA.Fitness(genA.BestCandidate, cities) {
		genA.BestCandidate = bestOfGeneration.Copy()
		genA.IterationsSinceChange = 0
	}
}

func (genA *GeneticAlgorithm) FillRandomPopulation(populationSize, candidateLength int, cities map[string]City) Population {
	candidatePool := make(Population, 0)
	for len(candidatePool) < populationSize {
		bitstring, err := genA.GenerateCandidate(candidateLength, cities)
		Check(err)
		candidatePool = append(candidatePool, Genome{bitstring})
	}
	return candidatePool
}

func (genA *GeneticAlgorithm) Summarise(title string, candidatePool Population, cities map[string]City) {
	output := ""
	output += title
	//output += "{"
	//for _, val := range candidatePool {
	//	output += "["
	//	if len(val.Sequence) <= 10 {
	//		output += val.Sequence.String()
	//	} else {
	//		output += fmt.Sprintf("%3.0f", math.Abs(genA.Fitness(val, cities)))
	//	}
	//	output += "]"
	//}
	//output += "}"
	output += " # Genes: " + strconv.Itoa(len(candidatePool))
	output += ", Max: " + strconv.FormatFloat(math.Abs(genA.MaxFitness(candidatePool, cities)), 'f', 2, 64) 
	output += ", Avg: " + strconv.FormatFloat(math.Abs(genA.AverageFitness(candidatePool, cities)), 'f', 2, 64)
	output += ", Best: " + genA.MaxFitnessCandidate(candidatePool, cities).String()
	genA.Output(output)
}

func (genA *GeneticAlgorithm) Run(cities map[string]City, populationSize, bitstringLength, generations int, crossover, mutate, terminateEarly bool) error {
	if genA.Output == nil {
		return errors.New("output func is nil")
	}
	if genA.RandomEngine == nil {
		return errors.New("random generator is not initialised")
	}

	// Init
	genA.Candidates = make(Population, 0)
	genA.Candidates = genA.FillRandomPopulation(populationSize, bitstringLength, cities)
	genA.BestCandidate = genA.Candidates[0].Copy()
	genA.UpdateBestCandidate(genA.MaxFitnessCandidate(genA.Candidates, cities), cities)

	// Run breeding cycles
	for y := 1; y <= generations; y++ {
		var bestCandidateOfGeneration Genome

		bestCandidateOfGeneration = genA.MaxFitnessCandidate(genA.Candidates, cities)
		genA.UpdateBestCandidate(bestCandidateOfGeneration, cities)
		genA.Output("Iteration", y)
		genA.Summarise("Start Population       :", genA.Candidates, cities)

		breedingGround := make(Population, 0)

		// Selection
		breedingGround = genA.Selection(genA.Candidates, cities)
		bestCandidateOfGeneration = genA.MaxFitnessCandidate(genA.Candidates, cities)
		genA.UpdateBestCandidate(bestCandidateOfGeneration, cities)
		genA.Summarise("Selection Offspring    :", breedingGround, cities)

		// Crossover
		if crossover {
			breedingGround = genA.Crossover(breedingGround)
			bestCandidateOfGeneration = genA.MaxFitnessCandidate(genA.Candidates, cities)
			genA.UpdateBestCandidate(bestCandidateOfGeneration, cities)
			genA.Summarise("Crossover Offspring    :", breedingGround, cities)
		}

		// Mutation
		if mutate {
			for index := range breedingGround {
				breedingGround[index] = genA.Mutate(breedingGround[index])
			}
			bestCandidateOfGeneration = genA.MaxFitnessCandidate(genA.Candidates, cities)
			genA.UpdateBestCandidate(bestCandidateOfGeneration, cities)
			genA.Summarise("Mutation Offspring     :", breedingGround, cities)
		}

		genA.Generations++
		genA.IterationsSinceChange++
		genA.Candidates = breedingGround.Copy()
		genA.Summarise("Final Population       :", breedingGround, cities)
		genA.Output()
		genA.Output()

		if terminateEarly && float32(genA.IterationsSinceChange) > float32(generations)*0.25 {
			genA.Output("Termination : Stagnating change")
			break
		}
	}

	genA.Output("Best Candidate Found:", genA.BestCandidate.Sequence, "Fitness:", genA.Fitness(genA.BestCandidate, cities))
	return nil
}
