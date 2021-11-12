package fuel

import (
	"github.com/qvistgaard/openrms/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFuelBurningCalculation(t *testing.T) {
	f1 := calculateFuelState(0.1, 0, 50)
	f2 := calculateFuelState(0.1, 0, 100)

	assert.Equal(t, types.Liter(0.05), f1)
	assert.Equal(t, types.Liter(0.1), f2)

}

func TestRefuelFuelCalculation(t *testing.T) {
	f1, full := calculateRefuellingValue(20, 10)
	assert.False(t, full)
	assert.Equal(t, types.Liter(10), f1)

	f2, full2 := calculateRefuellingValue(10, 10)
	assert.True(t, full2)
	assert.Equal(t, types.Liter(0), f2)

	f3, full3 := calculateRefuellingValue(5, 10)
	assert.True(t, full3)
	assert.Equal(t, types.Liter(0), f3)

}
