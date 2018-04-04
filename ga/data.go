package ga

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"regexp"
)

// Storage type for co-ordinates loaded from TSPLib files
type City struct {
	X, Y float64
}

// Function to parse and construct a candidate from an ".opt.tour" file in TSPLib format
//
// TSBLib can be found at:
//http://www.iwr.uni-heidelberg.de/groups/comopt/software/TSPLIB95/
func LoadTSPOptTour(filename string) Genome {
	sequence := make(Bitstring, 0)

	file, err := os.Open(filename)
	Check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// Regex-ers
	fileType := regexp.MustCompile(`^\s*TYPE\s*:`)
	tourSection := regexp.MustCompile(`^\s*TOUR_SECTION`)

	// Read file header
	for scanner.Scan() {
		Check(scanner.Err())

		line := scanner.Text()
		Check(scanner.Err())

		if fileType.MatchString(line) && !strings.Contains(line, "TOUR") {
			panic("Input file being loaded as optimal tour file")
		}

		if tourSection.MatchString(line) {
			break
		}
	}

	// Read file content
	re := regexp.MustCompile(`\s+`)
	for scanner.Scan() {
		Check(scanner.Err())

		line := scanner.Text()
		Check(scanner.Err())

		line = re.ReplaceAllString(line, " ")
		line = strings.Trim(line, " ")

		if strings.Compare(line, "EOF") == 0 {
			break
		}
		data := strings.Split(line, " ")
		for _, val := range data {
			_, convertErr := strconv.ParseInt(val, 10, 64)
			Check(convertErr)
			sequence = append(sequence, val)
		}
	}
	return Genome{Sequence: sequence}
}

// Function to parse and construct a map of cities from a ".tsp" file in TSPLib format.
//
// TSBLib can be found at:
//http://www.iwr.uni-heidelberg.de/groups/comopt/software/TSPLIB95/
//
// Key is a string due to the encoding of the chromosomes in the Bitstring as a string, each chromosome represents a visit
// to the city located at that index in the map.
func LoadTSPLib(filename string) map[string]City {
	cities := make(map[string]City, 0)
	file, err := os.Open(filename)
	Check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Regex-ers
	fileType := regexp.MustCompile(`^\s*TYPE\s*:`)
	fixedEdgesSection := regexp.MustCompile(`^\s*FIXED_EDGES_SECTION`)
	edgeWeightType := regexp.MustCompile(`^\s*EDGE_WEIGHT_TYPE\s*:`)
	nodeCoordSection := regexp.MustCompile(`^\s*NODE_COORD_SECTION`)

	// Read file header
	for scanner.Scan() {
		Check(scanner.Err())

		line := scanner.Text()
		Check(scanner.Err())

		if fileType.MatchString(line) && !strings.Contains(line, "TSP") {
			panic("Optimal tour file being loaded as input file")
		}

		if fixedEdgesSection.MatchString(line) {
			panic("Program does not support fixed edges")
		}

		if edgeWeightType.MatchString(line) && !strings.Contains(line, "EUC_2D") {
			panic("Data file does not contain co-ordinates compatible with EUC_2D fitness function")
		}

		if nodeCoordSection.MatchString(line) {
			break
		}
	}
	count := 0

	// Read file content
	re := regexp.MustCompile(`\s+`)
	for scanner.Scan() {
		Check(scanner.Err())

		line := scanner.Text()
		Check(scanner.Err())

		line = re.ReplaceAllString(line, " ")
		line = strings.Trim(line, " ")

		if strings.Compare(line, "EOF") == 0 {
			break
		} else {
			count++
		}
		data := strings.Split(line, " ")

		x, convertErr := strconv.ParseFloat(data[1], 64)
		Check(convertErr)

		y, convertErr := strconv.ParseFloat(data[2], 64)
		Check(convertErr)

		cities[strconv.Itoa(count)] = City{x, y}
	}

	return cities
}
