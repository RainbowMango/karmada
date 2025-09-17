/*
Copyright 2025 The Karmada Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helper

import (
	"container/heap"
)

// Party represents a party in the Webster apportionment method.
type Party struct {
	// Name is the name of the party.
	Name string
	// Votes is the number of votes the party received.
	Votes int64
	// Seats is the number of seats currently assigned to the party.
	Seats int32
}

// WebsterPriorityQueue implements heap.Interface for Party using the Webster (Sainte-Laguë) method.
// The party with the highest Webster priority (Votes/(2*Seats+1)) is at the top.
type WebsterPriorityQueue []Party

// Check if our WebsterPriorityQueue implements necessary interfaces
var _ heap.Interface = &WebsterPriorityQueue{}

// Len returns the number of parties in the queue.
func (pq *WebsterPriorityQueue) Len() int {
	return len(*pq)
}

// Less compares two parties by their Webster priority.
// The party with the higher priority comes first.
func (pq *WebsterPriorityQueue) Less(i, j int) bool {
	// In the Webster method, compare the priority of two parties: the one with the higher value of Votes/(2*Seats+1) gets the next seat.
	// We avoid using division here to prevent loss of precision and potential rounding errors
	// that can occur with floating point arithmetic, especially when dealing with large integers.
	// Instead, we use cross-multiplication to compare the two fractions directly:
	// pq[i].Votes * (2*pq[j].Seats+1) > pq[j].Votes * (2*pq[i].Seats+1)
	// This approach ensures the comparison is both accurate and efficient, as it only involves integer operations.
	// It also avoids the performance overhead and subtle bugs that can arise from floating point division.
	left := (*pq)[i].Votes * int64(2*(*pq)[j].Seats+1)
	right := (*pq)[j].Votes * int64(2*(*pq)[i].Seats+1)
	if left == right {
		return (*pq)[i].Votes >= (*pq)[j].Votes
	}
	return left > right
}

// Swap swaps two parties in the queue.
func (pq *WebsterPriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
}

// Push adds a new party to the queue.
func (pq *WebsterPriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(Party))
}

// Pop removes and returns the party with the highest priority.
func (pq *WebsterPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// AllocateWebsterSeats allocates new seats using the Webster method.
// newSeats: number of new seats to allocate across all parties.
// parties: slice of Party with Name and Votes set, Seats may already be assigned.
func AllocateWebsterSeats(newSeats int32, parties []Party) []Party {
	// Initialize queue with all parties, preserve existing Seats values
	pq := make(WebsterPriorityQueue, len(parties))
	nameToIndex := make(map[string]int, len(parties))
	for i, p := range parties {
		pq[i] = p
		nameToIndex[p.Name] = i
	}
	heap.Init(&pq)

	remaining := newSeats
	if remaining <= 0 {
		return parties
	}

	for s := int32(0); s < remaining; s++ {
		// Pop the party with the highest priority
		top := heap.Pop(&pq).(Party)
		top.Seats++
		heap.Push(&pq, top)
	}

	// Collect results in the same order as input
	result := make([]Party, len(parties))
	for _, p := range pq {
		if idx, ok := nameToIndex[p.Name]; ok {
			result[idx] = p
		}
	}
	return result
}
