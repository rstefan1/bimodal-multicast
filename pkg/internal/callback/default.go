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
	"github.com/rstefan1/bimodal-multicast/pkg/internal/round"
)

// DefaultRegistry is a default callbacks registry
type DefaultRegistry struct {
	callbacks map[string]func(buffer.Message, interface{}, *log.Logger) (bool, error)
}

// NewDefaultRegistry creates a default callback registry
func NewDefaultRegistry() (*DefaultRegistry, error) {
	r := &DefaultRegistry{}
	r.callbacks = map[string]func(buffer.Message, interface{}, *log.Logger) (bool, error){
		ADDPEER: func(msg buffer.Message, peersBuf interface{}, logger *log.Logger) (bool, error) {
			host := strings.Split(msg.Msg.(string), "/")
			addr := host[0]
			port := host[1]

			err := peersBuf.(*peer.Buffer).AddPeer(peer.NewPeer(addr, port))
			if err != nil {
				return false, err
			}

			logger.Printf("Peer %s/%s added in peers buffer", addr, port)
			return true, nil
		},
		// nolint:unparam
		REMOVEPEER: func(msg buffer.Message, peersBuf interface{}, logger *log.Logger) (bool, error) {
			host := strings.Split(msg.Msg.(string), "/")
			addr := host[0]
			port := host[1]

			peersBuf.(*peer.Buffer).RemovePeer(peer.NewPeer(addr, port))
			logger.Printf("Peer %s/%s removed from peers buffer", addr, port)

			return true, nil
		},
	}
	return r, nil
}

// GetDefaultCallback returns a default callback from registry
func (r *DefaultRegistry) GetDefaultCallback(t string) (func(buffer.Message, interface{}, *log.Logger) (bool, error), error) {
	if v, ok := r.callbacks[t]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("callback type doesn't exist in default registry")
}

// RunDefaultCallbacks runs default callbacks and adds given message in buffer.
// RunDefaultCallback returns true if message was added in buffer.
func (r *DefaultRegistry) RunDefaultCallbacks(m buffer.Message, hostAddr, hostPort string,
	logger *log.Logger, msgBuf *buffer.MessageBuffer, peerBuf *peer.Buffer, gossipRound *round.GossipRound) (bool, error) {

	// get callback from callbacks registry
	callbackFn, err := r.GetDefaultCallback(m.CallbackType)
	if err != nil {
		return false, err
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
	ok, err := callbackFn(m, p, logger)
	if err != nil {
		e := fmt.Errorf("BMMC %s:%s: Error at calling callback function: %s", hostAddr, hostPort, err)
		logger.Print(e)
	}

	// add message in buffer only if callback call returns true
	if ok {
		err = msgBuf.AddMessage(m)
		if err != nil {
			e := fmt.Errorf("BMMC %s:%s error at syncing buffer with message %s in round %d: %s", hostAddr, hostPort, m.ID, gossipRound.GetNumber(), err)
			logger.Print(e)
			return true, e
		}
		logger.Printf("BMMC %s:%s synced buffer with message %s in round %d", hostAddr, hostPort, m.ID, gossipRound.GetNumber())
		return true, nil
	}

	return false, nil
}
