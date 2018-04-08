package ga

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
	"math"
	"log"
)

// Main structure of program.
// Compartmentalised such that multiple instances can theoretically be initialised and run simultaneously, though the results are undefined.
//
// This structure keeps track of the candidates to be overwritten and passed between generations,
// the number of generations the algorithm was run for, the best candidate ever found, and the stagnation counter IterationsSinceChange.
//
// The structure also defines an Output closure, that can be used to change where the algorithm logs information to, as
// well as RandomEngine, which can be independently configured to change how the algorithm generates "randomness"
type GeneticAlgorithm struct {
	Candidates    Population
	BestCandidate Genome
	Generations   int

	IterationsSinceChange int

	Output       func(a ...interface{})
	RandomEngine *rand.Rand
}

// Initialises and returns a GeneticAlgorithm structure with all the default configurations applied
func NewGeneticAlgorithm() GeneticAlgorithm {
	var geneticAlgorithm GeneticAlgorithm
	geneticAlgorithm.SetOutputFunc(PrintToConsole)
	geneticAlgorithm.SetSeed(time.Now().Unix())

	return geneticAlgorithm
}

// Change the seed that the randomiser pulls random numbers out of.
func (genA *GeneticAlgorithm) SetSeed(seed int64) {
	genA.RandomEngine = rand.New(rand.NewSource(seed))
}

// Given the best genome found in a generation, if the new fitness is better - update the best candidate ever found stored in the GeneticAlgorithm
// structure and reset the IterationsSinceChange counter, used in the detection of improvement stagnation.
func (genA *GeneticAlgorithm) UpdateBestCandidate(bestOfGeneration Genome, cities map[string]City) {
	if genA.Fitness(bestOfGeneration, cities) > genA.Fitness(genA.BestCandidate, cities) {
		genA.BestCandidate = bestOfGeneration.Copy()
		genA.IterationsSinceChange = 0
	}
}

// Wrapper around GenerateCandidate(), creates a population of size "populationSize", where each candidate is of length
// "candidateLength". Returns the resulting candidate pool.
func (genA *GeneticAlgorithm) FillRandomPopulation(populationSize, candidateLength int, cities map[string]City) Population {
	candidatePool := make(Population, 0)
	for len(candidatePool) < populationSize {
		bitstring, err := genA.GenerateCandidate(candidateLength, cities)
		Check(err)
		candidatePool = append(candidatePool, Genome{bitstring})
	}
	return candidatePool
}

// Given a candidate pool or population, print out a summary of the contents:
//
// The text contained in "title",
// Number of Genes in pool,
// Average and Maximum fitness across the population,
// and a visual representation of the best candidate of the generation.
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
	maxCandidate := genA.MaxFitnessCandidate(candidatePool, cities)
	output += " # Genes: " + strconv.Itoa(len(candidatePool))
	output += ", Max: " + strconv.FormatFloat(math.Abs(genA.MaxFitness(candidatePool, cities)), 'f', 2, 64)
	output += ", Avg: " + strconv.FormatFloat(math.Abs(genA.AverageFitness(candidatePool, cities)), 'f', 2, 64)
	output += ", Best: " + strconv.FormatFloat(math.Abs(genA.Fitness(maxCandidate, cities)), 'f', 2, 64) + " " + maxCandidate.String()
	genA.Output(output)
}

// Run the algorithm, with the passed parameters. Performs the stages of
// Selection, Crossover and Mutation, as determined by the boolean flag parameters crossover and mutate.
// terminateEarly if set to true, will attempt to detect stagnation of improvement. If 25% of generations have passed
// and no improvement has been made, algorithm will terminate early and return the best thus far.
func (genA *GeneticAlgorithm) Run(cities map[string]City, populationSize, bitstringLength, generations int, crossover, mutate, terminateEarly bool, terminatePercentage float64, method string) error {
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
		switch method {
		case "aco":
			breedingGround = genA.ACOSelection(genA.Candidates, cities)
		case "roulette":
			breedingGround = genA.RouletteSelection(genA.Candidates, cities)
		case "tournament":
			breedingGround = genA.TournamentSelection(genA.Candidates, cities)
		default:
			log.Fatal("method not a recognised value")
		}
		bestCandidateOfGeneration = genA.MaxFitnessCandidate(genA.Candidates, cities)
		genA.UpdateBestCandidate(bestCandidateOfGeneration, cities)
		genA.Summarise("Selection Offspring    :", breedingGround, cities)

		// Crossover
		if crossover {
				breedingGround = genA.PartiallyMappedCrossover(genA.Candidates)
			bestCandidateOfGeneration = genA.MaxFitnessCandidate(genA.Candidates, cities)
			genA.UpdateBestCandidate(bestCandidateOfGeneration, cities)
			genA.Summarise("Crossover Offspring    :", breedingGround, cities)
		}

		// Mutation
		if mutate {
			for index := range breedingGround {
				breedingGround[index] = genA.InversionMutate(breedingGround[index])
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

		if terminateEarly && float64(genA.IterationsSinceChange) > float64(generations)*terminatePercentage {
			genA.Output("Termination : Stagnating change")
			break
		}
	}

	genA.Output("Best Candidate Found :", strconv.FormatFloat(math.Abs(genA.Fitness(genA.BestCandidate, cities)), 'f', 2, 64), genA.BestCandidate, )
	return nil
}
