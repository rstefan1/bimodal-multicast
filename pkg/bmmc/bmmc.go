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
	"log"
	"maps"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/callback"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

const (
	// NOCALLBACK is callback type for messages without callback.
	NOCALLBACK = callback.NOCALLBACK

	addPeerErrFmt    = "error at adding the peer %s: %w"
	removePeerErrFmt = "error at removing the peer %s: %w"

	runCallbackErrFmt = "error at calling callback at %s for message %s in round %d"

	createCallbacksRegistryErrFmt = "error at creating callbacks registry: %w"
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
	// callbacks registry
	callbacksRegistry *callback.Registry
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
	callbacksRegistry, err := callback.NewRegistry(cfg.Callbacks)
	if err != nil {
		return nil, fmt.Errorf(createCallbacksRegistryErrFmt, err)
	}

	// create an instance of the protocol
	//nolint: exhaustivestruct
	b := &BMMC{
		config:            cfg,
		peerBuffer:        peer.NewPeerBuffer(),
		messageBuffer:     buffer.NewBuffer(cfg.BufferSize),
		gossipRound:       NewGossipRound(),
		callbacksRegistry: callbacksRegistry,

		// TODO remove the following line
		selectedPeers: make([]bool, peer.MAXPEERS),
	}

	// add internal callbacks
	internalCallbacks := map[string]func(any, *log.Logger) error{
		callback.ADDPEER:    callback.AddPeerCallback,
		callback.REMOVEPEER: callback.RemovePeerCallback,
	}
	maps.Copy(b.callbacksRegistry.Callbacks, internalCallbacks) // maps.Copy(dst, src)

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

func (b *BMMC) runCallbacks(el buffer.Element) {
	if el.CallbackType == callback.NOCALLBACK {
		return
	}

	callbackFn := b.callbacksRegistry.GetCallback(el.CallbackType)
	if callbackFn == nil {
		return
	}

	var callbackData any

	if el.CallbackType == callback.ADDPEER || el.CallbackType == callback.REMOVEPEER {
		// internal callback
		callbackData = callback.PeerCallbackData{
			Element: el,
			Buffer:  b.peerBuffer,
		}
	} else {
		callbackData = el
	}

	if err := callbackFn(callbackData, b.config.Logger); err != nil {
		b.config.Logger.Printf(runCallbackErrFmt, b.config.Host.String(), el.ID, b.gossipRound.GetNumber())
	}
}
