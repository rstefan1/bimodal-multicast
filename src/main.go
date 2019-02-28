package bimodal

import (
	"sync"

	"github.com/rstefan1/bimodal-multicast/src/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/src/internal/config"
	"github.com/rstefan1/bimodal-multicast/src/internal/peer"
	"github.com/rstefan1/bimodal-multicast/src/protocol"
)

var (
	// shared buffer for peers
	peerBuffer = []peer.Peer{}
	// shared buffer for messages
	msgBuffer = buffer.MessageBuffer{}
)

func BMC() error {
	msgBuffer = msgBuffer.AddMutex(&sync.Mutex{})
	stop := make(chan struct{})
	cfg := config.Config{}

	if err := httpserver.Start(&peerBuffer, &msgBuffer, stop, cfg); err != nil {
		return err
	}
	protocol.Start(&peerBuffer, &msgBuffer)

	return nil
}
