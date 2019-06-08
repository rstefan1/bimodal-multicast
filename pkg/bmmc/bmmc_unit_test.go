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

const localhost = "localhost"

func newBMMC(addr, port string, peerBuffer *peer.Buffer, msgBuffer *buffer.MessageBuffer) *Bmmc {
	gossipRound := round.NewGossipRound()

	cbCustomRegistry, err := callback.NewCustomRegistry(map[string]func(interface{}, *log.Logger) error{})
	Expect(err).To(Succeed())
	cbDefaultRegistry, err := callback.NewDefaultRegistry()
	Expect(err).To(Succeed())

	logger := log.New(os.Stdout, "", 0)

	serverCfg := server.Config{
		Addr:             addr,
		Port:             port,
		PeerBuf:          peerBuffer,
		MsgBuf:           msgBuffer,
		GossipRound:      gossipRound,
		Logger:           logger,
		CustomCallbacks:  cbCustomRegistry,
		DefaultCallbacks: cbDefaultRegistry,
	}
	server := server.New(serverCfg)

	gossiperCfg := gossip.Config{
		Addr:        addr,
		Port:        port,
		PeerBuf:     peerBuffer,
		MsgBuf:      msgBuffer,
		Beta:        defaultBeta,
		GossipRound: gossipRound,
		Logger:      logger,
	}
	gossiper := gossip.New(gossiperCfg)

	return &Bmmc{
		peerBuffer:  peerBuffer,
		msgBuffer:   msgBuffer,
		gossipRound: gossipRound,
		server:      server,
		gossiper:    gossiper,
		serverCfg:   serverCfg,
		gossiperCfg: gossiperCfg,
		logger:      logger,
		stop:        make(chan struct{}),
	}
}

var _ = Describe("BMMC", func() {
	var (
		addr1, addr2 string
		port1, port2 string
		bmmc1, bmmc2 *Bmmc
		peerBuffer1  *peer.Buffer
		peerBuffer2  *peer.Buffer
		msgBuffer1   *buffer.MessageBuffer
		msgBuffer2   *buffer.MessageBuffer
	)

	BeforeEach(func() {
		addr1 = localhost
		port1 = testutil.SuggestPort()
		addr2 = localhost
		port2 = testutil.SuggestPort()

		peerBuffer1 = peer.NewPeerBuffer()
		Expect(peerBuffer1.AddPeer(peer.NewPeer(addr2, port2))).To(Succeed())
		msgBuffer1 = buffer.NewMessageBuffer()
		bmmc1 = newBMMC(addr1, port1, peerBuffer1, msgBuffer1)
		Expect(bmmc1.Start()).To(Succeed())

		peerBuffer2 = peer.NewPeerBuffer()
		Expect(peerBuffer2.AddPeer(peer.NewPeer(addr1, port1))).To(Succeed())
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

			err := bmmc1.AddMessage(msg, callback.NOCALLBACK)
			Expect(err).To(BeNil())

			Eventually(func() []string {
				buf := bmmc1.GetMessages()

				sbuf := make([]string, len(buf))
				for i, v := range buf {
					sbuf[i] = fmt.Sprint(v)
				}
				return sbuf
			}).Should(ConsistOf(expectedBuffer))
		})
	})

	When("add new peer", func() {
		It("adds given peer in peers list", func() {
			newAddr := localhost
			newPort := "49999"

			expectedBuffer1 := []string{
				fmt.Sprintf("%s/%s", newAddr, newPort),
				fmt.Sprintf("%s/%s", addr2, port2),
			}
			expectedBuffer2 := []string{
				fmt.Sprintf("%s/%s", newAddr, newPort),
				fmt.Sprintf("%s/%s", addr1, port1),
			}

			Expect(bmmc1.AddPeer(newAddr, newPort)).To(Succeed())

			Eventually(func() []string {
				return peerBuffer1.GetPeers()
			}).Should(ConsistOf(expectedBuffer1))
			Eventually(func() []string {
				return peerBuffer2.GetPeers()
			}).Should(ConsistOf(expectedBuffer2))
		})
	})

	When("remove a peer", func() {
		It("removes given peer from peers list", func() {
			newAddr := localhost
			newPort := "49999"
			Expect(peerBuffer2.AddPeer(peer.NewPeer(newAddr, newPort))).To(Succeed())
			expectedBuffer := []string{
				fmt.Sprintf("%s/%s", addr1, port1),
			}

			Expect(bmmc1.RemovePeer(newAddr, newPort)).To(Succeed())

			Eventually(func() []string {
				return peerBuffer2.GetPeers()
			}).Should(ConsistOf(expectedBuffer))
		})
	})
})
