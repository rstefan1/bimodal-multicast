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
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/httpmessage"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/round"
	"github.com/rstefan1/bimodal-multicast/pkg/peer"
)

type Gossip struct {
	// buffer with addresses of nodes in system
	peerBuffer []peer.Peer
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
func (g *Gossip) randomlySelectPeer() peer.Peer {
	for {
		r := rand.Intn(len(g.peerBuffer))
		if g.selectedPeers[r] {
			continue
		}
		g.selectedPeers[r] = true
		return (g.peerBuffer)[r]
	}
}

// resetSelectedPeers is a helper func that clear slice with selected peers in gossip round
func (g *Gossip) resetSelectedPeers() {
	for i := range g.selectedPeers {
		g.selectedPeers[i] = false
	}
}

// gossipRound is the gossip round that runs every 100ms
func (g *Gossip) gossipRound(stop <-chan struct{}) {
	var (
		dest peer.Peer
		path string
	)

	for {
		select {
		case <-stop:
			g.logger.Printf("End of gossip round from %s:%s", g.gossipAddr, g.gossipPort)
			return
		default:
			g.gossipRoundNumber.Increment()

			gossipMsg := httpmessage.HTTPGossip{
				Addr:        g.gossipAddr,
				Port:        g.gossipPort,
				RoundNumber: g.gossipRoundNumber,
				Digests:     *(g.msgBuffer).DigestBuffer(),
			}

			// gossipLen is number of nodes which will receive gossip message
			gossipLen := int(g.beta*float64(len(g.peerBuffer))/float64(g.gossipRoundNumber.GetNumber())) + 1

			// send gossip messages
			for i := 0; i < gossipLen; i++ {
				dest = g.randomlySelectPeer()
				path = fmt.Sprintf("http://%s:%s/gossip", dest.Addr, dest.Port)
				jsonGossip, err := json.Marshal(gossipMsg)
				if err != nil {
					g.logger.Printf("Gossiper from %s:%s can not marshal the gossip message: %s", g.gossipAddr, g.gossipPort, err)
					continue
				}

				// send the gossip message
				_, err = http.Post(path, "json", bytes.NewBuffer(jsonGossip))
				if err != nil {
					g.logger.Printf("Gossiper from %s:%s can not marshal the gossip message: %s", g.gossipAddr, g.gossipPort, err)
				}
			}

			(*g.msgBuffer).IncrementGossipCount()
			g.resetSelectedPeers()

			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (g *Gossip) Start(stop <-chan struct{}) {
	g.logger.Printf("Starting Gossiper on %s:%s", g.gossipAddr, g.gossipPort)
	g.gossipRound(stop)
}

func New(cfg Config) *Gossip {
	if cfg.Logger == nil {
		cfg.Logger = log.New(os.Stdout, "", 0)
	}

	return &Gossip{
		peerBuffer:        cfg.PeerBuf,
		msgBuffer:         cfg.MsgBuf,
		gossipAddr:        cfg.Addr,
		gossipPort:        cfg.Port,
		selectedPeers:     make([]bool, len(cfg.PeerBuf)),
		beta:              cfg.Beta,
		gossipRoundNumber: cfg.GossipRound,
		logger:            cfg.Logger,
	}
}
