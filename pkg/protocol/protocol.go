package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/rstefan1/bimodal-multicast/src/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/src/internal/httpmessage"
	"github.com/rstefan1/bimodal-multicast/src/internal/peer"
)

var (
	// buffer with addresses of nodes in system
	peerBuffer *[]peer.Peer
	// buffer with gossip messages
	msgBuffer *buffer.MessageBuffer
	// gossipAddr is the gossip server address
	gossipAddr string
	// port is the gossip node port
	gossipPort string
	// round number
	roundNumber int64
	// beta is the expected fanout for gossip
	beta = 0.5
	// selected peers for sending gossip message
	selectedPeers []bool
)

func randomlySelectedPeer() peer.Peer {
	for {
		r := rand.Intn(len(*peerBuffer))
		if selectedPeers[r] {
			continue
		}
		selectedPeers[r] = true
		return (*peerBuffer)[r]
	}
}

// GossipRound is the gossip round that runs every 100ms in out implementation
func gossipRound() {
	var (
		dest peer.Peer
		path string
	)
	for {
		// increment round number
		roundNumber++

		gossipMsg := httpmessage.HTTPGossip{
			Addr:        gossipAddr,
			Port:        gossipPort,
			RoundNumber: roundNumber,
			Digests:     (*msgBuffer).DigestBuffer(),
		}

		// gossipLen is number of nodes which will receive gossip message
		gossipLen := int(beta * float64(len(*peerBuffer)) / float64(roundNumber))

		for i := 0; i < gossipLen; i++ {
			dest = randomlySelectedPeer()
			path = fmt.Sprintf("http://%s:%s/gossip", dest.Addr, dest.Port)
			jsonGossip, err := json.Marshal(gossipMsg)
			if err != nil {
				log.Fatal(err)
				continue
			}

			// send the gossip message
			_, err = http.Post(path, "json", bytes.NewBuffer(jsonGossip))
			if err != nil {
				log.Fatal(err)
			}
		}

		(*msgBuffer).IncrementGossipCount()

		for i := range selectedPeers {
			selectedPeers[i] = false
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func Start(peerBuf *[]peer.Peer, msgBuf *buffer.MessageBuffer, addr, port string) {
	peerBuffer = peerBuf
	msgBuffer = msgBuf
	gossipAddr = addr
	gossipPort = port
	selectedPeers = make([]bool, len(*peerBuffer))
	roundNumber = 1

	gossipRound()
}
