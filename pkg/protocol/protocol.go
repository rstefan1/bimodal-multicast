package protocol

// import (
// 	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
// 	"github.com/rstefan1/bimodal-multicast/pkg/internal/config"
// 	"github.com/rstefan1/bimodal-multicast/pkg/internal/httpserver"
// 	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
// )

// var (
// 	// shared buffer for peers
// 	peerBuffer = []peer.Peer{}
// 	// shared buffer for messages
// 	msgBuffer = buffer.NewMessageBuffer()
// )

// func BMC() error {
// 	stop := make(chan struct{})
// 	cfg := config.Config{}

// 	if err := httpserver.Start(&peerBuffer, &msgBuffer, stop, cfg); err != nil {
// 		return err
// 	}
// 	gossipgossip.Start(&peerBuffer, &msgBuffer)

// 	return nil
// }
