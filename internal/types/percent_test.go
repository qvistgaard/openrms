package types

import (
	"testing"
)

func TestPercentFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected Percent
		err      error
	}{
		{"0", NewPercentFromUint8(0), nil},
		{"50", NewPercentFromUint8(50), nil},
		{"100", NewPercentFromUint8(100), nil},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := PercentFromString(test.input)
			if err != nil && test.err == nil {
				t.Errorf("Expected no error, but got an error: %v", err)
			} else if err == nil && test.err != nil {
				t.Errorf("Expected an error, but got no error")
			} else if err != nil && test.err != nil && err.Error() != test.err.Error() {
				t.Errorf("Expected error message '%s', but got '%s'", test.err.Error(), err.Error())
			}

			if result != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, result)
			}
		})
	}
}
