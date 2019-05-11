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

import "fmt"

// NOCALLBACK is type for messages without callback function
const NOCALLBACK = "no-callback"

type CallbackFn func(string) (bool, error)

type Registry struct {
	callbacks map[string]CallbackFn
}

// NewRegistry creates a callback registry
func NewRegistry(cb map[string]CallbackFn) (*Registry, error) {
	r := &Registry{}

	if cb == nil {
		return nil, fmt.Errorf("Callback map must not be empty")
	}

	r.callbacks = cb
	return r, nil
}

func (r *Registry) GetCallback(t string) (CallbackFn, error) {
	if v, ok := r.callbacks[t]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("callback type doesn't exist in registry")
}
