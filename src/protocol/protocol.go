package pbcast

import (
	"math/rand"
	"time"

	"github.com/rstefan1/bimodal-multicast/src/httpserver/send"
	. "github.com/rstefan1/bimodal-multicast/src/internal"
)

var (
	// buffer with gossip messages
	msgBuffer []Message
	// buffer with addresses of nodes in system
	nodeBuffer []Node
	// round numer
	roundNumber = 1
	// beta is the expected fanout for gossip
	beta = 0.5
)

// TODO implement a func that read a yml file with all node (addr & port)

func randomlySelectedNumber() int {
	// TODO check if node wasn't selected before
	return rand.Intn(len(nodeBuffer))
}

// newGossipMessage is a func that returns a new gossip message
func newGossipMessage() GossipMessage {
	return GossipMessage{
		roundNumber: roundNumber,
		msg:         msgBuffer,
	}
}

// GossipRound is the gossip round that runs every 100ms in out implementation
func GossipRound() {
	for {
		// increment the round number
		roundNumber++

		// create the gossip message
		gMsg := newGossipMessage()

		for i := 0; i < int(beta*float64(len(nodeBuffer))/float64(roundNumber)); i++ {
			// select a random node
			dest := nodeBuffer[randomlySelectedNumber()]
			send.GossipMessage(dest, gMsg)
		}

		for _, m := range msgBuffer {
			m.GossipCount++
		}

		time.Sleep(100 * time.Millisecond)
	}

	// TODO discard messages for which gossip_count
	// exceeds G, the garbage-collection limit
}
