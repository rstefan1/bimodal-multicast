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

type Registry struct {
	callbacks map[string]func(string) (bool, error)
}

// NewRegistry creates a callback registry
func NewRegistry() *Registry {
	r := &Registry{}
	r.callbacks = make(map[string]func(string) (bool, error))
	return r
}

// Register registers a new callback in registry
func (r *Registry) Register(t string, fn func(string) (bool, error)) error {
	if _, ok := r.callbacks[t]; ok {
		return fmt.Errorf("callback type already exists in registry")
	}

	r.callbacks[t] = fn
	return nil
}

func (r *Registry) GetCallback(t string) (func(string) (bool, error), error) {
	if v, ok := r.callbacks[t]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("callback type doesn't exist in registry")
}
