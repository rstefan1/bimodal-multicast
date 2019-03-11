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

type Protocol struct {
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
	beta float64
	// selected peers for sending gossip message
	selectedPeers []bool
}

// randomlySelectPeer is a helper func that returns a random peer
func (p Protocol) randomlySelectPeer() peer.Peer {
	for {
		r := rand.Intn(len(*p.peerBuffer))
		if p.selectedPeers[r] {
			continue
		}
		p.selectedPeers[r] = true
		return (*p.peerBuffer)[r]
	}
}

// resetSelectedPeers is a helper func that clear slice with selected peers in gossip round
func (p Protocol) resetSelectedPeers() {
	for i := range p.selectedPeers {
		p.selectedPeers[i] = false
	}
}

// gossipRound is the gossip round that runs every 100ms
func (p Protocol) gossipRound() {
	var (
		dest peer.Peer
		path string
	)
	for {
		// increment round number
		p.roundNumber++

		gossipMsg := httpmessage.HTTPGossip{
			Addr:        p.gossipAddr,
			Port:        p.gossipPort,
			RoundNumber: p.roundNumber,
			Digests:     (*p.msgBuffer).DigestBuffer(),
		}

		// gossipLen is number of nodes which will receive gossip message
		gossipLen := int(p.beta * float64(len(*p.peerBuffer)) / float64(p.roundNumber))

		// send gossip messages
		for i := 0; i < gossipLen; i++ {
			dest = p.randomlySelectPeer()
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

		(*p.msgBuffer).IncrementGossipCount()
		p.resetSelectedPeers()

		time.Sleep(100 * time.Millisecond)
	}
}

func (p Protocol) Start() {
	p.gossipRound()
}

func (p Protocol) NewProtocol(peerBuf *[]peer.Peer, msgBuf *buffer.MessageBuffer, addr, port string) Protocol {
	return Protocol{
		peerBuffer:    peerBuf,
		msgBuffer:     msgBuf,
		gossipAddr:    addr,
		gossipPort:    port,
		selectedPeers: make([]bool, len(*peerBuf)),
		roundNumber:   1,
	}
}
