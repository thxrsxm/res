// Package utils provides utility functions for the resolution theorem prover.
package utils

// Abs returns the absolute value of an integer x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// UnsignedSort sorts a slice of integers by their absolute values in ascending order.
// This is used to provide consistent ordering of literals in clause string representations.
//
// The sorting is stable and uses an insertion sort algorithm, which is efficient
// for small slices (typical for clauses with few literals).
//
// Example:
//
//	Input:  []int{-3, 1, -2, 4}
//	Output: []int{1, -2, -3, 4}  // Sorted by absolute values: 1, 2, 3, 4
func UnsignedSort(values []int) {
	for i := 1; i < len(values); i++ {
		key := values[i]
		j := i - 1
		for j >= 0 && Abs(values[j]) > Abs(key) {
			values[j+1] = values[j]
			j = j - 1
		}
		values[j+1] = key
	}
}
