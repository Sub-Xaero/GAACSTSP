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

	re := regexp.MustCompile(`\s+`)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "FIXED_EDGES_SECTION") {
			panic("Program does not support fixed edges")
		}

		if strings.Contains(line, "EDGE_WEIGHT_TYPE") && !strings.Contains(line, "EUC_2D") {
			panic("Data file does not contain co-ordinates compatible with EUC_2D fitness function")
		}

		if strings.Compare(line, "NODE_COORD_SECTION") == 0 {
			break
		}
	}
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
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

		Check(scanner.Err())
		cities[strconv.Itoa(count)] = City{x, y}
	}

	return cities
}
