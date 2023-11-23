package fuel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFuelBurningCalculation(t *testing.T) {
	f1 := calculateFuelState(0.1, 0, 50)
	f2 := calculateFuelState(0.1, 0, 100)

	assert.Equal(t, float32(0.05), f1)
	assert.Equal(t, float32(0.1), f2)

}
