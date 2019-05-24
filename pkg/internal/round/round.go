/*
Copyright 2019 Robert Andrei STEFAN

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

package round

import "sync"

// GossipRound is the number of gossiper rounds
type GossipRound struct {
	Number int64       `json:"gossip_round_number"`
	Mux    *sync.Mutex `json:"gossip_round_mux"`
}

// NewGossipRound creates new GossipRound
func NewGossipRound() *GossipRound {
	return &GossipRound{
		Number: int64(0),
		Mux:    &sync.Mutex{},
	}
}

// Increment increments the gossip round numbers
func (r *GossipRound) Increment() {
	r.Mux.Lock()
	defer r.Mux.Unlock()

	r.Number++
}

// GetNumber returns the gossip round numbers
func (r *GossipRound) GetNumber() int64 {
	r.Mux.Lock()
	defer r.Mux.Unlock()

	n := r.Number
	return n
}
