package cmd

import (
	"strings"
	"testing"
)

func TestGetStringFromIoReader(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      string
		expectedError error
	}{
		{
			name:     "Valid case",
			input:    "Hello, World!\n",
			expected: "Hello, World!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			result, err := getStringFromIoReader(r)

			if result != tt.expected {
				t.Errorf("Expected result %q, but got %q", tt.expected, result)
			}

			if (err == nil && tt.expectedError != nil) || (err != nil && tt.expectedError == nil) {
				t.Errorf("Expected error %v, but got %v", tt.expectedError, err)
			} else if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error %v, but got %v", tt.expectedError, err)
			}
		})
	}
}
