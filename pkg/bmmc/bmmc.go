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

const (
	defaultBeta = 0.3
	// NOCALLBACK is callback type for messages without callback
	NOCALLBACK = callback.NOCALLBACK
)

// Bmmc is the protocol
type Bmmc struct {
	// shared buffer with addresses of nodes in system
	peerBuffer *peer.Buffer
	// shared buffer with gossip messages
	msgBuffer *buffer.MessageBuffer
	// gossip round number
	gossipRound *round.GossipRound
	// http server
	server *server.Server
	// gossiper
	gossiper *gossip.Gossiper
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

	callbacks := cfg.Callbacks
	if callbacks == nil {
		callbacks = map[string]func(interface{}, *log.Logger) (bool, error){}
	}
	cbCustomRegistry, err := callback.NewCustomRegistry(callbacks)
	if err != nil {
		return nil, fmt.Errorf("Error at creating new custom callbacks registry: %s", err)
	}
	cbDefaultRegistry, err := callback.NewDefaultRegistry()
	if err != nil {
		return nil, fmt.Errorf("Error at creating new default callbacks registry: %s", err)
	}

	p := &Bmmc{
		peerBuffer:  peer.NewPeerBuffer(),
		msgBuffer:   buffer.NewMessageBuffer(),
		gossipRound: round.NewGossipRound(),
	}

	p.server = server.New(server.Config{
		Addr:             cfg.Addr,
		Port:             cfg.Port,
		PeerBuf:          p.peerBuffer,
		MsgBuf:           p.msgBuffer,
		GossipRound:      p.gossipRound,
		Logger:           cfg.Logger,
		CustomCallbacks:  cbCustomRegistry,
		DefaultCallbacks: cbDefaultRegistry,
	})

	p.gossiper = gossip.New(gossip.Config{
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
	if err := b.server.Start(b.stop); err != nil {
		return err
	}

	// start gossiper
	go func() {
		b.gossiper.Start(b.stop)
	}()

	return nil
}

// Stop stops the gossip server and the http server
func (b *Bmmc) Stop() {
	close(b.stop)
}

// AddMessage adds new message in messages buffer
func (b *Bmmc) AddMessage(msg interface{}, callbackType string) error {
	return b.msgBuffer.AddMessage(buffer.NewMessage(msg, callbackType))
}

// AddPeer adds new peer in peers buffer
func (b *Bmmc) AddPeer(addr, port string) error {
	err := b.peerBuffer.AddPeer(
		peer.NewPeer(addr, port),
	)
	if err != nil {
		return fmt.Errorf("Error at adding the peer %s/%s: %s", addr, port, err)
	}

	err = b.msgBuffer.AddMessage(
		buffer.NewMessage(
			fmt.Sprintf("%s/%s", addr, port),
			callback.ADDPEER,
		),
	)
	if err != nil {
		return fmt.Errorf("Error at adding the peer %s/%s: %s", addr, port, err)
	}

	return nil
}

// RemovePeer removes given peer from peers buffer
func (b *Bmmc) RemovePeer(addr, port string) error {
	b.peerBuffer.RemovePeer(
		peer.NewPeer(addr, port),
	)

	err := b.msgBuffer.AddMessage(
		buffer.NewMessage(
			fmt.Sprintf("%s/%s", addr, port),
			callback.REMOVEPEER,
		),
	)
	if err != nil {
		return fmt.Errorf("Error at removing the peer %s/%s: %s", addr, port, err)
	}

	return nil
}

// GetMessages returns a slice with all messages from messages buffer
func (b *Bmmc) GetMessages() []interface{} {
	return b.msgBuffer.UnwrapMessageBuffer()
}
