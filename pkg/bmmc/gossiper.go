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

package bmmc

import (
	"time"
)

const (
	startGossiperLogFmt = "Starting gossiper for %s:%s"
	stopGossiperLogFmt  = "End of gossip round from %s:%s"
)

// randomlySelectPeer is a helper func that returns a random peer
func (b *BMMC) randomlySelectPeer() (string, string) {
	for {
		addr, port, i := b.peerBuffer.GetRandom()
		if b.selectedPeers[i] {
			continue
		}
		b.selectedPeers[i] = true
		return addr, port
	}
}

// resetSelectedPeers is a helper func that clear slice with selected peers in gossip round
func (b *BMMC) resetSelectedPeers() {
	for i := range b.selectedPeers {
		b.selectedPeers[i] = false
	}
}

// gossipLen is number of nodes which will receive gossip message.
// It will be 0 if the node has empty peers buffer or if the node has
// empty message buffer.
func (b *BMMC) computeGossipLen() int {
	if b.peerBuffer.Length() == 0 || b.messageBuffer.Length() == 0 || b.config.Beta == 0 {
		return 0
	}
	return int(b.config.Beta*float64(b.peerBuffer.Length())) + 1
}

func (b *BMMC) round(stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			b.config.Logger.Printf(stopGossiperLogFmt, b.config.Addr, b.config.Port)
			return
		default:
			b.gossipRound.Increment()

			gossipLen := b.computeGossipLen()

			// send gossip messages
			for i := 0; i < gossipLen; i++ {
				destAddr, destPort := b.randomlySelectPeer()

				gossipMsg := HTTPGossip{
					Addr:        b.config.Addr,
					Port:        b.config.Port,
					RoundNumber: b.gossipRound,
					Digests:     *(b.messageBuffer).DigestBuffer(),
				}

				err := b.sendGossip(gossipMsg, destAddr, destPort)
				if err != nil {
					b.config.Logger.Printf("%s", err)
				}
			}

			(*b.messageBuffer).IncrementGossipCount()
			b.resetSelectedPeers()

			time.Sleep(b.config.RoundDuration)
		}
	}
}

func (b *BMMC) startGossiper(stop <-chan struct{}) {
	b.config.Logger.Printf(startGossiperLogFmt, b.config.Addr, b.config.Port)
	b.round(stop)
}
