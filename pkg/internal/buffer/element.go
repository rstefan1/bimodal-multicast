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
	"crypto/sha1" //nolint: gosec
	"encoding/hex"
	"fmt"
	"time"
)

// Element is an element from messages buffer.
type Element struct {
	ID           string    `json:"id"`
	Timestamp    time.Time `json:"timestamp"`
	Msg          any       `json:"msg"`
	CallbackType string    `json:"callbackType"`
	GossipCount  int64     `json:"gossipCount"` // number of rounds since the element is in buffer
	Internal     bool      `json:"internal"`    // true if the element is an internal element, not a user element
}

// generateIDFromMsg returns an ID consisting of a hash of the original string,
// a timestamp and a random number.
func generateIDFromMsg(s string) (string, error) {
	h := sha1.New() //nolint: gosec

	if _, err := h.Write([]byte(s)); err != nil {
		//nolint: wrapcheck
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// NewElement creates new buffer element with given message and callback type.
func NewElement(msg any, cbType string, internal bool) (Element, error) {
	id, err := generateIDFromMsg(fmt.Sprintf("%v", msg))
	if err != nil {
		return Element{}, err
	}

	return Element{
		ID:           id,
		Timestamp:    time.Now(),
		Msg:          msg,
		CallbackType: cbType,
		GossipCount:  0,
		Internal:     internal,
	}, nil
}
