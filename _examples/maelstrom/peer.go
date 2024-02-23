/*
Copyright 2024 Robert Andrei STEFAN

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

package main

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// Peer represents a maelstrom node.
type Peer struct {
	*maelstrom.Node
}

// String returns the peer as string.
func (p Peer) String() string {
	return p.Node.ID()
}

// Send sends a request.
func (p Peer) Send(msg []byte, route string, peerToSend string) error {
	body := map[string]string{}

	body[typeBodyKey] = route
	body[messageBodyKey] = string(msg)

	return p.Node.Send(peerToSend, body)
}
