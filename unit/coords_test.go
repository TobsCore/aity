package unit

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"math"
)

var karlsruhe = Coord{lat: 49.013482, lon: 8.404235}
var berlin = Coord{lat: 52.519067, lon: 13.406528}

func TestDistanceReversal(t *testing.T) {
	distance := Distance(karlsruhe, berlin)
	distanceRev := Distance(berlin, karlsruhe)

	if distance != distanceRev {
		t.Error("It shouldn't matter, from which direction to calc the distance")
	}
}

func TestDistanceKarlsruheBerlin(t *testing.T) {
	var distance = Distance(karlsruhe, berlin)
	distanceExpected := Meter(524000)
	delta := Meter(1000)
	assert.True(t, Meter(math.Abs(float64(distance)-float64(distanceExpected))) < delta, "The difference between the expected and the actual value is too large. Got " + distance.String())
}
