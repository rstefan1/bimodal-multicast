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

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/round"
)

// CustomRegistry is a custom callbacks registry
type CustomRegistry struct {
	callbacks map[string]func(interface{}, *log.Logger) (bool, error)
}

// NewCustomRegistry creates a custom callback registry
func NewCustomRegistry(cb map[string]func(interface{}, *log.Logger) (bool, error)) (*CustomRegistry, error) {
	r := &CustomRegistry{}

	if cb == nil {
		return nil, fmt.Errorf("callback map must not be empty")
	}

	r.callbacks = cb
	return r, nil
}

// GetCustomCallback returns a custom callback from registry
func (r *CustomRegistry) GetCustomCallback(t string) (func(interface{}, *log.Logger) (bool, error), error) {
	if v, ok := r.callbacks[t]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("callback type doesn't exist in custom registry")
}

// RunCustomCallbacks runs custom callbacks and adds given message in buffer.
// RunCustomCallbacks returns true if given message was added in buffer.
func (r *CustomRegistry) RunCustomCallbacks(m buffer.Message, hostAddr, hostPort string,
	logger *log.Logger, msgBuf *buffer.MessageBuffer, gossipRound *round.GossipRound) (bool, error) {

	// get callback from callbacks registry
	callbackFn, err := r.GetCustomCallback(m.CallbackType)
	if err != nil {
		return false, err
	}

	// run callback function
	ok, err := callbackFn(m.Msg, logger)
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
			fmt.Println()
			return true, e
		}
		logger.Printf("BMMC %s:%s synced buffer with message %s in round %d", hostAddr, hostPort, m.ID, gossipRound.GetNumber())
		return true, nil
	}

	return false, nil
}
