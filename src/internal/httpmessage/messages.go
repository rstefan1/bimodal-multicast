package httpmessage

import "github.com/rstefan1/bimodal-multicast/src/internal/buffer"

type HTTPGossip struct {
	Addr        string              `json:"http_gossip_addr"`
	Port        string              `json:"http_gossip_port"`
	RoundNumber int                 `json:"http_gossip_round_number"`
	Digests     buffer.DigestBuffer `json:"http_gossip_digests"`
}

type HTTPSolicitation struct {
	Addr        string              `json:"http_solicitation_addr"`
	Port        string              `json:"http_solicitation_port"`
	RoundNumber int                 `json:"http_solicitation_round_number"`
	Digests     buffer.DigestBuffer `json:"http_solicitation_digests"`
}
