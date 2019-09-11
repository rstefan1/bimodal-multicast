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
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Peer Buffer Interface", func() {
	When("Length() is called", func() {
		It("returns the length of given buffer", func() {
			peers := []Peer{
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "20000"},
			}
			pBuf := &Buffer{
				peers: peers,
				mux:   &sync.Mutex{},
			}

			Expect(pBuf.Length()).To(Equal(len(peers)))
		})
	})

	DescribeTable("when alreadyExists() is called",
		func(peers []Peer, p Peer, expected bool) {
			pBuf := &Buffer{
				peers: peers,
				mux:   &sync.Mutex{},
			}
			Expect(pBuf.alreadyExists(p)).To(Equal(expected))
		},

		Entry("returns true if peers is at the beginning",
			[]Peer{
				{addr: "localhost", port: "55555"},
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "20000"},
			},
			Peer{addr: "localhost", port: "55555"},
			true),
		Entry("returns true if peers is in the middle",
			[]Peer{
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "55555"},
				{addr: "localhost", port: "20000"},
			},
			Peer{addr: "localhost", port: "55555"},
			true),
		Entry("returns true if peers is at the end",
			[]Peer{
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "20000"},
				{addr: "localhost", port: "55555"},
			},
			Peer{addr: "localhost", port: "55555"},
			true),
		Entry("returns false if peer doesn't exist in buffer",
			[]Peer{
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "20000"},
			},
			Peer{addr: "localhost", port: "55555"},
			false),
	)

	DescribeTable("when AddPeer() is called",
		func(peers []Peer, p Peer, expectError bool, expectedPeers []Peer) {
			pBuf := &Buffer{
				peers: peers,
				mux:   &sync.Mutex{},
			}

			err := pBuf.AddPeer(p)
			if expectError {
				Expect(err).To(Succeed())
			} else {
				Expect(err).To(Not(Succeed()))
			}
			Expect(pBuf.peers).To(ConsistOf(expectedPeers))
		},

		Entry("adds peer in the peers buffer when it doesn't exist",
			[]Peer{
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "20000"},
			},
			Peer{addr: "localhost", port: "30000"},
			true,
			[]Peer{
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "20000"},
				{addr: "localhost", port: "30000"},
			},
		),
		Entry("doens't add peer in the peers buffer when it already exists",
			[]Peer{
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "20000"},
			},
			Peer{addr: "localhost", port: "10000"},
			false,
			[]Peer{
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "20000"},
			},
		),
	)

	When("GetPeers() is called", func() {
		It("returns a slice of strings with peers", func() {
			peers := []Peer{
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "20000"},
			}
			pBuf := &Buffer{
				peers: peers,
				mux:   &sync.Mutex{},
			}
			expectedPeers := []string{
				"localhost/10000",
				"localhost/20000",
			}

			Expect(pBuf.GetPeers()).To(ConsistOf(expectedPeers))
		})
	})

	DescribeTable("when RemovePeer() is called",
		func(peers []Peer, p Peer, expectedPeers []Peer) {
			pBuf := &Buffer{
				peers: peers,
				mux:   &sync.Mutex{},
			}
			pBuf.RemovePeer(p)
			Expect(pBuf.peers).To(ConsistOf(expectedPeers))
		},

		Entry("remove peer from begin",
			[]Peer{
				{addr: "localhost", port: "55555"},
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "20000"},
			},
			Peer{addr: "localhost", port: "55555"},
			[]Peer{
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "20000"},
			}),
		Entry("remove peer from middle",
			[]Peer{
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "55555"},
				{addr: "localhost", port: "20000"},
			},
			Peer{addr: "localhost", port: "55555"},
			[]Peer{
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "20000"},
			}),
		Entry("remove peer from end",
			[]Peer{
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "20000"},
				{addr: "localhost", port: "55555"},
			},
			Peer{addr: "localhost", port: "55555"},
			[]Peer{
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "20000"},
			}),
		Entry("doesn't remove peer it it doesn't exist in buffer",
			[]Peer{
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "20000"},
			},
			Peer{addr: "localhost", port: "55555"},
			[]Peer{
				{addr: "localhost", port: "10000"},
				{addr: "localhost", port: "20000"},
			}),
	)
})
