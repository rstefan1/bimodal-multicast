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
)

const (
	httpSynchronizationDecodeErrFmt  = "error at decoding http synchronization message in HTTP Server: %s"
	httpSynchronizationMarshalErrFmt = "error at marshal http synchronization message in HTTP Server: %s"
	httpSynchronizationSendErrFmt    = "error at sending HTTPSynchronization message in HTTP Server: %s"
)

// HTTPSynchronization is synchronization message for http server
type HTTPSynchronization struct {
	Addr     string               `json:"addr"`
	Port     string               `json:"port"`
	Messages buffer.MessageBuffer `json:"messages"`
}

func synchronizationHTTPPath(addr, port string) string {
	return fmt.Sprintf("http://%s:%s%s", addr, port, synchronizationRoute)
}

// receiveSynchronization receives http solicitation message
func (b *BMMC) receiveSynchronization(r *http.Request) (*buffer.MessageBuffer, string, string, error) {
	var t HTTPSynchronization

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		return nil, "", "", fmt.Errorf(httpSynchronizationDecodeErrFmt, err)
	}

	return &t.Messages, t.Addr, t.Port, nil
}

// sendSynchronization send http synchronization message
func (b *BMMC) sendSynchronization(synchronization HTTPSynchronization, addr, port string) error {
	jsonSynchronization, err := json.Marshal(synchronization)
	if err != nil {
		return fmt.Errorf(httpSynchronizationMarshalErrFmt, err)
	}

	go func() {
		resp, err := http.Post(synchronizationHTTPPath(addr, port), "json", bytes.NewBuffer(jsonSynchronization))
		if err != nil {
			b.config.Logger.Printf(httpSynchronizationSendErrFmt, err)
			return
		}
		defer resp.Body.Close() // nolint:errcheck
	}()

	return nil
}
