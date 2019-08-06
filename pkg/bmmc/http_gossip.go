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

package bmmc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/round"
)

// HTTPGossip is gossip message for http server
type HTTPGossip struct {
	Addr        string              `json:"addr"`
	Port        string              `json:"port"`
	RoundNumber *round.GossipRound  `json:"roundNumber"`
	Digests     buffer.DigestBuffer `json:"digests"`
}

func gossipHTTPPath(addr, port string) string {
	return fmt.Sprintf("http://%s:%s/gossip", addr, port)
}

// receiveGossip receives a HTPP gossip message
func receiveGossip(r *http.Request) (*buffer.DigestBuffer, string, string, *round.GossipRound, error) {
	var t HTTPGossip

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		return nil, "", "", nil, fmt.Errorf("Error at decoding http gossip message in HTTP Server: %s", err)
	}

	return &t.Digests, t.Addr, t.Port, t.RoundNumber, nil
}

// sendGossip sends a HTTP gossip message
func sendGossip(gossipMsg HTTPGossip, addr, port string) error {
	jsonGossip, err := json.Marshal(gossipMsg)
	if err != nil {
		return fmt.Errorf("Gossiper from %s:%s can not marshal the gossip message: %s", gossipMsg.Addr, gossipMsg.Port, err)
	}

	_, err = http.Post(gossipHTTPPath(addr, port), "json", bytes.NewBuffer(jsonGossip))
	if err != nil {
		return fmt.Errorf("Gossiper from %s:%s can not marshal the gossip message: %s", gossipMsg.Addr, gossipMsg.Port, err)
	}

	return nil
}
