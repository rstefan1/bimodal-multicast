package main

import (
	"github.com/rstefan1/bimodal-multicast/src/httpserver"
	. "github.com/rstefan1/bimodal-multicast/src/internal"
	"github.com/rstefan1/bimodal-multicast/src/protocol"
)

var (
	nodeBuffer []Node
	msgBuffer  []Message
)

func main() {
	httpserver.Start(&nodeBuffer, &msgBuffer)
	protocol.Start(&nodeBuffer, &msgBuffer)
}
