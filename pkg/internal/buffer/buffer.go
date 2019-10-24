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

package buffer

import "sync"

// Buffer is the buffer with messages
type Buffer struct {
	Elements []Element   `json:"elements"`
	Len      int         `json:"len"`
	Mux      *sync.Mutex `json:"mux"`
}

// NewBuffer creates new buffer
func NewBuffer(size int) *Buffer {
	return &Buffer{
		Elements: make([]Element, size),
		Len:      0,
		Mux:      &sync.Mutex{},
	}
}
