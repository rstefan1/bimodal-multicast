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
	"fmt"
	"log"
	"strings"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

const (
	invalidAddPeerMsgErr         = "invalid add peer message"
	invalidRemovePeerMsgErr      = "invalid remove peer message"
	inexistentDefaultCallbackErr = "callback doesn't exist in the default registry"

	peerAddedLogFmt   = "peer %s/%s added in the peers buffer"
	peerRemovedLogFmt = "peer %s/%s removed from the peers buffer"

	// addPrefix is the prefix for add peer messages
	addPrefix = "add"
	// removePrefix is the prefix for remove peer messages
	removePrefix = "remove"
)

// ComposeAddPeerMessage returns a `add peer` message with given addr and port
func ComposeAddPeerMessage(addr, port string) string {
	return fmt.Sprintf("%s/%s/%s", addPrefix, addr, port)
}

// DecomposeAddPeerMessage decomposes given `add peer` message to addr and port
func DecomposeAddPeerMessage(msg string) (string, string, error) {
	host := strings.Split(msg, "/")
	if len(host) != 3 {
		return "", "", errors.New(invalidAddPeerMsgErr)
	}

	if host[0] != addPrefix {
		return "", "", errors.New(invalidAddPeerMsgErr)
	}

	addr := host[1]
	port := host[2]
	return addr, port, nil
}

// ComposeRemovePeerMessage returns a `remove peer` message with given addr and port
func ComposeRemovePeerMessage(addr, port string) string {
	return fmt.Sprintf("%s/%s/%s", removePrefix, addr, port)
}

// DecomposeRemovePeerMessage decomposes given `remove peer` message to addr and port
func DecomposeRemovePeerMessage(msg string) (string, string, error) {
	host := strings.Split(msg, "/")
	if len(host) != 3 {
		return "", "", errors.New(invalidRemovePeerMsgErr)
	}

	if host[0] != removePrefix {
		return "", "", errors.New(invalidRemovePeerMsgErr)
	}

	addr := host[1]
	port := host[2]
	return addr, port, nil
}

// DefaultRegistry is a default callbacks registry
type DefaultRegistry struct {
	callbacks map[string]func(buffer.Message, interface{}, *log.Logger) error
}

// NewDefaultRegistry creates a default callback registry
func NewDefaultRegistry() (*DefaultRegistry, error) {
	r := &DefaultRegistry{}

	r.callbacks = map[string]func(buffer.Message, interface{}, *log.Logger) error{
		ADDPEER: func(msg buffer.Message, peersBuf interface{}, logger *log.Logger) error {
			// extract addr and peer from `add peer` message
			addr, port, err := DecomposeAddPeerMessage(msg.Msg.(string))
			if err != nil {
				return err
			}

			// add peer in buffer
			if err = peersBuf.(*peer.Buffer).AddPeer(peer.NewPeer(addr, port)); err != nil {
				return err
			}

			logger.Printf(peerAddedLogFmt, addr, port)
			return nil
		},

		// nolint:unparam
		REMOVEPEER: func(msg buffer.Message, peersBuf interface{}, logger *log.Logger) error {
			// extract addr and peer from `remove peer` message
			addr, port, err := DecomposeRemovePeerMessage(msg.Msg.(string))
			if err != nil {
				return err
			}

			// remove the peer from buffer
			peersBuf.(*peer.Buffer).RemovePeer(peer.NewPeer(addr, port))

			logger.Printf(peerRemovedLogFmt, addr, port)
			return nil
		},
	}

	return r, nil
}

// GetCallback returns a default callback from registry
func (r *DefaultRegistry) GetCallback(t string) (func(buffer.Message, interface{}, *log.Logger) error, error) {
	if v, ok := r.callbacks[t]; ok {
		return v, nil
	}
	return nil, errors.New(inexistentDefaultCallbackErr)
}

// RunCallbacks runs default callbacks.
func (r *DefaultRegistry) RunCallbacks(m buffer.Message, peerBuf *peer.Buffer, logger *log.Logger) error {

	// get callback from callbacks registry
	callbackFn, err := r.GetCallback(m.CallbackType)
	if err != nil {
		// dont't return err if default registry haven't given callback
		return nil
	}

	// TODO find a way to remove the following switch
	var p interface{}
	switch m.CallbackType {
	case ADDPEER:
		p = peerBuf
	case REMOVEPEER:
		p = peerBuf
	default:
		p = nil
	}

	// run callback function
	err = callbackFn(m, p, logger)
	if err != nil {
		return err
	}

	return nil
}
