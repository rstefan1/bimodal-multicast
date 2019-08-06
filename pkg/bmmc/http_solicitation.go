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

// HTTPSolicitation is solicitation message for http server
type HTTPSolicitation struct {
	Addr        string              `json:"addr"`
	Port        string              `json:"port"`
	RoundNumber *round.GossipRound  `json:"roundNumber"`
	Digests     buffer.DigestBuffer `json:"digests"`
}

func solicitationHTTPPath(addr, port string) string {
	return fmt.Sprintf("http://%s:%s/solicitation", addr, port)
}

// receiveSolicitation receives http solicitation message
func receiveSolicitation(r *http.Request) (*buffer.DigestBuffer, string, string, *round.GossipRound, error) {
	var t HTTPSolicitation

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		return nil, "", "", nil, fmt.Errorf("Error at decoding http solicitation message in HTTP Server: %s", err)
	}

	return &t.Digests, t.Addr, t.Port, t.RoundNumber, nil
}

// sendSolicitation send http solicitation message
func sendSolicitation(solicitation HTTPSolicitation, addr, port string) error {
	jsonSolicitation, err := json.Marshal(solicitation)
	if err != nil {
		return fmt.Errorf("Error at marshal http solicitation in HTTP Server: %s", err)
	}

	_, err = http.Post(solicitationHTTPPath(addr, port), "json", bytes.NewBuffer(jsonSolicitation))
	if err != nil {
		return fmt.Errorf("Error at sending http soliccitation message in HTTP Server: %s", err)
	}

	return nil
}
