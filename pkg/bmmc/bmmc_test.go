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

package bmmc_test

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/callback"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rstefan1/bimodal-multicast/pkg/bmmc"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

func getBuffer(node *bmmc.BMMC) []string {
	buf := node.GetMessages()

	sbuf := make([]string, len(buf))
	for i, v := range buf {
		sbuf[i] = fmt.Sprint(v)
	}

	return sbuf
}

func getBufferFn(node *bmmc.BMMC) func() []string {
	return func() []string {
		return getBuffer(node)
	}
}

func interfaceToString(b []interface{}) []string {
	s := make([]string, len(b))
	for i, v := range b {
		s[i] = fmt.Sprint(v)
	}

	return s
}

func fakeRegistry(cbType string, e error) map[string]func(interface{}, *log.Logger) error {
	return map[string]func(interface{}, *log.Logger) error{
		cbType: func(_ interface{}, _ *log.Logger) error {
			return e
		},
	}
}

// suggestPort suggests an unused port.
func suggestPort() string {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	Expect(err).ToNot(HaveOccurred())

	l, err := net.ListenTCP("tcp", addr)
	Expect(err).ToNot(HaveOccurred())

	defer l.Close() //nolint: errcheck

	addr, convOk := l.Addr().(*net.TCPAddr)
	Expect(convOk).To(BeTrue())

	return strconv.Itoa(addr.Port)
}

func newBMMC(p peer.Peer, cbCustomRegistry map[string]func(interface{}, *log.Logger) error) *bmmc.BMMC {
	b, err := bmmc.New(&bmmc.Config{
		Host:       p,
		Callbacks:  cbCustomRegistry,
		BufferSize: 32,
	})
	Expect(err).ToNot(HaveOccurred())

	return b
}

var _ = Describe("BMMC", func() {
	DescribeTable("protocol",
		func(
			cbCustomRegistry map[string]func(interface{}, *log.Logger) error,
			msg string,
			callbackType string,
			expectedBuf []string,
		) {
			peer1, err := peer.NewHTTPPeer("localhost", suggestPort(), &http.Client{})
			Expect(err).ToNot(HaveOccurred())
			peer2, err := peer.NewHTTPPeer("localhost", suggestPort(), &http.Client{})
			Expect(err).ToNot(HaveOccurred())

			node1 := newBMMC(peer1, cbCustomRegistry)
			node2 := newBMMC(peer2, cbCustomRegistry)

			Expect(node1.Start()).To(Succeed())
			Expect(node2.Start()).To(Succeed())

			Expect(node1.AddPeer(peer2.String())).To(Succeed())
			Expect(node2.AddPeer(peer1.String())).To(Succeed())

			// message for adding peers in buffer
			extraMsgBuffer := []string{
				peer1.String(),
				peer2.String(),
			}
			expectedBuf = append(expectedBuf, extraMsgBuffer...)

			// Add a message in first node.
			// Both nodes must have this message.
			Expect(node1.AddMessage(msg, callbackType)).To(Succeed())

			Eventually(getBufferFn(node1)).Should(ConsistOf(expectedBuf))
			Eventually(getBufferFn(node2)).Should(ConsistOf(expectedBuf))
		},
		Entry("sync buffers if callback returns error",
			fakeRegistry("my-callback", errors.New("invalid-callback")), //nolint: goerr113
			"awesome-message",
			"my-callback",
			[]string{"awesome-message"},
		),
		Entry("sync buffers if callback doesn't return error",
			fakeRegistry("my-callback", nil),
			"awesome-message",
			"my-callback",
			[]string{"awesome-message"},
		),
	)

	When("system has ten nodes", func() {
		const nodesLen = 10

		var (
			nodes          [nodesLen]*bmmc.BMMC
			hosts          [nodesLen]peer.Peer
			extraMsgBuffer []interface{}
			expectedBuf    []interface{}

			err error
		)

		BeforeEach(func() {
			extraMsgBuffer = make([]interface{}, nodesLen)
			expectedBuf = []interface{}{}

			for i := 0; i < nodesLen; i++ {

				hosts[i], err = peer.NewHTTPPeer("localhost", suggestPort(), &http.Client{})
				Expect(err).ToNot(HaveOccurred())

				extraMsgBuffer[i] = hosts[i].String()
			}

			// create a protocol for each node, and start it
			for i := 0; i < nodesLen; i++ {
				nodes[i] = newBMMC(hosts[i], map[string]func(interface{}, *log.Logger) error{})
				Expect(nodes[i].Start()).To(Succeed())
			}

			// add peers
			for i := 1; i < nodesLen; i++ {
				Expect(nodes[0].AddPeer(hosts[i].String())).To(Succeed())
			}
			Expect(nodes[1].AddPeer(hosts[0].String())).To(Succeed())

			for i := range nodes {
				Eventually(getBufferFn(nodes[i]), time.Second*5).Should(ConsistOf(extraMsgBuffer...))
			}
		})

		AfterEach(func() {
			for i := range nodes {
				nodes[i].Stop()
			}
		})

		When("one node has an message", func() {
			BeforeEach(func() {
				msg := "another-awesome-message"
				expectedBuf = append(expectedBuf, msg)

				randomNode := rand.Intn(nodesLen) //nolint: gosec
				err := nodes[randomNode].AddMessage(msg, callback.NOCALLBACK)
				Expect(err).ToNot(HaveOccurred())
				Eventually(getBufferFn(nodes[randomNode]), time.Second).Should(ConsistOf(append(expectedBuf, extraMsgBuffer...)...))
			})

			It("sync all nodes with the message", func() {
				for i := range nodes {
					Eventually(getBufferFn(nodes[i]), time.Second).Should(ConsistOf(append(expectedBuf, extraMsgBuffer...)...))
				}
			})
		})

		When("one node has more messages", func() {
			BeforeEach(func() {
				randomNode := rand.Intn(nodesLen) //nolint: gosec
				for i := 0; i < 3; i++ {
					msg := i
					expectedBuf = append(expectedBuf, msg)

					err := nodes[randomNode].AddMessage(msg, callback.NOCALLBACK)
					Expect(err).ToNot(HaveOccurred())
				}

				Eventually(getBufferFn(nodes[randomNode]), time.Second).Should(
					ConsistOf(interfaceToString(append(expectedBuf, extraMsgBuffer...))))
			})

			It("sync all nodes with all messages", func() {
				for i := range nodes {
					Eventually(getBufferFn(nodes[i]), time.Second).Should(
						ConsistOf(interfaceToString(append(expectedBuf, extraMsgBuffer...))))
				}
			})
		})

		When("three nodes have three different messages", func() {
			BeforeEach(func() {
				randomNodes := [3]int{2, 4, 6}

				for i := 0; i < 3; i++ {
					msg := i
					expectedBuf = append(expectedBuf, msg)

					err := nodes[randomNodes[i]].AddMessage(msg, callback.NOCALLBACK)
					Expect(err).ToNot(HaveOccurred())
				}
			})

			It("sync all nodes with all messages", func() {
				for i := range nodes {
					Eventually(getBufferFn(nodes[i]), time.Second).Should(
						ConsistOf(interfaceToString(append(expectedBuf, extraMsgBuffer...))))
				}
			})
		})
	})
})
