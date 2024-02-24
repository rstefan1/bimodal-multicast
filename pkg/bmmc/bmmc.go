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
	"fmt"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/callback"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

const (
	// NOCALLBACK is callback type for messages without callback.
	NOCALLBACK = callback.NOCALLBACK

	addPeerErrFmt    = "error at adding the peer %s: %w"
	removePeerErrFmt = "error at removing the peer %s: %w"

	runDefaultCallbackErrFmt = "error at calling default callback at %s for message %s in round %d"
	runCustomCallbackErrFmt  = "error at calling custom callback at %s for message %s in round %d"

	createCustomCRErrFmt  = "error at creating new custom callbacks registry: %w"
	createDefaultCRErrFmt = "error at creating new default callbacks registry: %w"
)

// BMMC is the bimodal multicast protocol.
type BMMC struct {
	// protocol config
	config *Config
	// shared buffer with peers
	peerBuffer *peer.Buffer
	// shared buffer with gossip messages
	messageBuffer *buffer.Buffer
	// gossip round number
	gossipRound *GossipRound
	// custom callback registry
	customCallbacks *callback.CustomRegistry
	// default callback registry
	defaultCallbacks *callback.DefaultRegistry
	// stop channel
	stop chan struct{}

	// TODO remove the following field
	selectedPeers []bool
}

// New creates a new instance for the protocol.
func New(cfg *Config) (*BMMC, error) {
	// validate given config
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	// fill optional fields of the config
	cfg.fillEmptyFields()

	// set callbacks
	cbCustomRegistry, err := callback.NewCustomRegistry(cfg.Callbacks)
	if err != nil {
		return nil, fmt.Errorf(createCustomCRErrFmt, err)
	}

	cbDefaultRegistry, err := callback.NewDefaultRegistry()
	if err != nil {
		return nil, fmt.Errorf(createDefaultCRErrFmt, err)
	}

	// create an instance of the protocol
	//nolint: exhaustivestruct
	b := &BMMC{
		config:           cfg,
		peerBuffer:       peer.NewPeerBuffer(),
		messageBuffer:    buffer.NewBuffer(cfg.BufferSize),
		gossipRound:      NewGossipRound(),
		customCallbacks:  cbCustomRegistry,
		defaultCallbacks: cbDefaultRegistry,

		// TODO remove the following line
		selectedPeers: make([]bool, peer.MAXPEERS),
	}

	return b, nil
}

// Start starts the gossip server and the http server.
func (b *BMMC) Start() error {
	b.stop = make(chan struct{})

	// start gossiper
	go func() {
		b.startGossiper(b.stop)
	}()

	return nil
}

// Stop stops the gossip server and the http server.
func (b *BMMC) Stop() {
	close(b.stop)
}

// AddMessage adds new message in messages buffer.
func (b *BMMC) AddMessage(msg any, callbackType string) error {
	m, err := buffer.NewElement(msg, callbackType, false)
	if err != nil {
		b.config.Logger.Printf(syncBufferLogErrFmt, b.config.Host.String(), m.ID, b.gossipRound.GetNumber(), err)

		return err //nolint: wrapcheck
	}

	if err := b.messageBuffer.Add(m); err != nil {
		b.config.Logger.Printf(syncBufferLogErrFmt, b.config.Host.String(), m.ID, b.gossipRound.GetNumber(), err)

		return err //nolint: wrapcheck
	}

	b.config.Logger.Printf(bufferSyncedLogFmt,
		b.config.Host.String(), m.ID, b.gossipRound.GetNumber())

	b.runCallbacks(m)

	return nil
}

// AddPeer adds new peer in peers buffer.
func (b *BMMC) AddPeer(p string) error {
	if err := b.peerBuffer.AddPeer(p); err != nil {
		return fmt.Errorf(addPeerErrFmt, p, err)
	}

	msg, err := buffer.NewElement(p, callback.ADDPEER, true)
	if err != nil {
		return fmt.Errorf(addPeerErrFmt, p, err)
	}

	if err = b.messageBuffer.Add(msg); err != nil {
		return fmt.Errorf(addPeerErrFmt, p, err)
	}

	return nil
}

// RemovePeer removes given peer from peers buffer.
func (b *BMMC) RemovePeer(p string) error {
	b.peerBuffer.RemovePeer(p)

	msg, err := buffer.NewElement(p, callback.REMOVEPEER, true)
	if err != nil {
		return fmt.Errorf(removePeerErrFmt, p, err)
	}

	if err := b.messageBuffer.Add(msg); err != nil {
		return fmt.Errorf(removePeerErrFmt, p, err)
	}

	return nil
}

// GetMessages returns a slice with all user messages from messages buffer.
func (b *BMMC) GetMessages() []any {
	return b.messageBuffer.Messages(false)
}

// GetPeers returns an array with all peers from peers buffer.
func (b *BMMC) GetPeers() []string {
	return b.peerBuffer.GetPeers()
}

func (b *BMMC) runCallbacks(m buffer.Element) {
	if m.CallbackType != callback.NOCALLBACK {
		if err := b.defaultCallbacks.RunCallbacks(m, b.peerBuffer, b.config.Logger); err != nil {
			b.config.Logger.Printf(runDefaultCallbackErrFmt, b.config.Host.String(), m.ID, b.gossipRound.GetNumber())
		}

		if err := b.customCallbacks.RunCallbacks(m, b.config.Logger); err != nil {
			b.config.Logger.Printf(runCustomCallbackErrFmt, b.config.Host.String(), m.ID, b.gossipRound.GetNumber())
		}
	}
}
