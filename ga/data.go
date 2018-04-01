package ga

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"regexp"
)

type City struct {
	x, y float64
}

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
