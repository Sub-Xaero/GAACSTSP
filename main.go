package main

import (
	"fmt"
	"flag"
	"strings"
	"os"
	. "subxaero/GAACSTSP/ga"
	"log"
	"math"
	"strconv"
)

func main() {
	var (
		inputFileNamePtr       = flag.String("input", "data/berlin52.tsp", "the path to a TSPLib input file \".tsp\" containing cities to find a solution for")
		optimalRoutePtr        = flag.String("optimal", "", "the path to a TSPLib optimal route file \".opt.tour\" containing an optimal solution to compare against")
		selectionMethodPtr     = flag.String("selectionMethod", "aco", "(default: aco) selection method to use, one of 'aco', 'tournament', 'roulette'")
		crossoverMethodPtr     = flag.String("xoMethod", "ox", "(default: ox) crossover method to use, one of 'ox' (ordered), 'pmx' (partially-mapped)")
		numGenerationsPtr      = flag.Int("generations", 500, "the number of generations to run for")
		sizePtr                = flag.Int("size", 50, "the number of candidates to have in the pool")
		crossoverPtr           = flag.Bool("crossover", true, "(default: false) whether or not the algorithm should use crossover operators")
		mutatePtr              = flag.Bool("mutate", true, "(default: false) whether or not the algorithm should use mutation operators")
		terminateEarlyPtr      = flag.Bool("terminateEarly", false, "(default: false) whether or not the algorithm should terminate early if stagnation is detected")
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

	var selectionMethod string
	switch strings.ToLower(*selectionMethodPtr) {
	case "aco":
		fallthrough
	case "roulette":
		fallthrough
	case "tournament":
		selectionMethod = strings.ToLower(*selectionMethodPtr)

	default:
		log.Fatal("method flag specified but was not a recognised value. Please use -h for help")
	}

	var crossoverMethod string
	switch strings.ToLower(*crossoverMethodPtr) {
	case "ox":
		fallthrough
	case "pmx":
		crossoverMethod = strings.ToLower(*crossoverMethodPtr)

	default:
		log.Fatal("xomethod flag specified but was not a recognised value. Please use -h for help")
	}

	var (
		populationSize      = *sizePtr
		generations         = *numGenerationsPtr
		strLength           = len(cities) + 1
		crossover           = *crossoverPtr
		mutate              = *mutatePtr
		terminateEarly      = *terminateEarlyPtr
		terminatePercentage = float64(*terminatePercentagePtr) / 100.0
	)
	ga.Run(cities, populationSize, strLength, generations, crossover, mutate, terminateEarly, terminatePercentage, selectionMethod, crossoverMethod)
	if len(optimal.Sequence) != 0 {
		fmt.Println("Optimal              :", strconv.FormatFloat(math.Abs(ga.Fitness(optimal, cities)), 'f', 2, 64))
	}
	fmt.Println("Best Found           :", strconv.FormatFloat(math.Abs(ga.Fitness(ga.BestCandidate, cities)), 'f', 2, 64))
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Println("populationSize", populationSize)
	fmt.Println("generations", generations)
	fmt.Println("strLength", strLength)
	fmt.Println("selection", selectionMethod)
	fmt.Println("crossover", crossover, crossoverMethod)
	fmt.Println("mutate", mutate, "inversion")
	fmt.Println("terminateEarly", terminateEarly)
	fmt.Println("terminatePercentage", terminatePercentage)
}
