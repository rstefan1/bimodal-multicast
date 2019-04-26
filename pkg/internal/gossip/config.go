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
	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/peer"
)

// Config has configs for gossip object
type Config struct {
	// Addr is the address for node which runs gossip round
	Addr string
	// Port is the port for node which runs gossip round
	Port string
	// PeerBuf is the list of peers
	PeerBuf []peer.Peer
	// MsgBuf is the list of messages
	MsgBuf *buffer.MessageBuffer
	// Beta is the expected fanout
	Beta float64
}
