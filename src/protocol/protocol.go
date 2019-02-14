package protocol

import (
	"math/rand"
	"time"

	"github.com/rstefan1/bimodal-multicast/src/httpserver/send"
	. "github.com/rstefan1/bimodal-multicast/src/internal"
)

var (
	// buffer with addresses of nodes in system
	nodeBuffer *[]Node
	// buffer with gossip messages
	msgBuffer *[]Message
	// round numer
	roundNumber = 1
	// beta is the expected fanout for gossip
	beta = 0.5
)

// TODO implement a func that read a yml file with all node (addr & port)

func randomlySelectedNumber() int {
	// TODO check if node wasn't selected before
	return rand.Intn(len(*nodeBuffer))
}

// GossipRound is the gossip round that runs every 100ms in out implementation
func gossipRound() {
	for {
		// increment the round number
		roundNumber++

		gMsg := GossipMessage{
			RoundNumber: roundNumber,
			Digest:      Digest(*msgBuffer),
		}

		length := int(beta * float64(len(*nodeBuffer)) / float64(roundNumber))
		for i := 0; i < length; i++ {
			dest := *nodeBuffer[randomlySelectedNumber()]
			send.Gossip(dest, gMsg)
		}

		for _, m := range *msgBuffer {
			m.GossipCount++
		}

		time.Sleep(100 * time.Millisecond)
	}

	// TODO discard messages for which gossip_count
	// exceeds G, the garbage-collection limit
}

func Start(nodeBuf *[]Node, msgBuf *[]Message) {
	nodeBuffer = nodeBuf
	msgBuffer = msgBuf
	gossipRound()
}
