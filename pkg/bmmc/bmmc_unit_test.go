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
	"fmt"
	"log"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/callback"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/gossip"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/round"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/server"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/testutil"
)

func newBMMC(addr, port string, peerBuffer *peer.PeerBuffer, msgBuffer *buffer.MessageBuffer) *Bmmc {
	gossipRound := round.NewGossipRound()

	cbCustomRegistry, err := callback.NewCustomRegistry(map[string]func(string, *log.Logger) (bool, error){})
	Expect(err).To(Succeed())
	cbDefaultRegistry, err := callback.NewDefaultRegistry()
	Expect(err).To(Succeed())

	server := server.New(server.Config{
		Addr:             addr,
		Port:             port,
		PeerBuf:          peerBuffer,
		MsgBuf:           msgBuffer,
		GossipRound:      gossipRound,
		Logger:           log.New(os.Stdout, "", 0),
		CustomCallbacks:  cbCustomRegistry,
		DefaultCallbacks: cbDefaultRegistry,
	})

	gossiper := gossip.New(gossip.Config{
		Addr:        addr,
		Port:        port,
		PeerBuf:     peerBuffer,
		MsgBuf:      msgBuffer,
		Beta:        defaultBeta,
		GossipRound: gossipRound,
		Logger:      log.New(os.Stdout, "", 0),
	})

	return &Bmmc{
		peerBuffer:   peerBuffer,
		msgBuffer:    msgBuffer,
		gossipRound:  gossipRound,
		httpServer:   server,
		gossipServer: gossiper,
		stop:         make(chan struct{}),
	}
}

var _ = Describe("BMMC", func() {
	var (
		addr1, addr2 string
		port1, port2 string
		bmmc1, bmmc2 *Bmmc
		peerBuffer1  *peer.PeerBuffer
		peerBuffer2  *peer.PeerBuffer
		msgBuffer1   *buffer.MessageBuffer
		msgBuffer2   *buffer.MessageBuffer
	)

	BeforeEach(func() {
		addr1 = "localhost"
		port1 = testutil.SuggestPort()
		addr2 = "localhost"
		port2 = testutil.SuggestPort()

		peerBuffer1 = peer.NewPeerBuffer()
		Expect(peerBuffer1.AddPeer(peer.NewPeer(addr2, port2))).To(BeTrue())
		msgBuffer1 = buffer.NewMessageBuffer()
		bmmc1 = newBMMC(addr1, port1, peerBuffer1, msgBuffer1)
		Expect(bmmc1.Start()).To(Succeed())

		peerBuffer2 = peer.NewPeerBuffer()
		Expect(peerBuffer2.AddPeer(peer.NewPeer(addr1, port1))).To(BeTrue())
		msgBuffer2 = buffer.NewMessageBuffer()
		bmmc2 = newBMMC(addr2, port2, peerBuffer2, msgBuffer2)
		Expect(bmmc2.Start()).To(Succeed())

		// wait after starting gossipers and servers
		time.Sleep(time.Millisecond)
	})

	AfterEach(func() {
		bmmc1.Stop()
		bmmc2.Stop()
	})

	When("add new message", func() {
		It("adds given message in message buffer", func() {
			msg := "awesome-message"
			expectedBuffer := []string{msg}

			bmmc1.AddMessage(msg, callback.NOCALLBACK)

			Eventually(func() []string {
				return bmmc1.GetMessages()
			}).Should(ConsistOf(expectedBuffer))
		})
	})

	When("add new peer", func() {
		It("adds given peer in peers list", func() {
			newAddr := "localhost"
			newPort := "49999"
			expectedBuffer := []string{
				fmt.Sprintf("%s/%s", newAddr, newPort),
				fmt.Sprintf("%s/%s", addr1, port1),
			}

			bmmc1.AddPeer(newAddr, newPort)

			Eventually(func() []string {
				return peerBuffer2.GetPeers()
			}).Should(ConsistOf(expectedBuffer))
		})
	})

	When("remove a peer", func() {
		It("removes given peer from peers list", func() {
			newAddr := "localhost"
			newPort := "49999"
			peerBuffer2.AddPeer(peer.NewPeer(newAddr, newPort))
			expectedBuffer := []string{
				fmt.Sprintf("%s/%s", addr1, port1),
			}

			bmmc1.RemovePeer(newAddr, newPort)

			Eventually(func() []string {
				return peerBuffer2.GetPeers()
			}).Should(ConsistOf(expectedBuffer))
		})
	})
})
