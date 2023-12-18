package observable

import (
	"fmt"
	"testing"
)

func TestFilterPercentageChange(t *testing.T) {
	tests := []struct {
		current, new uint8
		expected     bool
	}{
		{50, 50, false},   // Same value, should return false
		{0, 0, false},     // Same value, should return false
		{10, 20, true},    // Different value within the valid range, should return true
		{75, 75, false},   // Same value, should return false
		{80, 90, true},    // Different value within the valid range, should return true
		{100, 120, false}, // New value exceeds the valid range, should return false
		{100, 190, false}, // New value exceeds the valid range, should return false
	}

	filter := DistinctPercentageChange()

	for _, test := range tests {
		t.Run(fmt.Sprintf("Current:%d New:%d", test.current, test.new), func(t *testing.T) {
			result := filter(test.current, test.new)
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}

func TestDistinctBooleanChange(t *testing.T) {
	tests := []struct {
		current, new bool
		expected     bool
	}{
		{true, true, false},   // Same value, should return false
		{false, false, false}, // Same value, should return false
		{true, false, true},   // Different value, should return true
		{false, true, true},   // Different value, should return true
	}

	filter := DistinctBooleanChange()

	for _, test := range tests {
		t.Run(fmt.Sprintf("Current:%v New:%v", test.current, test.new), func(t *testing.T) {
			result := filter(test.current, test.new)
			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}
