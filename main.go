package main

import (
	. "subxaero/GAACSTSP/ga"
)

func main() {
	cities := LoadTSPLib("data/xqf131.tsp")
	var ga = NewGeneticAlgorithm()
	var (
		populationSize = 200
		generations    = 2000
		strLength      = len(cities)
	)

	ga.Run(cities, populationSize, strLength, generations, true, true)
}
