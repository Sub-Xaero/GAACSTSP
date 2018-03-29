package ga

import (
	"math"
)

func EuclideanDistance(c1, c2 City) float64 {
	var difX = math.Pow(c1.x-c2.x, 2)
	var difY = math.Pow(c1.y-c2.y, 2)
	var sqSum = math.Sqrt(difX + difY)
	return sqSum
}

