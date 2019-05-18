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

package gossip

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/callback"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/round"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/server"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/testutil"
)

const (
	timeout = time.Second * 2
)

func getFakeEmptyCallbackRegistry() *callback.CustomRegistry {
	cb, err := callback.NewCustomRegistry(map[string]func(string, *log.Logger) (bool, error){})
	Expect(err).To(Succeed())
	return cb
}

var _ = Describe("Gossiper", func() {
	var (
		gossip       *Gossiper
		gossipPort   string
		mockPort     string
		gossipPeers  *peer.PeerBuffer
		mockPeers    *peer.PeerBuffer
		gossipMsgBuf *buffer.MessageBuffer
		mockMsgBuf   *buffer.MessageBuffer
		gossipRound  *round.GossipRound
		mockRound    *round.GossipRound
		gossipCfg    Config
		httpCfg      server.Config
		mockCfg      server.Config
		gossipStop   chan struct{}
		httpStop     chan struct{}
		mockStop     chan struct{}
	)

	BeforeEach(func() {
		gossipRound = round.NewGossipRound()
		mockRound = round.NewGossipRound()

		gossipPort = testutil.SuggestPort()
		mockPort = testutil.SuggestPort()

		mockPeer := peer.NewPeer("localhost", mockPort)
		gossipPeer := peer.NewPeer("localhost", gossipPort)

		mockPeers = peer.NewPeerBuffer()
		ok := mockPeers.AddPeer(gossipPeer)
		Expect(ok).To(BeTrue())

		gossipPeers = peer.NewPeerBuffer()
		ok = gossipPeers.AddPeer(mockPeer)
		Expect(ok).To(BeTrue())

		gossipMsgBuf = buffer.NewMessageBuffer()
		gossipMsgBuf.AddMessage(buffer.NewMessage(
			fmt.Sprintf("%d", rand.Int31()),
			callback.NOCALLBACK,
		))
		mockMsgBuf = buffer.NewMessageBuffer()

		cbDefaultRegistry, err := callback.NewDefaultRegistry()
		Expect(err).To(Succeed())

		gossipCfg = Config{
			Addr:        "localhost",
			Port:        gossipPort,
			PeerBuf:     gossipPeers,
			MsgBuf:      gossipMsgBuf,
			GossipRound: gossipRound,
			Logger:      log.New(os.Stdout, "", 0),
		}
		httpCfg = server.Config{
			Addr:             "localhost",
			Port:             gossipPort,
			PeerBuf:          gossipPeers,
			MsgBuf:           gossipMsgBuf,
			GossipRound:      gossipRound,
			CustomCallbacks:  getFakeEmptyCallbackRegistry(),
			DefaultCallbacks: cbDefaultRegistry,
		}
		mockCfg = server.Config{
			Addr:             "localhost",
			Port:             mockPort,
			PeerBuf:          mockPeers,
			MsgBuf:           mockMsgBuf,
			GossipRound:      mockRound,
			CustomCallbacks:  getFakeEmptyCallbackRegistry(),
			DefaultCallbacks: cbDefaultRegistry,
		}

		gossipStop = make(chan struct{})
		httpStop = make(chan struct{})
		mockStop = make(chan struct{})

		gossip = New(gossipCfg)
	})

	AfterEach(func() {
		close(gossipStop)
		close(mockStop)
	})

	It("synchronize nodes with missing messages", func() {
		go func() {
			mockHTTPServer := server.New(mockCfg)
			err := mockHTTPServer.Start(mockStop)
			Expect(err).To(Succeed())
		}()

		go func() {
			gossipHTTPServer := server.New(httpCfg)
			err := gossipHTTPServer.Start(httpStop)
			Expect(err).To(Succeed())
		}()

		// wait for starting http servers
		time.Sleep(time.Second)

		go func() {
			gossip.Start(gossipStop)
		}()

		Eventually(func() bool {
			return gossipMsgBuf.SameMessages(mockMsgBuf)
		}, timeout).Should(Equal(true))
	})
})
