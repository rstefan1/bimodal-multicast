package main

import (
	"github.com/rstefan1/bimodal-multicast/src/httpserver"
	. "github.com/rstefan1/bimodal-multicast/src/internal"
	. "github.com/rstefan1/bimodal-multicast/src/internal/buffer"
	. "github.com/rstefan1/bimodal-multicast/src/internal/peer"
	"github.com/rstefan1/bimodal-multicast/src/protocol"
)

var (
	peerBuffer []Peer
	msgBuffer  MessageBuffer
)

func main() {
	httpserver.Start(&peerBuffer, &msgBuffer)
	protocol.Start(&peerBuffer, &msgBuffer)
}
