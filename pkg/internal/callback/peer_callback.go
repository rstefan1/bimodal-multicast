/*
Copyright 2024 Robert Andrei STEFAN

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

package callback

import (
	"errors"
	"log/slog"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

const (
	// ADDPEER is the type of messages used for adding new peer in peers buffer.
	ADDPEER = "add-peer"
	// REMOVEPEER is the type of messages used for deleting a peer from peers buffer.
	REMOVEPEER = "remove-peer"
)

var (
	errCannotConvertToString           = errors.New("cannot convert the given message to string")
	errCannotConvertToPeerCallbackData = errors.New("cannot convert the given message to PeerCallbackData")
)

// PeerCallbackData contains data for add-peer & remove-peer callbacks.
type PeerCallbackData struct {
	Element buffer.Element
	Buffer  *peer.Buffer
}

// AddPeerCallback is the callback for adding peers in peers buffer.
func AddPeerCallback(data any, logger *slog.Logger) error {
	peerCBData, convOk := data.(PeerCallbackData)
	if !convOk {
		return errCannotConvertToPeerCallbackData
	}

	p, convOk := peerCBData.Element.Msg.(string)
	if !convOk {
		return errCannotConvertToString
	}

	// add peer in buffer
	if added := peerCBData.Buffer.AddPeer(p); !added {
		logger.Debug("peer already exists", "peer", p)

		return nil
	}

	logger.Debug("new peer added", "peer", p)

	return nil
}

// RemovePeerCallback is the callback for removing peers from peers buffer.
func RemovePeerCallback(data any, logger *slog.Logger) error {
	peerCBData, convOk := data.(PeerCallbackData)
	if !convOk {
		return errCannotConvertToPeerCallbackData
	}

	p, convOk := peerCBData.Element.Msg.(string)
	if !convOk {
		return errCannotConvertToString
	}

	peerCBData.Buffer.RemovePeer(p)

	logger.Debug("peer removed", "peer", p)

	return nil
}
