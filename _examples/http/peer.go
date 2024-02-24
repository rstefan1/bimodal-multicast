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

package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// Peer decorates a Peer over HTTP.
type Peer struct {
	Addr       string
	Port       string
	httpClient *http.Client
}

// String prints the http peer.
func (p Peer) String() string {
	return fmt.Sprintf("%s:%s", p.Addr, p.Port)
}

// decodePeer returns addr and port of given http peer.
func decodePeer(encodedPeer string) (string, string) {
	decodedPeer := strings.Split(encodedPeer, ":")

	return decodedPeer[0], decodedPeer[1]
}

func httpPath(addr, port, route string) string {
	return fmt.Sprintf("http://%s%s", net.JoinHostPort(addr, port), route)
}

// Send sends a request.
func (p Peer) Send(msg []byte, route string, peerToSend string) error {
	addr, port := decodePeer(peerToSend)

	resp, err := p.httpClient.Post(httpPath(addr, port, route), "json", bytes.NewBuffer(msg)) //nolint: noctx
	if err != nil {
		return err //nolint: wrapcheck
	}

	return resp.Body.Close() //nolint: wrapcheck
}

// NewPeer creates a Peer.
func NewPeer(addr, port string, httpClient *http.Client) (Peer, error) {
	if err := addrValidator()(addr); err != nil {
		return Peer{}, err
	}

	if err := portAsStringValidator()(port); err != nil {
		return Peer{}, err
	}

	return Peer{
		Addr:       addr,
		Port:       port,
		httpClient: httpClient,
	}, nil
}
