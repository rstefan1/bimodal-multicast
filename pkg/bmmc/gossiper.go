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
			b.config.Logger.Info("ending gossiper")

			return
		default:
			b.gossipRound.Increment()

			gossipLen := b.computeGossipLen()

			randomlySelectedPeers := b.peerBuffer.GetRandomPeers(gossipLen)

			// send gossip messages
			for _, p := range randomlySelectedPeers {
				gossipMsg := Gossip{
					Host:        b.config.Host.String(),
					RoundNumber: b.gossipRound,
					Digest:      b.messageBuffer.Digest(),
				}

				b.sendGossip(gossipMsg, p) //nolint: errcheck
			}

			(*b.messageBuffer).IncrementGossipCount()

			time.Sleep(b.config.RoundDuration)
		}
	}
}

func (b *BMMC) startGossiper(stop <-chan struct{}) {
	b.config.Logger.Info("starting gossiper")

	b.round(stop)
}
