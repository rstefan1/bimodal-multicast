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

package gossipserver

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/httpserver"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

const (
	timeout = time.Second
)

func suggestPort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

var _ = Describe("Gossip Server", func() {
	var (
		gossip       Gossip
		gossipPort   string
		mockPort     string
		gossipPeers  []peer.Peer
		mockPeers    []peer.Peer
		gossipMsgBuf *buffer.MessageBuffer
		mockMsgBuf   *buffer.MessageBuffer
		gossipCfg    Config
		httpCfg      httpserver.Config
		mockCfg      httpserver.Config
		gossipStop   chan struct{}
		httpStop     chan struct{}
		mockStop     chan struct{}
	)

	BeforeEach(func() {
		gossipPort = strconv.Itoa(suggestPort())
		mockPort = strconv.Itoa(suggestPort())

		gossipPeers = append(gossipPeers, peer.Peer{Addr: "localhost", Port: mockPort})
		mockPeers = append(mockPeers, peer.Peer{Addr: "localhost", Port: gossipPort})

		gossipMsgBuf = buffer.NewMessageBuffer()
		gossipMsgBuf.AddMessage(buffer.Message{
			ID:          fmt.Sprintf("%d", rand.Int31()),
			Msg:         fmt.Sprintf("%d", rand.Int31()),
			GossipCount: rand.Int(),
		})
		mockMsgBuf = buffer.NewMessageBuffer()

		gossipCfg = Config{
			Addr:    "localhost",
			Port:    gossipPort,
			PeerBuf: gossipPeers,
			MsgBuf:  gossipMsgBuf,
		}
		httpCfg = httpserver.Config{
			Addr:    "localhost",
			Port:    gossipPort,
			PeerBuf: gossipPeers,
			MsgBuf:  gossipMsgBuf,
		}
		mockCfg = httpserver.Config{
			Addr:    "localhost",
			Port:    mockPort,
			PeerBuf: mockPeers,
			MsgBuf:  mockMsgBuf,
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
			mockHTTPServer := httpserver.New(mockCfg)
			err := mockHTTPServer.Start(mockStop)
			Expect(err).To(Succeed())
		}()

		go func() {
			gossipHTTPServer := httpserver.New(httpCfg)
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
