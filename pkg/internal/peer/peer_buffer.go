/*
Copyright 2024 Robert Andrei STEFAN

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

package peer

import (
	"fmt"
	"math/rand"
	"sync"
)

// Buffer is the buffer with encoded peers.
type Buffer struct {
	peers []string
	mux   *sync.Mutex
}

// NewPeerBuffer creates a PeerBuffer.
func NewPeerBuffer() *Buffer {
	return &Buffer{
		peers: []string{},
		mux:   &sync.Mutex{},
	}
}

// Length returns length of peers buffer.
func (peerBuffer *Buffer) Length() int {
	peerBuffer.mux.Lock()
	defer peerBuffer.mux.Unlock()

	l := len(peerBuffer.peers)

	return l
}

// alreadyExists return true if the peer already exists in peers buffer.
func (peerBuffer *Buffer) alreadyExists(peer string) bool {
	// Important! Whoever calls this function must LOCK the buffer
	for _, p := range peerBuffer.peers {
		if p == peer {
			return true
		}
	}

	return false
}

// AddPeer adds a peer in peers buffer.
func (peerBuffer *Buffer) AddPeer(peer string) error {
	peerBuffer.mux.Lock()
	defer peerBuffer.mux.Unlock()

	if peerBuffer.alreadyExists(peer) {
		return fmt.Errorf("peer %s already exists in peer buffer", peer) //nolint: goerr113
	}

	peerBuffer.peers = append(peerBuffer.peers, peer)

	return nil
}

// RemovePeer removes a peer from peers buffer.
func (peerBuffer *Buffer) RemovePeer(peer string) {
	peerBuffer.mux.Lock()
	defer peerBuffer.mux.Unlock()

	pos := -1 // out of buffer

	for i, p := range peerBuffer.peers {
		if p == peer {
			pos = i

			break
		}
	}

	if pos >= 0 {
		peerBuffer.peers[pos] = peerBuffer.peers[len(peerBuffer.peers)-1] // Copy last element to index pos.
		peerBuffer.peers = peerBuffer.peers[:len(peerBuffer.peers)-1]     // Truncate slice.
	}
}

// GetPeers returns a list of strings that contains peers.
func (peerBuffer *Buffer) GetPeers() []string {
	peerBuffer.mux.Lock()
	defer peerBuffer.mux.Unlock()

	p := make([]string, len(peerBuffer.peers))

	copy(p, peerBuffer.peers)

	return p
}

// GetRandomPeer returns random peer from peers buffer.
func (peerBuffer *Buffer) GetRandomPeer() string {
	peerBuffer.mux.Lock()
	defer peerBuffer.mux.Unlock()

	p := peerBuffer.peers[rand.Intn(len(peerBuffer.peers))] //nolint: gosec

	return p
}

// GetRandomPeers returns a list of random peers from peers buffer.
func (peerBuffer *Buffer) GetRandomPeers(noPeers int) []string {
	peerBuffer.mux.Lock()
	defer peerBuffer.mux.Unlock()

	selectedPeers := []string{}

	for len(selectedPeers) < noPeers {
		randomPeer := peerBuffer.peers[rand.Intn(len(peerBuffer.peers))] //nolint: gosec

		validPeer := true

		for _, selectedPeer := range selectedPeers {
			if selectedPeer == randomPeer {
				validPeer = false

				break
			}
		}

		if validPeer {
			selectedPeers = append(selectedPeers, randomPeer)
		}
	}

	return selectedPeers
}
