package fuel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFuelBurningCalculation(t *testing.T) {
	f1 := calculateFuelState(0.1, 100, 255)
	f2 := calculateFuelState(0.1, f1, 255)
	f3 := calculateFuelState(0.1, f2, 255)
	f4 := calculateFuelState(0.1, f3, 255)
	f5 := calculateFuelState(0.1, f3, 255)

	assert.Equal(t, Liter(74.5), f1)
	assert.Equal(t, Liter(49), f2)
	assert.Equal(t, Liter(23.5), f3)
	assert.Equal(t, Liter(0), f4)
	assert.Equal(t, Liter(0), f5)
}
