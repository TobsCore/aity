package unit

import (
	"math"
)

type Coord struct {
	lat float64
	lon float64
}

// Earth's radius in meters
const radius = Meter(6371000)

// degToRad converts degrees to radians
func degToRad(d float64) float64 {
	return d * math.Pi / 180
}

// Distance uses the haversine formula to calculate the shortest distance over the earth's surface between the two given coordinates.
func Distance(from, to Coord) Meter {
	lat1 := degToRad(from.lat)
	lat2 := degToRad(to.lat)

	latD := lat2 - lat1
	lonD := degToRad(to.lon - from.lon)

	a := math.Pow(math.Sin(latD/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(lonD/2), 2)

	return Meter(2 * float64(radius) * math.Atan2(math.Sqrt(a), math.Sqrt(1-a)))
}
