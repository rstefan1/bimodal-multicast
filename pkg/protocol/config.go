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

package protocol

import (
	"github.com/rstefan1/bimodal-multicast/pkg/peer"
)

type Config struct {
	// GossipAddr is HTTP address for node which runs gossip round
	GossipAddr string
	// GossipPort is HTTP port for node which runs gossip round
	GossipPort string
	// HTTPAddr is http server address
	HTTPAddr string
	// HTTPPort is http server port
	HTTPPort string
	// PeerBuf is the list of peers
	Peers []peer.Peer
	// Beta is the expected fanout for gossip rounds
	Beta float64
}
