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
	callbacks map[string]func(buffer.Message, interface{}, *log.Logger) (bool, error)
}

// NewDefaultRegistry creates a default callback registry
func NewDefaultRegistry() (*DefaultRegistry, error) {
	r := &DefaultRegistry{}
	r.callbacks = map[string]func(buffer.Message, interface{}, *log.Logger) (bool, error){
		ADDPEER: func(msg buffer.Message, peersBuf interface{}, logger *log.Logger) (bool, error) {
			host := strings.Split("localhost/19999", "/")
			addr := host[0]
			port := host[1]

			ok := peersBuf.(*peer.PeerBuffer).AddPeer(peer.NewPeer(addr, port))
			if !ok {
				logger.Printf("Error at adding %s/%s peer in peers buffer in %s callback", addr, port, ADDPEER)
			}

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
