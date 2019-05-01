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
	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/gossip"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/round"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/server"
	"github.com/rstefan1/bimodal-multicast/pkg/peer"
)

type Bmmc struct {
	// shared buffer with addresses of nodes in system
	peerBuffer []peer.Peer
	// shared buffer with gossip messages
	msgBuffer *buffer.MessageBuffer
	// gossip round number
	gossipRound *round.GossipRound
	// http server
	httpServer *server.HTTP
	// gossip server
	gossipServer *gossip.Gossip
	// stop channel
	stop chan struct{}
}

// New creates a new instance for the protocol
func New(cfg Config) *Bmmc {
	p := &Bmmc{
		peerBuffer:  cfg.Peers,
		msgBuffer:   buffer.NewMessageBuffer(),
		gossipRound: round.NewGossipRound(),
	}

	p.httpServer = server.New(server.Config{
		Addr:        cfg.Addr,
		Port:        cfg.Port,
		PeerBuf:     p.peerBuffer,
		MsgBuf:      p.msgBuffer,
		GossipRound: p.gossipRound,
		Logger:      cfg.Logger,
	})

	p.gossipServer = gossip.New(gossip.Config{
		Addr:        cfg.Addr,
		Port:        cfg.Port,
		PeerBuf:     p.peerBuffer,
		MsgBuf:      p.msgBuffer,
		Beta:        cfg.Beta,
		GossipRound: p.gossipRound,
		Logger:      cfg.Logger,
	})

	return p
}

// Start starts the gossip server and the http server
func (b *Bmmc) Start() error {
	b.stop = make(chan struct{})

	// start http server
	if err := b.httpServer.Start(b.stop); err != nil {
		return err
	}

	// start gossip server
	go func() {
		b.gossipServer.Start(b.stop)
	}()

	return nil
}

// Stop stops the gossip server and the http server
func (b *Bmmc) Stop() {
	close(b.stop)
}

func (b *Bmmc) AddMessage(msg string) {
	b.msgBuffer.AddMessage(buffer.NewMessage(msg))
}

func (b *Bmmc) GetMessages() []string {
	return b.msgBuffer.UnwrapMessageBuffer()
}
