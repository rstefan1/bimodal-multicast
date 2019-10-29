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

// elementPosition gets the element position in buffer
func (buf *Buffer) elementPosition(el Element) int {
	for i := 0; i < buf.Len; i++ {
		if el.Timestamp.String() >= buf.Elements[i].Timestamp.String() {
			return i
		}
	}

	return -1
}

// shiftElements shifts elements from given index to right
func (buf *Buffer) shiftElements(index int) {
	lastElPos := buf.Len

	if buf.Len == len(buf.Elements) {
		lastElPos--
	}

	for i := lastElPos; i > index; i-- {
		buf.Elements[i] = buf.Elements[i-1]
	}
}

// Add adds the given element in buffer
func (buf *Buffer) Add(el Element) {
	buf.Mux.Lock()
	defer buf.Mux.Unlock()

	pos := buf.elementPosition(el)
	if pos == -1 {
		return
	}

	buf.shiftElements(pos)

	buf.Elements[pos] = el
}
