package main

import (
	"sync"

	"github.com/rstefan1/bimodal-multicast/src/httpserver"
	. "github.com/rstefan1/bimodal-multicast/src/internal/buffer"
	. "github.com/rstefan1/bimodal-multicast/src/internal/peer"
	"github.com/rstefan1/bimodal-multicast/src/protocol"
)

var (
	// shared buffer for peers
	peerBuffer = []Peer{}
	// shared buffer for messages
	msgBuffer = MessageBuffer{}
)

func main() {
	msgBuffer = msgBuffer.AddMutex(&sync.Mutex{})
	httpserver.Start(&peerBuffer, &msgBuffer)
	protocol.Start(&peerBuffer, &msgBuffer)
}
