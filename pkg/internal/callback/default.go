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

package callback

import (
	"errors"
	"log"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

const (
	peerAddedLogFmt   = "peer %s added in the peers buffer"
	peerRemovedLogFmt = "peer %s removed from the peers buffer"
)

var (
	//nolint: gochecknoglobals
	defaultCallbacks = map[string]func(buffer.Element, *peer.Buffer, *log.Logger) error{
		ADDPEER:    addPeerCallback,
		REMOVEPEER: removePeerCallback,
	}

	errCannotConvertToString = errors.New("cannot convert the given message to string")
)

// DefaultRegistry is a default callbacks registry.
type DefaultRegistry struct {
	callbacks map[string]func(buffer.Element, *peer.Buffer, *log.Logger) error
}

// NewDefaultRegistry creates a default callback registry.
func NewDefaultRegistry() (*DefaultRegistry, error) {
	return &DefaultRegistry{
		callbacks: defaultCallbacks,
	}, nil
}

// GetCallback returns a default callback from registry.
func (r *DefaultRegistry) GetCallback(t string) func(buffer.Element, *peer.Buffer, *log.Logger) error {
	if v, ok := r.callbacks[t]; ok {
		return v
	}

	return nil
}

// RunCallbacks runs default callbacks.
func (r *DefaultRegistry) RunCallbacks(m buffer.Element, peerBuf *peer.Buffer, logger *log.Logger) error {
	callbackFn := r.GetCallback(m.CallbackType)
	if callbackFn == nil {
		return nil
	}

	return callbackFn(m, peerBuf, logger)
}

func addPeerCallback(msg buffer.Element, peersBuf *peer.Buffer, logger *log.Logger) error {
	peer, convOk := msg.Msg.(string)
	if !convOk {
		return errCannotConvertToString
	}

	// add peer in buffer
	if err := peersBuf.AddPeer(peer); err != nil {
		return err //nolint: wrapcheck
	}

	logger.Printf(peerAddedLogFmt, peer)

	return nil
}

func removePeerCallback(msg buffer.Element, peersBuf *peer.Buffer, logger *log.Logger) error {
	peer, convOk := msg.Msg.(string)
	if !convOk {
		return errCannotConvertToString
	}

	peersBuf.RemovePeer(peer)

	logger.Printf(peerRemovedLogFmt, peer)

	return nil
}
