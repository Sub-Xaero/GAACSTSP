package main

import (
	"fmt"
	"flag"
	"strings"
	"os"
	. "subxaero/GAACSTSP/ga"
	"log"
)

func main() {
	var (
		inputFileNamePtr       = flag.String("input", "data/berlin52.tsp", "the path to a TSPLib input file \".tsp\" containing cities to find a solution for")
		optimalRoutePtr        = flag.String("optimal", "", "the path to a TSPLib optimal route file \".opt.tour\" containing an optimal solution to compare against")
		methodPtr              = flag.String("method", "ACO", "Selection method to use, one of 'ACO', 'Tournament', 'Roulette'")
		numGenerationsPtr      = flag.Int("generations", 500, "the number of generations to run for")
		lengthPtr              = flag.Int("length", 53, "the length of a candidate solution. Must be equal to the number of cities + 1")
		sizePtr                = flag.Int("size", 50, "the number of candidates to have in the pool")
		crossoverPtr           = flag.Bool("crossover", true, "whether or not the algorithm should use crossover operators")
		mutatePtr              = flag.Bool("mutate", true, "whether or not the algorithm should use mutation operators")
		terminateEarlyPtr      = flag.Bool("terminateEarly", false, "whether or not the algorithm should terminate early if stagnation is detected")
		terminatePercentagePtr = flag.Int("terminatePercentage", 25, "percentage of the specified no. of generations (default 500), should the algorithm terminate if change has not been detected in that time")
	)

	flag.Parse()

	cities := LoadTSPLib(*inputFileNamePtr)
	var ga = NewGeneticAlgorithm()

	var optimal Genome
	if *optimalRoutePtr != "" {
		optimal = LoadTSPOptTour(*optimalRoutePtr)
	} else {
		autoFileName := strings.Replace(*inputFileNamePtr, ".tsp", ".opt.tour", -1)
		if _, err := os.Stat(autoFileName); err == nil {
			optimal = LoadTSPOptTour(autoFileName)
		}
	}

	var method string

	switch *methodPtr {
	case "ACO":
		fallthrough
	case "Roulette":
		fallthrough
	case "Tournament":
		method = *methodPtr

	default:
		log.Fatal("method flag specified but was not a recognised value. Please use -h for help")
	}

	var (
		populationSize      = *sizePtr
		generations         = *numGenerationsPtr
		strLength           = *lengthPtr
		crossover           = *crossoverPtr
		mutate              = *mutatePtr
		terminateEarly      = *terminateEarlyPtr
		terminatePercentage = float64(*terminatePercentagePtr) / 100.0
	)
	ga.Run(cities, populationSize, strLength, generations, crossover, mutate, terminateEarly, terminatePercentage, method)
	if len(optimal.Sequence) != 0 {
		fmt.Println("Optimal    : ", ga.Fitness(optimal, cities))
	}
	fmt.Println("Best Found : ", ga.Fitness(ga.BestCandidate, cities))
}
