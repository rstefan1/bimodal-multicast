/*
Copyright 2019 Robert Andrei STEFAN

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package httpmessage

import (
	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/round"
)

// HTTPGossip is gossip message for http server
type HTTPGossip struct {
	Addr        string              `json:"http_gossip_addr"`
	Port        string              `json:"http_gossip_port"`
	RoundNumber *round.GossipRound  `json:"http_gossip_round_number"`
	Digests     buffer.DigestBuffer `json:"http_gossip_digests"`
}

// HTTPSolicitation is solicitation message for http server
type HTTPSolicitation struct {
	Addr        string              `json:"http_solicitation_addr"`
	Port        string              `json:"http_solicitation_port"`
	RoundNumber *round.GossipRound  `json:"http_solicitation_round_number"`
	Digests     buffer.DigestBuffer `json:"http_solicitation_digests"`
}

// HTTPSynchronization is synchronization message for http server
type HTTPSynchronization struct {
	Addr     string               `json:"http_synchronization_addr"`
	Port     string               `json:"http_synchronization_port"`
	Messages buffer.MessageBuffer `json:"http_synchronization_digests"`
}
