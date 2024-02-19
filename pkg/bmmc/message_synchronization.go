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

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
)

const (
	synchronizationDecodeErrFmt  = "error at decoding synchronization message in Server: %w"
	synchronizationMarshalErrFmt = "error at marshal synchronization message in Server: %w"
	synchronizationSendErrFmt    = "error at sending Synchronization message in Server: %s"
)

// Synchronization is synchronization message for server.
type Synchronization struct {
	Host     string           `json:"host"`
	Elements []buffer.Element `json:"elements"`
}

// receiveSynchronization receives http solicitation message.
func (b *BMMC) receiveSynchronization(msg []byte) ([]buffer.Element, string, error) {
	var body Synchronization

	if err := json.Unmarshal(msg, &body); err != nil {
		return nil, "", fmt.Errorf(synchronizationDecodeErrFmt, err)
	}

	return body.Elements, body.Host, nil
}

// sendSynchronization send http synchronization message.
func (b *BMMC) sendSynchronization(synchronization Synchronization, peerToSend string) error {
	jsonSynchronization, err := json.Marshal(synchronization)
	if err != nil {
		return fmt.Errorf(synchronizationMarshalErrFmt, err)
	}

	go func() {
		if err := b.config.Host.Send(jsonSynchronization, synchronizationRoute, peerToSend); err != nil {
			b.config.Logger.Printf(synchronizationSendErrFmt, err)
		}
	}()

	return nil
}
