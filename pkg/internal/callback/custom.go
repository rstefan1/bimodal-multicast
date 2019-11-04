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
)

const (
	nilCallbackMapErr           = "callback map must not be nil"
	inexistentCustomCallbackErr = "callback doesn't exist in the custom registry"
)

// CustomRegistry is a custom callbacks registry
type CustomRegistry struct {
	callbacks map[string]func(interface{}, *log.Logger) error
}

// NewCustomRegistry creates a custom callback registry
func NewCustomRegistry(cb map[string]func(interface{}, *log.Logger) error) (*CustomRegistry, error) {
	if cb == nil {
		return nil, errors.New(nilCallbackMapErr)
	}

	r := &CustomRegistry{}
	r.callbacks = cb

	return r, nil
}

// GetCallback returns a custom callback from registry
func (r *CustomRegistry) GetCallback(t string) (func(interface{}, *log.Logger) error, error) {
	if v, ok := r.callbacks[t]; ok {
		return v, nil
	}

	return nil, errors.New(inexistentCustomCallbackErr)
}

// RunCallbacks runs custom callbacks.
func (r *CustomRegistry) RunCallbacks(m buffer.Element, logger *log.Logger) error {
	// get callback from callbacks registry
	callbackFn, err := r.GetCallback(m.CallbackType)
	if err != nil {
		// dont't return err if custom registry haven't given callback
		return nil
	}

	// run callback function
	if err = callbackFn(m.Msg, logger); err != nil {
		return err
	}

	return nil
}
