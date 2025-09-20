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
	"fmt"
	"testing"
)

func ExampleAllocateWebsterSeats() {
	parties := []Party{
		{Name: "Alpha", Votes: 1200},
		{Name: "Beta", Votes: 900},
		{Name: "Gamma", Votes: 400},
	}
	totalSeats := int32(4)
	result := AllocateWebsterSeats(totalSeats, parties)
	for _, p := range result {
		fmt.Printf("%s: %d seats\n", p.Name, p.Seats)
	}
	// Output:
	// Alpha: 2 seats
	// Beta: 1 seats
	// Gamma: 1 seats
}

// This test checks the Webster (Sainte-Laguë) seat allocation for the classic example from
// https://en.wikipedia.org/wiki/Sainte-Lagu%C3%AB_method, where 230,000 voters allocate 8 seats
// among 4 Parties (PartyA: 100,000 votes, PartyB: 80,000 votes, PartyC: 30,000 votes, PartyD: 20,000 votes).
// The expected seat distribution after each round matches the step-by-step allocation shown in the Wikipedia article.
func TestAllocateWebsterSeats(t *testing.T) {
	parties := []Party{
		{Name: "PartyA", Votes: 100000},
		{Name: "PartyB", Votes: 80000},
		{Name: "PartyC", Votes: 30000},
		{Name: "PartyD", Votes: 20000},
	}
	totalSeats := int32(8)

	// Track seat assignments after each round
	type roundResult struct {
		round int
		seats []int32
	}
	var rounds []roundResult

	// Copy of Parties for mutation
	current := make([]Party, len(parties))
	copy(current, parties)

	for s := int32(1); s <= totalSeats; s++ {
		// Allocate up to s seats
		allocated := AllocateWebsterSeats(s, parties)
		// Record the seat count for each party in order
		seats := make([]int32, len(allocated))
		for i, p := range allocated {
			seats[i] = p.Seats
		}
		rounds = append(rounds, roundResult{round: int(s), seats: seats})
	}

	// Expected seat assignments after each round
	expected := [][]int32{
		// PartyA, PartyB, PartyC, PartyD
		{1, 0, 0, 0}, // 1 seat: PartyA
		{1, 1, 0, 0}, // 2 seats: PartyB
		{2, 1, 0, 0}, // 3 seats: PartyA
		{2, 1, 1, 0}, // 4 seats: PartyC
		{2, 2, 1, 0}, // 5 seats: PartyB
		{3, 2, 1, 0}, // 6 seats: PartyA
		{3, 2, 1, 1}, // 7 seats: PartyD
		{3, 3, 1, 1}, // 8 seats: PartyB
	}

	for i, r := range rounds {
		if len(r.seats) != len(expected[i]) {
			t.Fatalf("round %d: seat count mismatch: got %v, want %v", r.round, r.seats, expected[i])
		}
		for j := range r.seats {
			if r.seats[j] != expected[i][j] {
				t.Errorf("round %d, party %d: got %d seats, want %d", r.round, j, r.seats[j], expected[i][j])
			}
		}
	}
}
