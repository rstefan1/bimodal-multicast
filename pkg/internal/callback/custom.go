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
)

// NOCALLBACK is type for messages without callback function
const NOCALLBACK = "no-callback"

// CustomRegistry is a custom callbacks registry
type CustomRegistry struct {
	callbacks map[string]func(string) (bool, error)
}

// NewCustomRegistry creates a custom callback registry
func NewCustomRegistry(cb map[string]func(string) (bool, error)) (*CustomRegistry, error) {
	r := &CustomRegistry{}

	if cb == nil {
		return nil, fmt.Errorf("Callback map must not be empty")
	}

	r.callbacks = cb
	return r, nil
}

// GetCustomCallback returns a custom callback from registry
func (r *CustomRegistry) GetCustomCallback(t string) (func(string) (bool, error), error) {
	if v, ok := r.callbacks[t]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("callback type doesn't exist in registry")
}
