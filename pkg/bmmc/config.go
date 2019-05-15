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
	"log"
)

type Peer struct {
	Addr string
	Port string
}

type Config struct {
	// Addr is HTTP address for node which runs gossip and http servers
	// Required field
	Addr string
	// Port is HTTP port for node which runs gossip  and http servers
	// Required field
	Port string
	// PeerBuf is the list of peers
	// Optional field
	Peers []Peer
	// Beta is the expected fanout for gossip rounds
	// Optional field
	Beta float64
	// Logger
	// Optional field
	Logger *log.Logger
	// Callbacks funtions
	// Optional field
	Callbacks map[string]func(string) (bool, error)
}
