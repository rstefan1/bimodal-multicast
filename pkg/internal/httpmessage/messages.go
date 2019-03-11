package httpmessage

import "github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"

// HTTPGossip is gossip message for http server
type HTTPGossip struct {
	Addr        string              `json:"http_gossip_addr"`
	Port        string              `json:"http_gossip_port"`
	RoundNumber int64               `json:"http_gossip_round_number"`
	Digests     buffer.DigestBuffer `json:"http_gossip_digests"`
}

// HTTPSolicitation is solicitation message for http server
type HTTPSolicitation struct {
	Addr        string              `json:"http_solicitation_addr"`
	Port        string              `json:"http_solicitation_port"`
	RoundNumber int64               `json:"http_solicitation_round_number"`
	Digests     buffer.DigestBuffer `json:"http_solicitation_digests"`
}

// HTTPSynchronization is synchronization message for http server
type HTTPSynchronization struct {
	Addr     string               `json:"http_synchronization_addr"`
	Port     string               `json:"http_synchronization_port"`
	Messages buffer.MessageBuffer `json:"http_synchronization_digests"`
}
