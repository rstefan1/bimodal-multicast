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
	solicitationDecodingErrFmt = "error at decoding http solicitation message in Server: %w"
	solicitationMarshalErrFmt  = "error at marshal http solicitation in Server: %w"
	solicitationSendLogFmt     = "error at sending http solicitation message in Server: %s"
)

// Solicitation is solicitation message for server.
type Solicitation struct {
	Host        string       `json:"host"`
	RoundNumber *GossipRound `json:"roundNumber"`
	Digest      []string     `json:"digest"`
}

// receiveSolicitation receives http solicitation message.
func (b *BMMC) receiveSolicitation(msg []byte) ([]string, string, *GossipRound, error) {
	var body Solicitation

	if err := json.Unmarshal(msg, &body); err != nil {
		return nil, "", nil, fmt.Errorf(solicitationDecodingErrFmt, err)
	}

	return body.Digest, body.Host, body.RoundNumber, nil
}

// sendSolicitation send http solicitation message.
func (b *BMMC) sendSolicitation(solicitation Solicitation, peerToSend string) error {
	jsonSolicitation, err := json.Marshal(solicitation)
	if err != nil {
		return fmt.Errorf(solicitationMarshalErrFmt, err)
	}

	go func() {
		if err := b.config.Host.Send(jsonSolicitation, SolicitationRoute, peerToSend); err != nil {
			b.config.Logger.Printf(solicitationSendLogFmt, err)
		}
	}()

	return nil
}
