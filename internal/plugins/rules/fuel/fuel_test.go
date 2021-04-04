package fuel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFuelBurningCalculation(t *testing.T) {
	f1 := calculateFuelState(1, 100, 255)
	f2 := calculateFuelState(1, f1, 255)
	f3 := calculateFuelState(1, f2, 255)
	f4 := calculateFuelState(1, f3, 255)

	assert.Equal(t, float32(99), f1)
	assert.Equal(t, float32(98), f2)
	assert.Equal(t, float32(97), f3)
	assert.Equal(t, float32(96), f4)
}
