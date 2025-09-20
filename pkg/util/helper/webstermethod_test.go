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
	"reflect"
	"testing"
)

func ExampleAllocateWebsterSeats() {
	partiesVotes := map[string]int64{
		"Alpha": 1200,
		"Beta":  900,
		"Gamma": 400,
	}
	totalSeats := int32(4)
	result := AllocateWebsterSeats(totalSeats, partiesVotes, nil, nil)
	for _, p := range result {
		fmt.Printf("%s: %d seats\n", p.Name, p.Seats)
	}
	// Output:
	// Alpha: 2 seats
	// Beta: 1 seats
	// Gamma: 1 seats
}

func TestAllocateWebsterSeats(t *testing.T) {
	// This is used for test the classic example from https://en.wikipedia.org/wiki/Sainte-Lagu%C3%AB_method,
	// where 230,000 voters allocate 8 seats among 4 Parties.
	// The expected seat distribution after each round matches the step-by-step allocation shown in the Wikipedia.
	partiesVotes := map[string]int64{
		"PartyA": 100000,
		"PartyB": 80000,
		"PartyC": 30000,
		"PartyD": 20000,
	}

	tests := []struct {
		name               string
		newSeats           int32
		partyVotes         map[string]int64
		initialAssignments map[string]int32
		tieBreaker         func(a, b Party) bool
		expected           []Party
	}{
		{
			name:               "classic example, round 1",
			newSeats:           1,
			partyVotes:         partiesVotes,
			initialAssignments: nil,
			tieBreaker:         nil,
			expected: []Party{
				{Name: "PartyA", Votes: 100000, Seats: 1},
				{Name: "PartyB", Votes: 80000, Seats: 0},
				{Name: "PartyC", Votes: 30000, Seats: 0},
				{Name: "PartyD", Votes: 20000, Seats: 0},
			},
		},
		{
			name:               "classic example, round 2",
			newSeats:           2,
			partyVotes:         partiesVotes,
			initialAssignments: nil,
			tieBreaker:         nil,
			expected: []Party{
				{Name: "PartyA", Votes: 100000, Seats: 1},
				{Name: "PartyB", Votes: 80000, Seats: 1},
				{Name: "PartyC", Votes: 30000, Seats: 0},
				{Name: "PartyD", Votes: 20000, Seats: 0},
			},
		},
		{
			name:               "classic example, round 3",
			newSeats:           3,
			partyVotes:         partiesVotes,
			initialAssignments: nil,
			tieBreaker:         nil,
			expected: []Party{
				{Name: "PartyA", Votes: 100000, Seats: 2},
				{Name: "PartyB", Votes: 80000, Seats: 1},
				{Name: "PartyC", Votes: 30000, Seats: 0},
				{Name: "PartyD", Votes: 20000, Seats: 0},
			},
		},
		{
			name:               "classic example, round 4",
			newSeats:           4,
			partyVotes:         partiesVotes,
			initialAssignments: nil,
			tieBreaker:         nil,
			expected: []Party{
				{Name: "PartyA", Votes: 100000, Seats: 2},
				{Name: "PartyB", Votes: 80000, Seats: 1},
				{Name: "PartyC", Votes: 30000, Seats: 1},
				{Name: "PartyD", Votes: 20000, Seats: 0},
			},
		},
		{
			name:               "classic example, round 5",
			newSeats:           5,
			partyVotes:         partiesVotes,
			initialAssignments: nil,
			tieBreaker:         nil,
			expected: []Party{
				{Name: "PartyA", Votes: 100000, Seats: 2},
				{Name: "PartyB", Votes: 80000, Seats: 2},
				{Name: "PartyC", Votes: 30000, Seats: 1},
				{Name: "PartyD", Votes: 20000, Seats: 0},
			},
		},
		{
			name:               "classic example, round 6",
			newSeats:           6,
			partyVotes:         partiesVotes,
			initialAssignments: nil,
			tieBreaker:         nil,
			expected: []Party{
				{Name: "PartyA", Votes: 100000, Seats: 3},
				{Name: "PartyB", Votes: 80000, Seats: 2},
				{Name: "PartyC", Votes: 30000, Seats: 1},
				{Name: "PartyD", Votes: 20000, Seats: 0},
			},
		},
		{
			name:               "classic example, round 7",
			newSeats:           7,
			partyVotes:         partiesVotes,
			initialAssignments: nil,
			tieBreaker:         nil,
			expected: []Party{
				{Name: "PartyA", Votes: 100000, Seats: 3},
				{Name: "PartyB", Votes: 80000, Seats: 2},
				{Name: "PartyC", Votes: 30000, Seats: 1},
				{Name: "PartyD", Votes: 20000, Seats: 1},
			},
		},
		{
			name:               "classic example, round 8",
			newSeats:           8,
			partyVotes:         partiesVotes,
			initialAssignments: nil,
			tieBreaker:         nil,
			expected: []Party{
				{Name: "PartyA", Votes: 100000, Seats: 3},
				{Name: "PartyB", Votes: 80000, Seats: 3},
				{Name: "PartyC", Votes: 30000, Seats: 1},
				{Name: "PartyD", Votes: 20000, Seats: 1},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := AllocateWebsterSeats(test.newSeats, test.partyVotes, test.initialAssignments, test.tieBreaker)
			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, got)
			}
		})
	}
}
