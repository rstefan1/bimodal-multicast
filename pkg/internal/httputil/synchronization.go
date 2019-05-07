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
)

// httppSynchronization is synchronization message for http server
type HTTPSynchronization struct {
	Addr     string               `json:"http_synchronization_addr"`
	Port     string               `json:"http_synchronization_port"`
	Messages buffer.MessageBuffer `json:"http_synchronization_digests"`
}

// ReceiveSynchronization receives http solicitation message
func ReceiveSynchronization(r *http.Request) (*buffer.MessageBuffer, string, string, error) {
	decoder := json.NewDecoder(r.Body)
	var t HTTPSynchronization
	err := decoder.Decode(&t)
	if err != nil {
		return nil, "", "", fmt.Errorf("Error at decoding http synchronization message in HTTP Server: %s", err)
	}

	return &t.Messages, t.Addr, t.Port, nil
}

// SendSynchronization send http synchronization message
func SendSynchronization(synchronization HTTPSynchronization, tAddr, tPort string) error {
	jsonSynchronization, err := json.Marshal(synchronization)
	if err != nil {
		return fmt.Errorf("Error at marshal http synchronization message in HTTP Server: %s", err)
	}

	path := fmt.Sprintf("http://%s:%s/synchronization", tAddr, tPort)

	_, err = http.Post(path, "json", bytes.NewBuffer(jsonSynchronization))
	if err != nil {
		return fmt.Errorf("Error at sending HTTPSynchronization message in HTTP Server: %s", err)
	}

	return nil
}
