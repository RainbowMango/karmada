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

// WebsterPriorityQueue implements heap.Interface using the Webster (Sainte-Laguë) method.
// The party with the highest Webster priority (Votes/(2*Seats+1)) is at the top.
type WebsterPriorityQueue struct {
	parties []Party
}

// Check if our WebsterPriorityQueue implements necessary interfaces
var _ heap.Interface = &WebsterPriorityQueue{}

// Len returns the number of parties in the queue.
func (pq *WebsterPriorityQueue) Len() int {
	return len(pq.parties)
}

// Less compares two parties by their Webster priority.
// The party with the higher priority comes first.
func (pq *WebsterPriorityQueue) Less(i, j int) bool {
	// In the Webster method, compare the priority of two parties:
	// the one with the higher value of Votes/(2*Seats+1) gets the next seat.
	iPriority := pq.parties[i].Votes / (int64(2*pq.parties[i].Seats + 1))
	jPriority := pq.parties[j].Votes / (int64(2*pq.parties[j].Seats + 1))
	if iPriority == jPriority {
		return (*pq).parties[i].Votes >= (*pq).parties[j].Votes // TODO: replace with tie-breaker
	}
	return iPriority > jPriority
}

// Swap swaps two parties in the queue.
func (pq *WebsterPriorityQueue) Swap(i, j int) {
	pq.parties[i], pq.parties[j] = pq.parties[j], pq.parties[i]
}

// Push adds a new party to the queue.
func (pq *WebsterPriorityQueue) Push(x interface{}) {
	pq.parties = append(pq.parties, x.(Party))
}

// Pop removes and returns the party with the highest priority.
func (pq *WebsterPriorityQueue) Pop() interface{} {
	old := pq.parties
	n := len(old)
	item := old[n-1]
	pq.parties = old[0 : n-1]
	return item
}

// AllocateWebsterSeats allocates new seats using the Webster method.
// newSeats: number of new seats to allocate across all parties.
// parties: slice of Party with Name and Votes set, Seats may already be assigned.
func AllocateWebsterSeats(newSeats int32, partyCandidates []Party) []Party {
	// Initialize queue with all parties, preserve existing Seats values
	pq := WebsterPriorityQueue{
		parties: make([]Party, len(partyCandidates)),
	}
	nameToIndex := make(map[string]int, len(partyCandidates))
	for i, p := range partyCandidates {
		pq.parties[i] = p
		nameToIndex[p.Name] = i
	}
	heap.Init(&pq)

	remaining := newSeats
	if remaining <= 0 {
		return partyCandidates
	}

	for s := int32(0); s < remaining; s++ {
		// Pop the party with the highest priority
		top := heap.Pop(&pq).(Party)
		top.Seats++
		heap.Push(&pq, top)
	}

	// Collect results in the same order as input
	result := make([]Party, len(partyCandidates))
	for _, p := range pq.parties {
		if idx, ok := nameToIndex[p.Name]; ok {
			result[idx] = p
		}
	}
	return result
}
