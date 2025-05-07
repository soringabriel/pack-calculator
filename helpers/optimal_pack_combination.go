package helpers

import (
	"container/heap"
)

type Solution struct {
	TotalItems int
	TotalPacks int
	PackCounts map[int]int
}

// PriorityQueue implements heap.Interface for Items sorted by (items, then packs).
type PriorityQueue []*Solution

// Implement heap.Interface methods
func (pq PriorityQueue) Len() int { return len(pq) }

// Less prioritizes lower TotalItems, then fewer TotalPacks if tied
func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].TotalItems == pq[j].TotalItems {
		return pq[i].TotalPacks < pq[j].TotalPacks
	}
	return pq[i].TotalItems < pq[j].TotalItems
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(*Solution))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// FindOptimalPackCombination returns the optimal combination of packs that sum to at least `amount`,
// with minimal total items over `amount`, and minimal number of packs.
func FindOptimalPackCombination(packSizes []int, amount int) map[int]int {
	// Initialize priority queue
	pq := &PriorityQueue{}
	heap.Init(pq)

	// visited maps totalItems to the best solution number of packages, to avoid revisiting worse states
	visited := make(map[int]int)
	packCounts := make(map[int]int)

	// Start from 0 items and 0 packs
	heap.Push(pq, &Solution{
		TotalItems: 0,
		TotalPacks: 0,
		PackCounts: packCounts,
	})

	// Start searching solution
	for pq.Len() > 0 {
		solution := heap.Pop(pq).(*Solution) // This will always extract the minimal TotalItems, TotalPacks solution from priorityQueue (based on the Less method implementation)

		// When we reach or exceed the amount we got our solution (since the popped one from the queue its always the minimal TotalItems and TotalPacks)
		if solution.TotalItems >= amount {
			return solution.PackCounts
		}

		// Add new solutions to the queue
		for _, size := range packSizes {
			nextTotal := solution.TotalItems + size
			nextPacks := solution.TotalPacks + 1

			// Skip if we’ve seen this total before with fewer packs (=> this solution is a worse solution)
			if prevPacks, ok := visited[nextTotal]; ok && prevPacks <= nextPacks {
				continue
			}
			visited[nextTotal] = nextPacks

			// Copy and update pack usage
			newCounts := make(map[int]int)
			for k, v := range solution.PackCounts {
				newCounts[k] = v
			}
			newCounts[size]++

			// Push new solution into the queue
			heap.Push(pq, &Solution{
				TotalItems: nextTotal,
				TotalPacks: nextPacks,
				PackCounts: newCounts,
			})
		}
	}

	// No solution found
	return nil
}
