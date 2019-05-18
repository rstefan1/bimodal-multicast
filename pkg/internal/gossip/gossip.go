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

package gossip

import (
	"log"
	"time"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/httputil"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/round"
)

type Gossiper struct {
	// buffer with addresses of nodes in system
	peerBuffer *peer.PeerBuffer
	// buffer with gossip messages
	msgBuffer *buffer.MessageBuffer
	// gossipAddr is the gossip server address
	gossipAddr string
	// port is the gossip node port
	gossipPort string
	// gossip round number
	gossipRoundNumber *round.GossipRound
	// beta is the expected fanout for gossip
	beta float64
	// selected peers for sending gossip message
	selectedPeers []bool
	// logger
	logger *log.Logger
}

// randomlySelectPeer is a helper func that returns a random peer
func (g *Gossiper) randomlySelectPeer() (string, string) {
	for {
		addr, port, i := g.peerBuffer.GetRandom()
		if g.selectedPeers[i] {
			continue
		}
		g.selectedPeers[i] = true
		return addr, port
	}
}

// resetSelectedPeers is a helper func that clear slice with selected peers in gossip round
func (g *Gossiper) resetSelectedPeers() {
	for i := range g.selectedPeers {
		g.selectedPeers[i] = false
	}
}

// gossipRound is the gossip round that runs every 100ms
func (g *Gossiper) gossipRound(stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			g.logger.Printf("End of gossip round from %s:%s", g.gossipAddr, g.gossipPort)
			return
		default:
			g.gossipRoundNumber.Increment()

			// gossipLen is number of nodes which will receive gossip message.
			// It will be 0 if the node has empty peers list.
			var gossipLen int
			if g.peerBuffer.Length() == 0 {
				gossipLen = 0
			} else {
				gossipLen = int(g.beta*float64(g.peerBuffer.Length())/float64(g.gossipRoundNumber.GetNumber())) + 1
			}

			// send gossip messages
			for i := 0; i < gossipLen; i++ {
				destAddr, destPort := g.randomlySelectPeer()

				gossipMsg := httputil.HTTPGossip{
					Addr:        g.gossipAddr,
					Port:        g.gossipPort,
					RoundNumber: g.gossipRoundNumber,
					Digests:     *(g.msgBuffer).DigestBuffer(),
				}

				err := httputil.SendGossip(gossipMsg, destAddr, destPort)
				if err != nil {
					g.logger.Printf("%s", err)
				}
			}

			(*g.msgBuffer).IncrementGossipCount()
			g.resetSelectedPeers()

			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (g *Gossiper) Start(stop <-chan struct{}) {
	g.logger.Printf("Starting Gossiper on %s:%s", g.gossipAddr, g.gossipPort)
	g.gossipRound(stop)
}

func New(cfg Config) *Gossiper {
	return &Gossiper{
		peerBuffer:        cfg.PeerBuf,
		msgBuffer:         cfg.MsgBuf,
		gossipAddr:        cfg.Addr,
		gossipPort:        cfg.Port,
		selectedPeers:     make([]bool, peer.MAXPEERS),
		beta:              cfg.Beta,
		gossipRoundNumber: cfg.GossipRound,
		logger:            cfg.Logger,
	}
}
