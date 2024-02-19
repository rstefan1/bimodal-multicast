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

package peer

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/validators"
)

// HTTPPeer is a peer over http.
type HTTPPeer struct {
	addr       string
	port       string
	httpClient *http.Client
}

// String prints the http peer.
func (p HTTPPeer) String() string {
	return fmt.Sprintf("%s/%s", p.addr, p.port)
}

// decodeHTTPPeer returns addr and port of given http peer.
func decodeHTTPPeer(encodedPeer string) (string, string) {
	decodedPeer := strings.Split(encodedPeer, "/")

	return decodedPeer[0], decodedPeer[1]
}

func httpPath(addr, port, route string) string {
	return fmt.Sprintf("http://%s%s", net.JoinHostPort(addr, port), route)
}

// Send sends a request.
func (p HTTPPeer) Send(msg []byte, route string, peerToSend string) error {
	addr, port := decodeHTTPPeer(peerToSend)

	resp, err := p.httpClient.Post(httpPath(addr, port, route), "json", bytes.NewBuffer(msg)) //nolint: noctx

	defer resp.Body.Close() //nolint: errcheck, govet, staticcheck

	return err //nolint: wrapcheck
}

// Addr returns addr of http peer.
func (p HTTPPeer) Addr() string {
	return p.addr
}

// Port returns port of http peer.
func (p HTTPPeer) Port() string {
	return p.port
}

// NewHTTPPeer creates a HTTPPeer.
func NewHTTPPeer(addr, port string, httpClient *http.Client) (Peer, error) { //nolint: ireturn
	if err := validators.AddrValidator()(addr); err != nil {
		return HTTPPeer{}, err
	}

	if err := validators.PortAsStringValidator()(port); err != nil {
		return HTTPPeer{}, err
	}

	return HTTPPeer{
		addr:       addr,
		port:       port,
		httpClient: httpClient,
	}, nil
}
