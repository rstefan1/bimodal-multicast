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
	"net/http"
	"time"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/callback"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

const (
	// NOCALLBACK is callback type for messages without callback
	NOCALLBACK = callback.NOCALLBACK

	addPeerErrFmt    = "error at adding the peer %s/%s: %w"
	removePeerErrFmt = "error at removing the peer %s/%s: %w"

	runDefaultCallbackErrFmt = "error at calling default callback at %s:%s for message %s in round %d"
	runCustomCallbackErrFmt  = "error at calling custom callback at %s:%s for message %s in round %d"

	createCustomCRErrFmt  = "error at creating new custom callbacks registry: %w"
	createDefaultCRErrFmt = "error at creating new default callbacks registry: %w"

	// netClientTimeout is the timeout for http client
	netClientTimeout = time.Second * 10
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
	// http server
	server *http.Server
	// custom callback registry
	customCallbacks *callback.CustomRegistry
	// default callback registry
	defaultCallbacks *callback.DefaultRegistry
	// stop channel
	stop chan struct{}
	// netClient is the http client
	netClient *http.Client

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
	b := &BMMC{
		config:           cfg,
		peerBuffer:       peer.NewPeerBuffer(),
		messageBuffer:    buffer.NewBuffer(cfg.BufferSize),
		gossipRound:      NewGossipRound(),
		customCallbacks:  cbCustomRegistry,
		defaultCallbacks: cbDefaultRegistry,
		netClient: &http.Client{
			Timeout: netClientTimeout,
		},

		// TODO remove the following line
		selectedPeers: make([]bool, peer.MAXPEERS),
	}

	b.server = b.newServer()

	return b, nil
}

// Start starts the gossip server and the http server.
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

// Stop stops the gossip server and the http server.
func (b *BMMC) Stop() {
	close(b.stop)
}

// AddMessage adds new message in messages buffer.
func (b *BMMC) AddMessage(msg interface{}, callbackType string) error {
	m, err := buffer.NewElement(msg, callbackType)
	if err != nil {
		b.config.Logger.Printf(syncBufferLogErrFmt, b.config.Addr, b.config.Port, m.ID, b.gossipRound.GetNumber(), err)
		return err
	}

	if err := b.messageBuffer.Add(m); err != nil {
		b.config.Logger.Printf(syncBufferLogErrFmt, b.config.Addr, b.config.Port, m.ID, b.gossipRound.GetNumber(), err)
		return err
	}

	b.config.Logger.Printf(bufferSyncedLogFmt,
		b.config.Addr, b.config.Port, m.ID, b.gossipRound.GetNumber())

	b.runCallbacks(m, b.config.Addr, b.config.Port)

	return nil
}

// AddPeer adds new peer in peers buffer.
func (b *BMMC) AddPeer(addr, port string) error {
	p, err := peer.NewPeer(addr, port)
	if err != nil {
		return fmt.Errorf(addPeerErrFmt, addr, port, err)
	}

	if err = b.peerBuffer.AddPeer(p); err != nil {
		return fmt.Errorf(addPeerErrFmt, addr, port, err)
	}

	msg, err := buffer.NewElement(
		callback.ComposeAddPeerMessage(addr, port),
		callback.ADDPEER,
	)
	if err != nil {
		return fmt.Errorf(addPeerErrFmt, addr, port, err)
	}

	if err = b.messageBuffer.Add(msg); err != nil {
		return fmt.Errorf(addPeerErrFmt, addr, port, err)
	}

	return nil
}

// RemovePeer removes given peer from peers buffer.
func (b *BMMC) RemovePeer(addr, port string) error {
	p, err := peer.NewPeer(addr, port)
	if err != nil {
		return fmt.Errorf(removePeerErrFmt, addr, port, err)
	}

	b.peerBuffer.RemovePeer(p)

	msg, err := buffer.NewElement(
		callback.ComposeRemovePeerMessage(addr, port),
		callback.REMOVEPEER,
	)
	if err != nil {
		return fmt.Errorf(removePeerErrFmt, addr, port, err)
	}

	if err := b.messageBuffer.Add(msg); err != nil {
		return fmt.Errorf(removePeerErrFmt, addr, port, err)
	}

	return nil
}

// GetMessages returns a slice with all messages from messages buffer.
func (b *BMMC) GetMessages() []interface{} {
	return b.messageBuffer.Messages()
}

// GetPeers returns an array with all peers from peers buffer.
func (b *BMMC) GetPeers() []string {
	return b.peerBuffer.GetPeers()
}

func (b *BMMC) runCallbacks(m buffer.Element, hostAddr, hostPort string) {
	// TODO remove hostAddr and hostport from func args. These are used only for logging
	if m.CallbackType != callback.NOCALLBACK {
		if err := b.defaultCallbacks.RunCallbacks(m, b.peerBuffer, b.config.Logger); err != nil {
			b.config.Logger.Printf(runDefaultCallbackErrFmt, hostAddr, hostPort, m.ID, b.gossipRound.GetNumber())
		}

		if err := b.customCallbacks.RunCallbacks(m, b.config.Logger); err != nil {
			b.config.Logger.Printf(runCustomCallbackErrFmt, hostAddr, hostPort, m.ID, b.gossipRound.GetNumber())
		}
	}
}
