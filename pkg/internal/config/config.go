package config

import (
	"github.com/rstefan1/bimodal-multicast/src/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/src/internal/peer"
)

// HTTPConfig has configs for http server
type HTTPConfig struct {
	// Addr is http server address
	Addr string
	// Port is http server port
	Port string
	// PeerBuf is the list of peers
	PeerBuf *[]peer.Peer
	// MsgBuf is the list of messages
	MsgBuf *buffer.MessageBuffer
}

// GossipConfig has configs for gossip object
type GossipConfig struct {
	// Addr is HTTP address for node which runs gossip round
	Addr string
	// Port is HTTP port for node which runs gossip round
	Port string
	// PeerBuf is the list of peers
	PeerBuf *[]peer.Peer
	// MsgBuf is the list of messages
	MsgBuf *buffer.MessageBuffer
}
