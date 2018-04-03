package main

import (
	. "subxaero/GAACSTSP/ga"
	"fmt"
)

func main() {
	cities := LoadTSPLib("data/berlin52.tsp")

	var ga = NewGeneticAlgorithm()
	optimal := LoadTSPOptTour("data/berlin52.opt.tour")

	var (
		populationSize = 200
		generations    = 500
		strLength      = len(cities) + 1
	)
	ga.Run(cities, populationSize, strLength, generations, true, true, false)
	fmt.Println("Optimal    : ", ga.Fitness(optimal, cities))
	fmt.Println("Best Found : ", ga.Fitness(ga.BestCandidate, cities))
}
