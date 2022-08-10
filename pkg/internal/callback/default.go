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

package callback

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

const (
	peerAddedLogFmt   = "peer %s/%s added in the peers buffer"
	peerRemovedLogFmt = "peer %s/%s removed from the peers buffer"

	// addPrefix is the prefix for add peer messages.
	addPrefix = "add"
	// removePrefix is the prefix for remove peer messages.
	removePrefix = "remove"

	hostBlocksLen = 3
)

var (
	//nolint: gochecknoglobals
	defaultCallbacks = map[string]func(buffer.Element, *peer.Buffer, *log.Logger) error{
		ADDPEER:    addPeerCallback,
		REMOVEPEER: removePeerCallback,
	}

	errInvalidAddPeerMsg     = errors.New("invalid add peer message")
	errInvalidRemovePeerMsg  = errors.New("invalid remove peer message")
	errCannotConvertToString = errors.New("cannot convert the given message to string")
)

// ComposeAddPeerMessage returns a `add peer` message with given addr and port.
func ComposeAddPeerMessage(addr, port string) string {
	return fmt.Sprintf("%s/%s/%s", addPrefix, addr, port)
}

// DecomposeAddPeerMessage decomposes given `add peer` message to addr and port.
func DecomposeAddPeerMessage(msg string) (string, string, error) {
	host := strings.Split(msg, "/")
	if len(host) != hostBlocksLen {
		return "", "", errInvalidAddPeerMsg
	}

	if host[0] != addPrefix {
		return "", "", errInvalidAddPeerMsg
	}

	addr := host[1]
	port := host[2]

	return addr, port, nil
}

// ComposeRemovePeerMessage returns a `remove peer` message with given addr and port.
func ComposeRemovePeerMessage(addr, port string) string {
	return fmt.Sprintf("%s/%s/%s", removePrefix, addr, port)
}

// DecomposeRemovePeerMessage decomposes given `remove peer` message to addr and port.
func DecomposeRemovePeerMessage(msg string) (string, string, error) {
	host := strings.Split(msg, "/")
	if len(host) != hostBlocksLen {
		return "", "", errInvalidRemovePeerMsg
	}

	if host[0] != removePrefix {
		return "", "", errInvalidRemovePeerMsg
	}

	addr := host[1]
	port := host[2]

	return addr, port, nil
}

// DefaultRegistry is a default callbacks registry.
type DefaultRegistry struct {
	callbacks map[string]func(buffer.Element, *peer.Buffer, *log.Logger) error
}

// NewDefaultRegistry creates a default callback registry.
func NewDefaultRegistry() (*DefaultRegistry, error) {
	return &DefaultRegistry{
		callbacks: defaultCallbacks,
	}, nil
}

// GetCallback returns a default callback from registry.
func (r *DefaultRegistry) GetCallback(t string) func(buffer.Element, *peer.Buffer, *log.Logger) error {
	if v, ok := r.callbacks[t]; ok {
		return v
	}

	return nil
}

// RunCallbacks runs default callbacks.
func (r *DefaultRegistry) RunCallbacks(m buffer.Element, peerBuf *peer.Buffer, logger *log.Logger) error {
	callbackFn := r.GetCallback(m.CallbackType)
	if callbackFn == nil {
		return nil
	}

	// run callback function
	return callbackFn(m, peerBuf, logger)
}

func addPeerCallback(msg buffer.Element, peersBuf *peer.Buffer, logger *log.Logger) error {
	strMsg, convOk := msg.Msg.(string)
	if !convOk {
		return errCannotConvertToString
	}

	// extract addr and peer from `add peer` message
	addr, port, err := DecomposeAddPeerMessage(strMsg)
	if err != nil {
		return err
	}

	// add peer in buffer
	p, err := peer.NewPeer(addr, port)
	if err != nil {
		return err //nolint: wrapcheck
	}

	if err = peersBuf.AddPeer(p); err != nil {
		return err //nolint: wrapcheck
	}

	logger.Printf(peerAddedLogFmt, addr, port)

	return nil
}

func removePeerCallback(msg buffer.Element, peersBuf *peer.Buffer, logger *log.Logger) error {
	strMsg, convOk := msg.Msg.(string)
	if !convOk {
		return errCannotConvertToString
	}

	// extract addr and peer from `remove peer` message
	addr, port, err := DecomposeRemovePeerMessage(strMsg)
	if err != nil {
		return err
	}

	// remove the peer from buffer
	p, err := peer.NewPeer(addr, port)
	if err != nil {
		return err //nolint: wrapcheck
	}

	peersBuf.RemovePeer(p)

	logger.Printf(peerRemovedLogFmt, addr, port)

	return nil
}
