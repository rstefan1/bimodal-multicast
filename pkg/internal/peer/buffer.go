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

	if len(peerBuffer.peers)+1 >= MAXPEERS {
		return fmt.Errorf("the buffer is full. Can add up to %d peers", MAXPEERS) //nolint: goerr113
	}

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

// GetRandom returns random peer from peers buffer.
func (peerBuffer *Buffer) GetRandom() (string, int) {
	peerBuffer.mux.Lock()
	defer peerBuffer.mux.Unlock()

	r := rand.Intn(len(peerBuffer.peers)) //nolint: gosec

	return peerBuffer.peers[r], r
}
