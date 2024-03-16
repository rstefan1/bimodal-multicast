/*
Copyright 2024 Robert Andrei STEFAN

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
	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
)

const (
	// GossipRoute is the route for gossip messages.
	GossipRoute = "/gossip"
	// SolicitationRoute is the route for solicitation messages.
	SolicitationRoute = "/solicitation"
	// SynchronizationRoute is the route for synchronization messages.
	SynchronizationRoute = "/synchronization"
)

// GossipHandler handles a gossip message.
func (b *BMMC) GossipHandler(body []byte) {
	gossipDigest, p, roundNumber, err := b.receiveGossip(body)
	if err != nil {
		return
	}

	digest := b.messageBuffer.Digest()
	missingDigest := buffer.MissingStrings(gossipDigest, digest)

	if len(missingDigest) > 0 {
		solicitationMsg := Solicitation{
			Host:        b.config.Host.String(),
			RoundNumber: roundNumber,
			Digest:      missingDigest,
		}

		if err = b.sendSolicitation(solicitationMsg, p); err != nil {
			return
		}
	}
}

// SolicitationHandler handles a solicitation message.
func (b *BMMC) SolicitationHandler(body []byte) {
	missingDigest, p, _, err := b.receiveSolicitation(body)
	if err != nil {
		return
	}

	missingElements := b.messageBuffer.ElementsFromIDs(missingDigest)

	synchronizationMsg := Synchronization{
		Host:     b.config.Host.String(),
		Elements: missingElements,
	}

	if err = b.sendSynchronization(synchronizationMsg, p); err != nil {
		return
	}
}

// SynchronizationHandler handles a synchronization message.
func (b *BMMC) SynchronizationHandler(body []byte) {
	rcvElements, _, err := b.receiveSynchronization(body)
	if err != nil {
		return
	}

	for _, m := range rcvElements {
		err = b.messageBuffer.Add(m)
		if err != nil {
			b.config.Logger.Error("failed to sync buffer with message", "err", err, "msg", m.Msg)
		} else {
			b.config.Logger.Debug("buffer successfully synced with message", "msg", m.Msg)

			b.runCallbacks(m)
		}
	}
}
