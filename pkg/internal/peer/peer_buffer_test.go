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

package peer

import (
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTTP Peer Buffer Interface", func() {
	When("Length() is called", func() {
		It("returns the length of given buffer", func() {
			pBuf := &Buffer{
				peers: []string{
					"localhost/10000",
					"localhost/20000",
				},
				mux: &sync.Mutex{},
			}

			Expect(pBuf.Length()).To(Equal(2))
		})
	})

	DescribeTable("when alreadyExists() is called",
		func(peers []string, p string, expected bool) {
			pBuf := &Buffer{
				peers: peers,
				mux:   &sync.Mutex{},
			}
			Expect(pBuf.alreadyExists(p)).To(Equal(expected))
		},

		Entry("returns true if peer is at the beginning",
			[]string{
				"localhost/55555",
				"localhost/10000",
				"localhost/20000",
			},
			"localhost/55555",
			true,
		),

		Entry("returns true if peer is in the middle",
			[]string{
				"localhost/10000",
				"localhost/55555",
				"localhost/20000",
			},
			"localhost/55555",
			true,
		),

		Entry("returns true if peer is at the end",
			[]string{
				"localhost/10000",
				"localhost/20000",
				"localhost/55555",
			},
			"localhost/55555",
			true,
		),

		Entry("returns false if peer doesn't exist in buffer",
			[]string{
				"localhost/10000",
				"localhost/20000",
			},
			"localhost/55555",
			false,
		),
	)

	DescribeTable("when AddPeer() is called",
		func(peers []string, p string, expectError bool, expectedPeers []string) {
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
			[]string{
				"localhost/10000",
				"localhost/20000",
			},
			"localhost/30000",
			true,
			[]string{
				"localhost/10000",
				"localhost/20000",
				"localhost/30000",
			},
		),

		Entry("doens't add peer in the peers buffer when it already exists",
			[]string{
				"localhost/10000",
				"localhost/20000",
			},
			"localhost/10000",
			false,
			[]string{
				"localhost/10000",
				"localhost/20000",
			},
		),
	)

	When("GetPeers() is called", func() {
		It("returns a slice of strings with peers", func() {
			peers := []string{
				"localhost/10000",
				"localhost/20000",
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
		func(peers []string, p string, expectedPeers []string) {
			pBuf := &Buffer{
				peers: peers,
				mux:   &sync.Mutex{},
			}
			pBuf.RemovePeer(p)
			Expect(pBuf.peers).To(ConsistOf(expectedPeers))
		},

		Entry("remove peer from begin",
			[]string{
				"localhost/55555",
				"localhost/10000",
				"localhost/20000",
			},
			"localhost/55555",
			[]string{
				"localhost/10000",
				"localhost/20000",
			},
		),

		Entry("remove peer from middle",
			[]string{
				"localhost/10000",
				"localhost/55555",
				"localhost/20000",
			},
			"localhost/55555",
			[]string{
				"localhost/10000",
				"localhost/20000",
			},
		),

		Entry("remove peer from end",
			[]string{
				"localhost/10000",
				"localhost/20000",
				"localhost/55555",
			},
			"localhost/55555",
			[]string{
				"localhost/10000",
				"localhost/20000",
			},
		),

		Entry("doesn't remove peer if it doesn't exist in buffer",
			[]string{
				"localhost/10000",
				"localhost/20000",
			},
			"localhost/55555",
			[]string{
				"localhost/10000",
				"localhost/20000",
			},
		),
	)
})
