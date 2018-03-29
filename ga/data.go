package ga

import (
	"os"
	"bufio"
	"strings"
	"strconv"
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

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Compare(line, "NODE_COORD_SECTION") == 0 {
			break
		}
	}
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Compare(line, "EOF") == 0 {
			break
		} else {
			count++
		}
		data := strings.Split(scanner.Text(), " ")
		x, convertErr := strconv.ParseFloat(data[1], 64)
		Check(convertErr)

		y, convertErr := strconv.ParseFloat(data[2], 64)
		Check(convertErr)

		Check(scanner.Err())
		cities[strconv.Itoa(count)] = City{x, y}
	}

	return cities
}
