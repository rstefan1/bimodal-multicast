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
	"encoding/json"
	"fmt"
)

const (
	gossipDecodingErrFmt = "error at decoding gossip message in Server: %w"
	gossipMarshalErrFmt  = "gossiper from %s can not marshal the gossip message: %w"
)

// Gossip is gossip message for http server.
type Gossip struct {
	Host        string       `json:"host"`
	RoundNumber *GossipRound `json:"roundNumber"`
	Digest      []string     `json:"digest"`
}

// receiveGossip receives a gossip message.
func (b *BMMC) receiveGossip(msg []byte) ([]string, string, *GossipRound, error) {
	var body Gossip

	if err := json.Unmarshal(msg, &body); err != nil {
		b.config.Logger.Error("cannot decode gossip message", "err", err)

		return nil, "", nil, fmt.Errorf(gossipDecodingErrFmt, err)
	}

	return body.Digest, body.Host, body.RoundNumber, nil
}

// sendGossip sends a gossip message.
func (b *BMMC) sendGossip(gossipMsg Gossip, peerToSend string) error {
	jsonGossip, err := json.Marshal(gossipMsg)
	if err != nil {
		b.config.Logger.Error("cannot marshal gossip message", "err", err)

		return fmt.Errorf(gossipMarshalErrFmt, gossipMsg.Host, err)
	}

	go func() {
		if err := b.config.Host.Send(jsonGossip, GossipRoute, peerToSend); err != nil {
			b.config.Logger.Error("cannot send gossip message to peer", "err", err)
		}
	}()

	return nil
}
