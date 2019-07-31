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
	"net/http"
	"os"
	"time"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/callback"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/round"
)

const (
	defaultBeta = 0.3
	// NOCALLBACK is callback type for messages without callback
	NOCALLBACK = callback.NOCALLBACK
)

// Config is the config for the protocol
type Config struct {
	// Address is HTTP address for node which runs http servers
	// Required
	Address string
	// Port is HTTP port for node which runs http servers
	// Required
	Port string
	// Beta is the expected fanout for gossip rounds
	// Optional
	Beta float64
	// Logger
	// Optional
	Logger *log.Logger
	// Callbacks funtions
	// Optional
	Callbacks map[string]func(interface{}, *log.Logger) error
	// Gossip round duration
	// Optional
	RoundDuration time.Duration
}

// BMMC is the protocol
type BMMC struct {
	// protocol config
	config *Config
	// shared buffer with peers
	peerBuffer *peer.Buffer
	// shared buffer with gossip messages
	messageBuffer *buffer.MessageBuffer
	// gossip round number
	gossipRound *round.GossipRound
	// http server
	server *http.Server
	// custom callback registry
	customCallbacks *callback.CustomRegistry
	// default callback registry
	defaultCallbacks *callback.DefaultRegistry
	// stop channel
	stop chan struct{}

	// TODO remove the following field
	selectedPeers []bool
}

// validateConfig validates given config.
// Also, it sets default values for optional fields.
func validateConfig(cfg *Config) error {
	if len(cfg.Address) == 0 {
		return fmt.Errorf("Address must not be empty")
	}
	if len(cfg.Port) == 0 {
		return fmt.Errorf("Port must not be empty")
	}
	if cfg.Beta == 0 {
		cfg.Beta = defaultBeta
	}
	if cfg.Logger == nil {
		cfg.Logger = log.New(os.Stdout, "", 0)
	}
	if cfg.RoundDuration == 0 {
		cfg.RoundDuration = time.Millisecond * 100
	}

	return nil
}

// New creates a new instance for the protocol
func New(cfg *Config) (*BMMC, error) {
	err := validateConfig(cfg)
	if err != nil {
		return nil, err
	}

	// set callbacks
	callbacks := cfg.Callbacks
	if callbacks == nil {
		callbacks = map[string]func(interface{}, *log.Logger) error{}
	}
	cbCustomRegistry, err := callback.NewCustomRegistry(callbacks)
	if err != nil {
		return nil, fmt.Errorf("Error at creating new custom callbacks registry: %s", err)
	}
	cbDefaultRegistry, err := callback.NewDefaultRegistry()
	if err != nil {
		return nil, fmt.Errorf("Error at creating new default callbacks registry: %s", err)
	}

	b := &BMMC{
		config:           cfg,
		peerBuffer:       peer.NewPeerBuffer(),
		messageBuffer:    buffer.NewMessageBuffer(),
		gossipRound:      round.NewGossipRound(),
		customCallbacks:  cbCustomRegistry,
		defaultCallbacks: cbDefaultRegistry,

		// TODO remove the following line
		selectedPeers: make([]bool, peer.MAXPEERS),
	}
	b.server = b.newServer()

	return b, nil
}

// Start starts the gossip server and the http server
func (b *BMMC) Start() error {
	b.stop = make(chan struct{})

	// start http server
	if err := b.startServer(b.stop); err != nil {
		return err
	}

	// start gossiper
	go func() {
		b.startGossiper(b.stop)
	}()

	return nil
}

// Stop stops the gossip server and the http server
func (b *BMMC) Stop() {
	close(b.stop)
}

// AddMessage adds new message in messages buffer.
func (b *BMMC) AddMessage(msg interface{}, callbackType string) error {
	m := buffer.NewMessage(msg, callbackType)

	err := b.messageBuffer.AddMessage(m)
	if err != nil {
		b.config.Logger.Printf("BMMC %s:%s error at syncing buffer with message %s in round %d: %s",
			b.config.Address, b.config.Port, m.ID, b.gossipRound.GetNumber(), err)
		return err
	}

	b.config.Logger.Printf("BMMC %s:%s synced buffer with message %s in round %d",
		b.config.Address, b.config.Port, m.ID, b.gossipRound.GetNumber())

	// run callback function for messages with a callback registered
	if callbackType != callback.NOCALLBACK {
		err = b.defaultCallbacks.RunDefaultCallbacks(m, b.peerBuffer, b.config.Logger)
		if err != nil {
			b.config.Logger.Printf("Error at calling default callback at %s:%s for message %s in round %d",
				b.config.Address, b.config.Port, m.ID, b.gossipRound.GetNumber())
		}

		err = b.customCallbacks.RunCustomCallbacks(m, b.config.Logger)
		if err != nil {
			b.config.Logger.Printf("Error at calling custom callback at %s:%s for message %s in round %d",
				b.config.Address, b.config.Port, m.ID, b.gossipRound.GetNumber())
		}
	}
	return nil
}

// AddPeer adds new peer in peers buffer
func (b *BMMC) AddPeer(addr, port string) error {
	err := b.peerBuffer.AddPeer(
		peer.NewPeer(addr, port),
	)
	if err != nil {
		return fmt.Errorf("Error at adding the peer %s/%s: %s", addr, port, err)
	}

	err = b.messageBuffer.AddMessage(
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
func (b *BMMC) RemovePeer(addr, port string) error {
	b.peerBuffer.RemovePeer(
		peer.NewPeer(addr, port),
	)

	err := b.messageBuffer.AddMessage(
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
func (b *BMMC) GetMessages() []interface{} {
	return b.messageBuffer.UnwrapMessageBuffer()
}
