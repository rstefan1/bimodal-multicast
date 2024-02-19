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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

var _ = Describe("Gossiper", func() {
	Describe("computeGossipLen function", func() {
		var b *BMMC

		BeforeEach(func() {
			peerBuf := peer.NewPeerBuffer()
			Expect(peerBuf.AddPeer("localhost/19999")).To(Succeed())

			msgBuf := buffer.NewBuffer(25)

			msg, err := buffer.NewElement("awesome message", "awesome-callback")
			Expect(err).ToNot(HaveOccurred())

			Expect(msgBuf.Add(msg)).To(Succeed())

			b = &BMMC{
				peerBuffer:    peerBuf,
				messageBuffer: msgBuf,
				config: &Config{
					Beta: 0.5,
				},
			}
		})

		It("returns 0 if peerBuffer's legth is 0", func() {
			b.peerBuffer = peer.NewPeerBuffer()
			Expect(b.computeGossipLen()).To(Equal(0))
		})

		It("returns 0 if messageBuffer's length is 0", func() {
			b.messageBuffer = buffer.NewBuffer(25)
			Expect(b.computeGossipLen()).To(Equal(0))
		})

		It("returns 0 if beta is 0", func() {
			b.config.Beta = 0
			Expect(b.computeGossipLen()).To(Equal(0))
		})

		It("returns proper gossip len if no field are 0", func() {
			Expect(b.computeGossipLen()).To(Equal(int(b.config.Beta*float64(b.peerBuffer.Length())) + 1))
		})
	})
})
