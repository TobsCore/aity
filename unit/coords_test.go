package unit

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

var karlsruhe = Coordinate{Latitude: 49.013482, Longitude: 8.404235}
var berlin = Coordinate{Latitude: 52.519067, Longitude: 13.406528}

func TestDistanceReversal(t *testing.T) {
	distance := Distance(karlsruhe, berlin)
	distanceRev := Distance(berlin, karlsruhe)

	if distance != distanceRev {
		t.Error("It shouldn't matter, from which direction to calc the distance")
	}
}

func TestDistanceKarlsruheBerlin(t *testing.T) {
	var distance = Distance(karlsruhe, berlin)
	distanceExpected := 524000
	var delta int64 = 1000
	assert.True(t, int64(math.Abs(float64(distance)-float64(distanceExpected))) < delta, fmt.Sprintf("The difference between the expected and the actual value is too large. Got %d", distance))
}
