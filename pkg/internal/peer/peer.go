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

package peer

import (
	"fmt"
	"math/rand"
	"sync"
)

// MAXPEERS is the maximum number of peers in buffer
const MAXPEERS = 4096

// Peer is a peer
type Peer struct {
	addr string
	port string
}

// Buffer is the buffer with peers
type Buffer struct {
	peers []Peer
	mux   *sync.Mutex
}

// NewPeer creates a Peer
func NewPeer(addr, port string) Peer {
	return Peer{
		addr: addr,
		port: port,
	}
}

// NewPeerBuffer creates a PeerBuffer
func NewPeerBuffer() *Buffer {
	return &Buffer{
		peers: []Peer{},
		mux:   &sync.Mutex{},
	}
}

// Length returns length of peers buffer
func (peerBuffer *Buffer) Length() int {
	peerBuffer.mux.Lock()
	defer peerBuffer.mux.Unlock()

	l := len(peerBuffer.peers)
	return l
}

// alreadyExists return true if the peer already exists in peers buffer
func (peerBuffer *Buffer) alreadyExists(peer Peer) bool {
	// Important! Whoever calls this function must LOCK the buffer
	for _, p := range peerBuffer.peers {
		if p.addr == peer.addr && p.port == peer.port {
			return true
		}
	}
	return false
}

// AddPeer adds a peer in peers buffer
func (peerBuffer *Buffer) AddPeer(peer Peer) error {
	peerBuffer.mux.Lock()
	defer peerBuffer.mux.Unlock()

	if len(peerBuffer.peers)+1 >= MAXPEERS {
		return fmt.Errorf("the buffer is full. Can add up to %d peers", MAXPEERS)
	}

	if peerBuffer.alreadyExists(peer) {
		return fmt.Errorf("peer %s/%s already exists in peer buffer", peer.addr, peer.port)
	}

	peerBuffer.peers = append(peerBuffer.peers, peer)
	return nil
}

// RemovePeer removes a peer from peers buffer
func (peerBuffer *Buffer) RemovePeer(peer Peer) {
	peerBuffer.mux.Lock()
	defer peerBuffer.mux.Unlock()

	pos := -1

	for i, p := range peerBuffer.peers {
		if p.addr == peer.addr && p.port == peer.port {
			pos = i
			break
		}
	}

	if pos >= 0 {
		peerBuffer.peers[pos] = peerBuffer.peers[len(peerBuffer.peers)-1] // Copy last element to index pos.
		peerBuffer.peers[len(peerBuffer.peers)-1] = Peer{}                // Erase last element (write zero value).
		peerBuffer.peers = peerBuffer.peers[:len(peerBuffer.peers)-1]     // Truncate slice.
	}
}

// GetPeers returns a list of strings that contains peers
func (peerBuffer *Buffer) GetPeers() []string {
	peerBuffer.mux.Lock()
	defer peerBuffer.mux.Unlock()

	p := make([]string, len(peerBuffer.peers))
	for i := range peerBuffer.peers {
		p[i] = fmt.Sprintf("%s/%s", peerBuffer.peers[i].addr, peerBuffer.peers[i].port)
	}

	return p
}

// GetRandom returns random peer from peers buffer
func (peerBuffer *Buffer) GetRandom() (string, string, int) {
	peerBuffer.mux.Lock()
	defer peerBuffer.mux.Unlock()

	r := rand.Intn(len(peerBuffer.peers))
	return peerBuffer.peers[r].addr, peerBuffer.peers[r].port, r
}
