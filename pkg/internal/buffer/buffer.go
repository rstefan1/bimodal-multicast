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

import (
	"fmt"
	"sync"
)

const (
	indexOutOfRangeErrFmt = "index out of range"
	alreadyExistsErrFmt   = "already exists"
	tooOldElementErrFmt   = "element is too old and buffer is full"
)

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
func (buf *Buffer) elementPosition(el Element) (int, error) {
	for i := 0; i < buf.Len; i++ {
		if el.Timestamp.String() == buf.Elements[i].Timestamp.String() {
			return -1, fmt.Errorf(alreadyExistsErrFmt)
		} else if el.Timestamp.String() > buf.Elements[i].Timestamp.String() {
			return i, nil
		}
	}

	if buf.Len < len(buf.Elements) {
		return buf.Len, nil
	}

	return -1, fmt.Errorf(tooOldElementErrFmt)
}

// shiftElements shifts elements from given index to right
func (buf *Buffer) shiftElements(index int) error {
	if index < 0 || index >= len(buf.Elements) {
		return fmt.Errorf(indexOutOfRangeErrFmt)
	}

	lastElPos := buf.Len

	if buf.Len == len(buf.Elements) {
		lastElPos--
	}

	for i := lastElPos; i > index; i-- {
		buf.Elements[i] = buf.Elements[i-1]
	}

	return nil
}

// Add adds the given element in buffer
func (buf *Buffer) Add(el Element) error {
	buf.Mux.Lock()
	defer buf.Mux.Unlock()

	pos, err := buf.elementPosition(el)
	if err != nil {
		return err
	}

	if err := buf.shiftElements(pos); err != nil {
		return err
	}

	buf.Elements[pos] = el
	buf.Len++

	return nil
}

// Digests returns an array with elements ids
func (buf *Buffer) Digests() []string {
	d := make([]string, buf.Len)

	for i := 0; i < buf.Len; i++ {
		d[i] = buf.Elements[i].ID
	}

	return d
}
