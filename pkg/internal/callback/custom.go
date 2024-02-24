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

var (
	errNilCallbackMap        = errors.New("callback map must not be nil")
	errNotAlowedCallbackType = errors.New("callback type is not allowed")
)

// CustomRegistry is a custom callback registry.
type CustomRegistry struct {
	callbacks map[string]func(any, *log.Logger) error
}

// NewCustomRegistry creates a custom callback registry.
func NewCustomRegistry(cb map[string]func(any, *log.Logger) error) (*CustomRegistry, error) {
	if cb == nil {
		return nil, errNilCallbackMap
	}

	return &CustomRegistry{
		callbacks: cb,
	}, nil
}

// GetCallback returns a custom callback from registry.
func (r *CustomRegistry) GetCallback(t string) func(any, *log.Logger) error {
	if v, ok := r.callbacks[t]; ok {
		return v
	}

	return nil
}

// RunCallbacks runs custom callbacks.
func (r *CustomRegistry) RunCallbacks(m buffer.Element, logger *log.Logger) error {
	// get callback from callbacks registry
	callbackFn := r.GetCallback(m.CallbackType)
	if callbackFn == nil {
		return nil
	}

	// run callback function
	return callbackFn(m.Msg, logger)
}

// ValidateCustomCallbacks validates custom callbacks.
func ValidateCustomCallbacks(customCallbacks map[string]func(any, *log.Logger) error) error {
	// don't allow to use default callbacks types as custom callback types
	for customType := range customCallbacks {
		if _, exists := defaultCallbacks[customType]; exists {
			return errNotAlowedCallbackType
		}
	}

	return nil
}
