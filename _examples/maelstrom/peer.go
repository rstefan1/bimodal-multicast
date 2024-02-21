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
