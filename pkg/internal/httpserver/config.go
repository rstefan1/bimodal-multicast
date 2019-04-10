package httpserver

import (
	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

// Config has configs for http server
type Config struct {
	// Addr is http server address
	Addr string
	// Port is http server port
	Port string
	// PeerBuf is the list of peers
	PeerBuf []peer.Peer
	// MsgBuf is the list of messages
	MsgBuf *buffer.MessageBuffer
}
