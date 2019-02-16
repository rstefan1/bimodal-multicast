package gossip

import "github.com/rstefan1/bimodal-multicast/src/internal/buffer"

type GossipMessage struct {
	RoundNumber int
	Digest      buffer.DigestBuffer
}
