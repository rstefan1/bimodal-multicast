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
	"fmt"
	"log"
	"os"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/callback"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/gossip"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/round"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/server"
)

const defaultBeta = 0.3

type Bmmc struct {
	// shared buffer with addresses of nodes in system
	peerBuffer *peer.PeerBuffer
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
func New(cfg *Config) (*Bmmc, error) {
	if len(cfg.Addr) == 0 {
		return nil, fmt.Errorf("Address must not be empty")
	}
	if len(cfg.Port) == 0 {
		return nil, fmt.Errorf("Port must not be empty")
	}
	if cfg.Beta == 0 {
		cfg.Beta = defaultBeta
	}
	if cfg.Logger == nil {
		cfg.Logger = log.New(os.Stdout, "", 0)
	}

	cbCustomRegistry, err := callback.NewCustomRegistry(cfg.Callbacks)
	if err != nil {
		return nil, fmt.Errorf("Error at creating new custom callbacks registry: %s", err)
	}

	peerBuf := peer.NewPeerBuffer()
	for _, p := range cfg.Peers {
		pp := peer.NewPeer(p.Addr, p.Port)
		_ = peerBuf.AddPeer(pp)
	}

	p := &Bmmc{
		peerBuffer:  peerBuf,
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
		Callbacks:   cbCustomRegistry,
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

	return p, nil
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

func (b *Bmmc) AddMessage(msg, callbackType string) {
	b.msgBuffer.AddMessage(buffer.NewMessage(msg, callbackType))
}

func (b *Bmmc) GetMessages() []string {
	return b.msgBuffer.UnwrapMessageBuffer()
}
