package gossipserver

import (
	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

// Config has configs for gossip object
type Config struct {
	// Addr is HTTP address for node which runs gossip round
	Addr string
	// Port is HTTP port for node which runs gossip round
	Port string
	// PeerBuf is the list of peers
	PeerBuf []peer.Peer
	// MsgBuf is the list of messages
	MsgBuf *buffer.MessageBuffer
}
