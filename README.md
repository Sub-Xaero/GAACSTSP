# A Genetic Algorithm for the Traveling Salesman Problem using an Ant Colony Optimisation inspired selection operator, and Partially Mapped Crossover

Typical representations of the Travelling Salesman Problem (TSP), as a specialised
instance of a combinatorial optimisation problem, suffer a handicap in comparison
to other artificial intelligence techniques when encoded in a Genetic Algorithm
(GA). The handicap is such that encoding genes as a permutation constrained
string restricts use of the GA’s core operators, crossover and mutation. Classical GA
operators are typically constraint agnostic, or deal specifically with the constraints
relating to the value, range, and encoding of individual chromosomes rather than
the ’gene’ as a whole such as TSPs require. Therefore, alternative operators that
consider the constraints of the problem at hand are required in order to use a GA for
the purposes of solving a TSP.

This paper sets forth a solution that makes use of Partially-Mapped Crossover as
proposed and implemented by Goldberg, Lingle et al. (1985), and a new operator
for parent selection that takes inspiration from the way in which Ant Colony
Optimisation algorithms build a candidate solution by way of iterative traversal.

This algorithm processes data in TSPLib format (Reinelt, 1991). The main examples of which can be found at the official homepage:
http://comopt.ifi.uni-heidelberg.de/software/TSPLIB95/

## Useage
The executable can be invoked with the following command-line flags
```
  -crossover
        whether or not the algorithm should use crossover operators (default true)
  -generations int
        the number of generations to run for (default 500)
  -input string
        the path to a TSPLib input file ".tsp" containing cities to find a solution for (default "data/berlin52.tsp")
  -length int
        the length of a candidate solution. Must be equal to the number of cities + 1 (default 53)
  -method string
        Selection method to use, one of 'ACO', 'Tournament', 'Roulette' (default "ACO")
  -mutate
        whether or not the algorithm should use mutation operators (default true)
  -optimal string
        the path to a TSPLib optimal route file ".opt.tour" containing an optimal solution to compare against
  -size int
        the number of candidates to have in the pool (default 50)
  -terminateEarly
        whether or not the algorithm should terminate early if stagnation is detected
  -terminatePercentage int
        percentage of the specified no. of generations (default 500), should the algorithm terminate if change has not been detected in that time (default 25)
```