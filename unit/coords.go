package unit

import (
	"errors"
	"math"
)

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

func NewCoord(latitude, longitude float64) (Coordinate, error) {
	if latitude >= 100.0 || latitude <= -100.0 {
		return Coordinate{}, errors.New("latitude out of bounds")
	}
	if longitude >= 100.0 || longitude <= -100.0 {
		return Coordinate{}, errors.New("longitude out of bounds")
	}
	return Coordinate{Latitude: latitude, Longitude: longitude}, nil
}

// Earth's radius in meters
const radius = 6371000

// degToRad converts degrees to radians
func degToRad(d float64) float64 {
	return d * math.Pi / 180
}

// Distance uses the haversine formula to calculate the shortest distance over the earth's surface between the two given coordinates.
func Distance(from, to Coordinate) int64 {
	lat1 := degToRad(from.Latitude)
	lat2 := degToRad(to.Latitude)

	latD := lat2 - lat1
	lonD := degToRad(to.Longitude - from.Longitude)

	a := math.Pow(math.Sin(latD/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(lonD/2), 2)

	return int64(2 * float64(radius) * math.Atan2(math.Sqrt(a), math.Sqrt(1-a)))
}
