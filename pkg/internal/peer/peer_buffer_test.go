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
	When("Length() func is called", func() {
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

	When("alreadyExists() func is called", func() {
		It("returns true if peer is at the beginning", func() {
			pBuf := &Buffer{
				peers: []string{
					"localhost/55555",
					"localhost/10000",
					"localhost/20000",
				},
				mux: &sync.Mutex{},
			}

			Expect(pBuf.alreadyExists("localhost/55555")).To(BeTrue())
		})

		It("returns true if peer is in the middle", func() {
			pBuf := &Buffer{
				peers: []string{
					"localhost/10000",
					"localhost/55555",
					"localhost/20000",
				},
				mux: &sync.Mutex{},
			}

			Expect(pBuf.alreadyExists("localhost/55555")).To(BeTrue())
		})

		It("returns true if peer is at the end", func() {
			pBuf := &Buffer{
				peers: []string{
					"localhost/10000",
					"localhost/20000",
					"localhost/55555",
				},
				mux: &sync.Mutex{},
			}

			Expect(pBuf.alreadyExists("localhost/55555")).To(BeTrue())
		})

		It("returns false if peer doesn't exist in buffer", func() {
			pBuf := &Buffer{
				peers: []string{
					"localhost/10000",
					"localhost/20000",
				},
				mux: &sync.Mutex{},
			}

			Expect(pBuf.alreadyExists("localhost/55555")).To(BeFalse())
		})
	})

	When("AddPeer() func is called", func() {
		It("adds peer in the peers buffer when it doesn't exist", func() {
			pBuf := &Buffer{
				peers: []string{
					"localhost/10000",
					"localhost/20000",
				},
				mux: &sync.Mutex{},
			}
			newPeer := "localhost/30000"
			expectedPeers := []string{
				"localhost/10000",
				"localhost/20000",
				"localhost/30000",
			}

			Expect(pBuf.AddPeer(newPeer)).To(Succeed())
			Expect(pBuf.peers).To(ConsistOf(expectedPeers))
		})

		It("doesn't add peer in the peers buffer when it already exists", func() {
			pBuf := &Buffer{
				peers: []string{
					"localhost/10000",
					"localhost/20000",
				},
				mux: &sync.Mutex{},
			}
			newPeer := "localhost/10000"
			expectedPeers := []string{
				"localhost/10000",
				"localhost/20000",
			}

			Expect(pBuf.AddPeer(newPeer)).NotTo(Succeed())
			Expect(pBuf.peers).To(ConsistOf(expectedPeers))
		})
	})

	When("GetPeers() func is called", func() {
		It("returns a slice of strings with peers", func() {
			pBuf := &Buffer{
				peers: []string{
					"localhost/10000",
					"localhost/20000",
					"localhost/30000",
				},
				mux: &sync.Mutex{},
			}
			expectedPeers := []string{
				"localhost/10000",
				"localhost/20000",
				"localhost/30000",
			}

			Expect(pBuf.GetPeers()).To(ConsistOf(expectedPeers))
		})
	})

	When("RemovePeer() func is called", func() {
		It("remove peer if it is the first element of buffer", func() {
			pBuf := &Buffer{
				peers: []string{
					"localhost/66666",
					"localhost/10000",
					"localhost/20000",
				},
				mux: &sync.Mutex{},
			}
			peerToRemove := "localhost/66666"
			expectedPeers := []string{
				"localhost/10000",
				"localhost/20000",
			}

			pBuf.RemovePeer(peerToRemove)
			Expect(pBuf.peers).To(ConsistOf(expectedPeers))
		})

		It("remove peer if it is in the middle of buffer", func() {
			pBuf := &Buffer{
				peers: []string{
					"localhost/10000",
					"localhost/77777",
					"localhost/20000",
				},
				mux: &sync.Mutex{},
			}
			peerToRemove := "localhost/77777"
			expectedPeers := []string{
				"localhost/10000",
				"localhost/20000",
			}

			pBuf.RemovePeer(peerToRemove)
			Expect(pBuf.peers).To(ConsistOf(expectedPeers))
		})

		It("remove peer if it is the last element of buffer", func() {
			pBuf := &Buffer{
				peers: []string{
					"localhost/10000",
					"localhost/20000",
					"localhost/88888",
				},
				mux: &sync.Mutex{},
			}
			peerToRemove := "localhost/88888"
			expectedPeers := []string{
				"localhost/10000",
				"localhost/20000",
			}

			pBuf.RemovePeer(peerToRemove)
			Expect(pBuf.peers).To(ConsistOf(expectedPeers))
		})

		It("doesn't remove peer if it doesn't exist in buffer", func() {
			pBuf := &Buffer{
				peers: []string{
					"localhost/10000",
					"localhost/20000",
				},
				mux: &sync.Mutex{},
			}
			peerToRemove := "localhost/99999"
			expectedPeers := []string{
				"localhost/10000",
				"localhost/20000",
			}

			pBuf.RemovePeer(peerToRemove)
			Expect(pBuf.peers).To(ConsistOf(expectedPeers))
		})
	})
})
