package tests

import (
	"packcalculator/helpers"
	"reflect"
	"testing"
)

// Unit tests for the algorithm
func TestFindBestPackCombination(t *testing.T) {
	tests := []struct {
		name   string
		packs  []int
		amount int
		expect map[int]int
	}{
		{
			name:   "Exact match",
			packs:  []int{250, 500, 1000, 2000, 5000},
			amount: 5000,
			expect: map[int]int{5000: 1},
		},
		{
			name:   "Minimal overfill",
			packs:  []int{250, 500, 1000, 2000, 5000},
			amount: 260,
			expect: map[int]int{500: 1},
		},
		{
			name:   "Fewest packs when tied",
			packs:  []int{250, 500, 1000, 2000, 5000},
			amount: 750,
			expect: map[int]int{500: 1, 250: 1},
		},
		{
			name:   "Edge case from spec",
			packs:  []int{23, 31, 53},
			amount: 500000,
			expect: map[int]int{23: 2, 31: 7, 53: 9429},
		},
		{
			name:   "No solution",
			packs:  []int{},
			amount: 100,
			expect: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := helpers.FindOptimalPackCombination(tc.packs, tc.amount)
			if !reflect.DeepEqual(result, tc.expect) {
				t.Errorf("Test %s failed: expected %v, got %v", tc.name, tc.expect, result)
			}
		})
	}
}
