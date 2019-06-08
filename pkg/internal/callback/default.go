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
	"fmt"
	"log"
	"strings"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

// DefaultRegistry is a default callbacks registry
type DefaultRegistry struct {
	callbacks map[string]func(buffer.Message, interface{}, *log.Logger) error
}

// NewDefaultRegistry creates a default callback registry
func NewDefaultRegistry() (*DefaultRegistry, error) {
	r := &DefaultRegistry{}
	r.callbacks = map[string]func(buffer.Message, interface{}, *log.Logger) error{
		ADDPEER: func(msg buffer.Message, peersBuf interface{}, logger *log.Logger) error {
			host := strings.Split(msg.Msg.(string), "/")
			addr := host[0]
			port := host[1]

			err := peersBuf.(*peer.Buffer).AddPeer(peer.NewPeer(addr, port))
			if err != nil {
				return err
			}

			logger.Printf("Peer %s/%s added in peers buffer", addr, port)
			return nil
		},
		// nolint:unparam
		REMOVEPEER: func(msg buffer.Message, peersBuf interface{}, logger *log.Logger) error {
			host := strings.Split(msg.Msg.(string), "/")
			addr := host[0]
			port := host[1]

			peersBuf.(*peer.Buffer).RemovePeer(peer.NewPeer(addr, port))
			logger.Printf("Peer %s/%s removed from peers buffer", addr, port)

			return nil
		},
	}
	return r, nil
}

// GetDefaultCallback returns a default callback from registry
func (r *DefaultRegistry) GetDefaultCallback(t string) (func(buffer.Message, interface{}, *log.Logger) error, error) {
	if v, ok := r.callbacks[t]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("callback type doesn't exist in default registry")
}

// RunDefaultCallbacks runs default callbacks.
func (r *DefaultRegistry) RunDefaultCallbacks(m buffer.Message, peerBuf *peer.Buffer, logger *log.Logger) error {

	// get callback from callbacks registry
	callbackFn, err := r.GetDefaultCallback(m.CallbackType)
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
