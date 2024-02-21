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
	"errors"
	"math"
	"sync"
)

var (
	errIndexOutOfRange = errors.New("index out of range")
	errAlreadyExists   = errors.New("already exists")
	errTooOldElement   = errors.New("element is too old and buffer is full")
)

// Buffer is the buffer with messages.
type Buffer struct {
	Elements []Element   `json:"elements"`
	Len      int         `json:"len"`
	Mux      *sync.Mutex `json:"mux"`
}

// NewBuffer creates new buffer.
func NewBuffer(size int) *Buffer {
	return &Buffer{
		Elements: make([]Element, size),
		Len:      0,
		Mux:      &sync.Mutex{},
	}
}

// contains returns a boolean representing if given element already exists in buffer or not
// and an int representing element position in buffer.
func (buf *Buffer) contains(el Element) (bool, int) {
	for i := 0; i < buf.Len; i++ {
		if buf.Elements[i].ID == el.ID {
			return true, i
		}
	}

	return false, -1
}

// elementPosition gets the element position in buffer.
func (buf *Buffer) elementPosition(el Element) (int, error) {
	for i := 0; i < buf.Len; i++ {
		if el.Timestamp.String() >= buf.Elements[i].Timestamp.String() {
			return i, nil
		}
	}

	if buf.Len < len(buf.Elements) {
		return buf.Len, nil
	}

	return -1, errTooOldElement
}

// shiftElements shifts elements from given index to right.
func (buf *Buffer) shiftElements(index int) error {
	if index < 0 || index >= len(buf.Elements) {
		return errIndexOutOfRange
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

// Add adds the given element in buffer.
func (buf *Buffer) Add(el Element) error {
	buf.Mux.Lock()
	defer buf.Mux.Unlock()

	if e, _ := buf.contains(el); e {
		return errAlreadyExists
	}

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

// Digest returns a slice with elements ids.
func (buf *Buffer) Digest() []string {
	buf.Mux.Lock()
	defer buf.Mux.Unlock()

	d := make([]string, buf.Len)

	for i := 0; i < buf.Len; i++ {
		d[i] = buf.Elements[i].ID
	}

	return d
}

// IncrementGossipCount increments gossip count for each elements from buffer.
func (buf *Buffer) IncrementGossipCount() {
	buf.Mux.Lock()
	defer buf.Mux.Unlock()

	for i := 0; i < buf.Len; i++ {
		if buf.Elements[i].GossipCount == math.MaxInt64 {
			buf.Elements[i].GossipCount = int64(0)

			continue
		}

		buf.Elements[i].GossipCount++
	}
}

// Messages returns a slice with messages for each element in buffer.
// If withInternals parameter is false, Messages returns only user (not internal) messages.
func (buf *Buffer) Messages(withInternals bool) []interface{} {
	buf.Mux.Lock()
	defer buf.Mux.Unlock()

	msgs := []interface{}{}

	for i := 0; i < buf.Len; i++ {
		if !withInternals && buf.Elements[i].Internal {
			continue // don't add internal messages if `withInternal` param si false
		}

		msgs = append(msgs, buf.Elements[i].Msg)
	}

	return msgs
}

// Length returns number of elements in buffer.
func (buf *Buffer) Length() int {
	buf.Mux.Lock()
	defer buf.Mux.Unlock()

	l := buf.Len

	return l
}

// ElementsFromIDs returns a slice with elements from given IDs list.
func (buf *Buffer) ElementsFromIDs(digest []string) []Element {
	buf.Mux.Lock()
	defer buf.Mux.Unlock()

	el := []Element{}

	for i := 0; i < buf.Len; i++ {
		if ContainsString(digest, buf.Elements[i].ID) {
			el = append(el, buf.Elements[i])
		}
	}

	return el
}
