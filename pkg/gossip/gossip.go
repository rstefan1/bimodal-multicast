package gossip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/config"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/httpmessage"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
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
	// round number
	roundNumber int64
	// beta is the expected fanout for gossip
	beta float64
	// selected peers for sending gossip message
	selectedPeers []bool
}

// randomlySelectPeer is a helper func that returns a random peer
func (g Gossip) randomlySelectPeer() peer.Peer {
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
func (g Gossip) resetSelectedPeers() {
	for i := range g.selectedPeers {
		g.selectedPeers[i] = false
	}
}

// gossipRound is the gossip round that runs every 100ms
func (g Gossip) gossipRound(stop <-chan struct{}) {
	var (
		dest peer.Peer
		path string
	)

	for {
		select {
		case <-stop:
			log.Printf("End of gossip round from %s:%s", g.gossipAddr, g.gossipPort)
			return
		default:
			// increment round number
			g.roundNumber++

			gossipMsg := httpmessage.HTTPGossip{
				Addr:        g.gossipAddr,
				Port:        g.gossipPort,
				RoundNumber: g.roundNumber,
				Digests:     *(g.msgBuffer).DigestBuffer(),
			}

			// gossipLen is number of nodes which will receive gossip message
			gossipLen := int(g.beta*float64(len(g.peerBuffer))/float64(g.roundNumber)) + 1

			// send gossip messages
			for i := 0; i < gossipLen; i++ {
				dest = g.randomlySelectPeer()
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

			(*g.msgBuffer).IncrementGossipCount()
			g.resetSelectedPeers()

			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (g Gossip) Start(stop <-chan struct{}) {
	log.Printf("Starting Gossip server on %s:%s", g.gossipAddr, g.gossipPort)
	g.gossipRound(stop)
}

func New(cfg config.GossipConfig) Gossip {
	return Gossip{
		peerBuffer:    cfg.PeerBuf,
		msgBuffer:     cfg.MsgBuf,
		gossipAddr:    cfg.Addr,
		gossipPort:    cfg.Port,
		selectedPeers: make([]bool, len(cfg.PeerBuf)),
		roundNumber:   1,
	}
}
