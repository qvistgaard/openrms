package types

import (
	"testing"
)

func TestIdFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected Id
		err      error
	}{
		{"0", Id(0), nil},
		{"12345", Id(12345), nil},
		{"4294967295", Id(4294967295), nil}, // Max uint32 value
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := IdFromString(test.input)
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
