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
	"log/slog"
)

// NOCALLBACK is the type of messages without callback.
const NOCALLBACK = "no-callback"

var (
	errNilCallbackMap         = errors.New("callback map must not be nil")
	errNotAllowedCallbackType = errors.New("callback type is not allowed")
)

// Registry is a registry for callbacks.
type Registry struct {
	Callbacks map[string]func(any, *slog.Logger) error
}

// NewRegistry creates a callback registry.
func NewRegistry(cb map[string]func(any, *slog.Logger) error) (*Registry, error) {
	if cb == nil {
		return nil, errNilCallbackMap
	}

	return &Registry{
		Callbacks: cb,
	}, nil
}

// GetCallback returns a callback from registry.
func (r *Registry) GetCallback(t string) func(any, *slog.Logger) error {
	if v, ok := r.Callbacks[t]; ok {
		return v
	}

	return nil
}

// ValidateCustomCallbacks validates custom callbacks.
// NOTE: must be called before adding internal callbacks.
func ValidateCustomCallbacks(customCallbacks map[string]func(any, *slog.Logger) error) error {
	// don't allow to use internal callbacks types as custom callback types
	for customType := range customCallbacks {
		if customType == ADDPEER || customType == REMOVEPEER {
			return errNotAllowedCallbackType
		}
	}

	return nil
}
