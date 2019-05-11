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

type CallbacksRegistry struct {
	callbacks map[string]func(string) error
}

// NewRegistry creates a callback registry
func NewRegistry() *CallbacksRegistry {
	r := &CallbacksRegistry{}
	r.callbacks = make(map[string]func(string) error)
	return r
}

// Register registers a new callback in registry
func (r *CallbacksRegistry) Register(t string, fn func(string) error) error {
	if _, ok := r.callbacks[t]; ok {
		return fmt.Errorf("callback type already exists in registry")
	}

	r.callbacks[t] = fn
	return nil
}

func (r *CallbacksRegistry) Get(t string) (func(string) error, error) {
	if _, ok := r.callbacks[t]; ok {
		return r.callbacks[t], nil
	}
	return nil, fmt.Errorf("callback type doesn't exist in registry")
}
