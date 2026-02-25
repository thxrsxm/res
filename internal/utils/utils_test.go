package utils

import "testing"

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"positive number", 5, 5},
		{"negative number", -5, 5},
		{"zero", 0, 0},
		{"large positive", 1000000, 1000000},
		{"large negative", -1000000, 1000000},
		{"min int value", -922337203685477580, 922337203685477580},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Abs(tt.input)
			if result != tt.expected {
				t.Errorf("Abs(%d) = %d; want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnsignedSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "single element",
			input:    []int{5},
			expected: []int{5},
		},
		{
			name:     "already sorted positive",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "reverse sorted positive",
			input:    []int{5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "already sorted negative",
			input:    []int{-1, -2, -3, -4, -5},
			expected: []int{-1, -2, -3, -4, -5},
		},
		{
			name:     "reverse sorted negative",
			input:    []int{-5, -4, -3, -2, -1},
			expected: []int{-1, -2, -3, -4, -5},
		},
		{
			name:     "mixed positive and negative",
			input:    []int{3, -2, 5, -1, 4},
			expected: []int{-1, -2, 3, 4, 5},
		},
		{
			name:     "with zeros",
			input:    []int{0, -3, 2, -1, 0, 4},
			expected: []int{0, 0, -1, 2, -3, 4},
		},
		{
			name:     "duplicate values",
			input:    []int{-2, 1, -2},
			expected: []int{1, -2, -2},
		},
		{
			name:     "large numbers",
			input:    []int{1000000, -2000000, 500000, -1500000},
			expected: []int{500000, 1000000, -1500000, -2000000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := make([]int, len(tt.input))
			copy(actual, tt.input)
			UnsignedSort(actual)
			if len(actual) != len(tt.expected) {
				t.Errorf("length mismatch: got %d, want %d", len(actual), len(tt.expected))
				return
			}
			for i := range actual {
				if actual[i] != tt.expected[i] {
					t.Errorf("UnsignedSort(%v) = %v; want %v", tt.input, actual, tt.expected)
					return
				}
			}
		})
	}
}
