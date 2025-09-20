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
	Parties []Party
}

// Check if our WebsterPriorityQueue implements necessary interfaces
var _ heap.Interface = &WebsterPriorityQueue{}

// Len returns the number of Parties in the queue.
func (pq *WebsterPriorityQueue) Len() int {
	return len(pq.Parties)
}

// Less compares two Parties by their Webster priority.
// The party with the higher priority comes first.
func (pq *WebsterPriorityQueue) Less(i, j int) bool {
	// In the Webster method, compare the priority of two Parties:
	// the one with the higher value of Votes/(2*Seats+1) gets the next seat.
	iPriority := pq.Parties[i].Votes / (int64(2*pq.Parties[i].Seats + 1))
	jPriority := pq.Parties[j].Votes / (int64(2*pq.Parties[j].Seats + 1))
	if iPriority == jPriority {
		return (*pq).Parties[i].Votes >= (*pq).Parties[j].Votes // TODO: replace with tie-breaker
	}
	return iPriority > jPriority
}

// Swap swaps two Parties in the queue.
func (pq *WebsterPriorityQueue) Swap(i, j int) {
	pq.Parties[i], pq.Parties[j] = pq.Parties[j], pq.Parties[i]
}

// Push adds a new party to the queue.
func (pq *WebsterPriorityQueue) Push(x interface{}) {
	pq.Parties = append(pq.Parties, x.(Party))
}

// Pop removes and returns the party with the highest priority.
func (pq *WebsterPriorityQueue) Pop() interface{} {
	old := pq.Parties
	n := len(old)
	item := old[n-1]
	pq.Parties = old[0 : n-1]
	return item
}

// AllocateWebsterSeats allocates new seats using the Webster method.
// newSeats: number of new seats to allocate across all Parties.
// Parties: slice of Party with Name and Votes set, Seats may already be assigned.
func AllocateWebsterSeats(newSeats int32, partyCandidates []Party) []Party {
	// Initialize queue with all Parties, preserve existing Seats values
	pq := WebsterPriorityQueue{
		Parties: make([]Party, len(partyCandidates)),
	}
	nameToIndex := make(map[string]int, len(partyCandidates))
	for i, p := range partyCandidates {
		pq.Parties[i] = p
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
	for _, p := range pq.Parties {
		if idx, ok := nameToIndex[p.Name]; ok {
			result[idx] = p
		}
	}
	return result
}
