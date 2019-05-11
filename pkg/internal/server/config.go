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

package server

import (
	"log"

	"github.com/rstefan1/bimodal-multicast/pkg/callback"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/round"
	"github.com/rstefan1/bimodal-multicast/pkg/peer"
)

// Config has configs for http server
type Config struct {
	// Addr is http server address
	Addr string
	// Port is http server port
	Port string
	// PeerBuf is the list of peers
	PeerBuf []peer.Peer
	// MsgBuf is the list of messages
	MsgBuf *buffer.MessageBuffer
	// GossipRound is the gossip round number
	GossipRound *round.GossipRound
	// Logger
	Logger *log.Logger
	// Loss is the loss for gossip handler
	Loss float64
	// CallbackRegistry
	Callbacks *callback.CallbacksRegistry
}
