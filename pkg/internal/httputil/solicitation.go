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

package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/round"
)

// HTTPSolicitation is solicitation message for http server
type HTTPSolicitation struct {
	Addr        string              `json:"http_solicitation_addr"`
	Port        string              `json:"http_solicitation_port"`
	RoundNumber *round.GossipRound  `json:"http_solicitation_round_number"`
	Digests     buffer.DigestBuffer `json:"http_solicitation_digests"`
}

// ReceiveSolicitation receives http solicitation message
func ReceiveSolicitation(r *http.Request) (*buffer.DigestBuffer, string, string, *round.GossipRound, error) {
	decoder := json.NewDecoder(r.Body)
	var t HTTPSolicitation
	err := decoder.Decode(&t)
	if err != nil {
		return nil, "", "", nil, fmt.Errorf("Error at decoding http solicitation message in HTTP Server: %s", err)
	}

	return &t.Digests, t.Addr, t.Port, t.RoundNumber, nil
}

// SendSolicitation send http solicitation message
func SendSolicitation(solicitation HTTPSolicitation, tAddr, tPort string) error {
	jsonSolicitation, err := json.Marshal(solicitation)
	if err != nil {
		return fmt.Errorf("Error at marshal http solicitation in HTTP Server: %s", err)
	}

	path := fmt.Sprintf("http://%s:%s/solicitation", tAddr, tPort)

	_, err = http.Post(path, "json", bytes.NewBuffer(jsonSolicitation))
	if err != nil {
		return fmt.Errorf("Error at sending http soliccitation message in HTTP Server: %s", err)
	}

	return nil
}
