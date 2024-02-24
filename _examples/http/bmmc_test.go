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

package main

import (
	"errors"
	"fmt"
	"log"
	"maps"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rstefan1/bimodal-multicast/pkg/bmmc"
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

func anySliceToAnyString(b []any) []string {
	s := make([]string, len(b))
	for i, v := range b {
		s[i] = fmt.Sprint(v)
	}

	return s
}

func fakeRegistry(cbType string, e error) map[string]func(any, *log.Logger) error {
	return map[string]func(any, *log.Logger) error{
		cbType: func(_ any, _ *log.Logger) error {
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

func newBMMC(p Peer, callbacksRegistry map[string]func(any, *log.Logger) error) *bmmc.BMMC {
	cbRegistry := map[string]func(any, *log.Logger) error{}
	maps.Copy(cbRegistry, callbacksRegistry)

	b, err := bmmc.New(&bmmc.Config{
		Host:       p,
		Callbacks:  cbRegistry,
		BufferSize: 32,
	})
	Expect(err).ToNot(HaveOccurred())

	return b
}

var _ = Describe("BMMC", func() {
	testLog := log.New(os.Stdout, "", 0)

	DescribeTable("protocol",
		func(
			callbacksRegistry map[string]func(any, *log.Logger) error,
			msg string,
			callbackType string,
			expectedBuf []string,
		) {
			// server 1
			addr1 := "localhost"
			port1 := suggestPort()
			stop1 := make(chan struct{})

			peer1, err := NewPeer(addr1, port1, &http.Client{})
			Expect(err).ToNot(HaveOccurred())

			bmmc1 := newBMMC(peer1, callbacksRegistry)
			srv1 := NewServer(bmmc1, addr1, port1, testLog)

			// server 2
			addr2 := "localhost"
			port2 := suggestPort()
			stop2 := make(chan struct{})

			peer2, err := NewPeer(addr2, port2, &http.Client{})
			Expect(err).ToNot(HaveOccurred())

			bmmc2 := newBMMC(peer2, callbacksRegistry)
			srv2 := NewServer(bmmc2, addr2, port2, testLog)

			// start server 1 and server 2
			Expect(bmmc1.Start()).To(Succeed())
			Expect(srv1.Start(stop1, testLog)).To(Succeed())

			Expect(bmmc2.Start()).To(Succeed())
			Expect(srv2.Start(stop2, testLog)).To(Succeed())

			Expect(bmmc1.AddPeer(peer2.String())).To(Succeed())
			Expect(bmmc2.AddPeer(peer1.String())).To(Succeed())

			// Add a message in first node.
			// Both nodes must have this message.
			Expect(bmmc1.AddMessage(msg, callbackType)).To(Succeed())

			Eventually(getBufferFn(bmmc1)).Should(ConsistOf(expectedBuf))
			Eventually(getBufferFn(bmmc2)).Should(ConsistOf(expectedBuf))

			bmmc1.Stop()
			close(stop1)
			bmmc2.Stop()
			close(stop2)
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
			nodes       [nodesLen]*bmmc.BMMC
			srvs        [nodesLen]*Server
			peers       [nodesLen]Peer
			stops       [nodesLen]chan struct{}
			expectedBuf []any

			err error
		)

		BeforeEach(func() {
			expectedBuf = []any{}

			for i := 0; i < nodesLen; i++ {
				peers[i], err = NewPeer("localhost", suggestPort(), &http.Client{})
				Expect(err).ToNot(HaveOccurred())

				stops[i] = make(chan struct{})
			}

			// create a protocol for each node, and start it
			for i := 0; i < nodesLen; i++ {
				nodes[i] = newBMMC(peers[i], map[string]func(any, *log.Logger) error{})
				srvs[i] = NewServer(nodes[i], peers[i].Addr(), peers[i].Port(), testLog)

				Expect(nodes[i].Start()).To(Succeed())
				Expect(srvs[i].Start(stops[i], testLog)).To(Succeed())
			}

			// add peers
			for i := 1; i < nodesLen; i++ {
				Expect(nodes[0].AddPeer(peers[i].String())).To(Succeed())
			}
			Expect(nodes[1].AddPeer(peers[0].String())).To(Succeed())

			for i := 1; i < nodesLen; i++ {
				Eventually(getBufferFn(nodes[i]), time.Second).Should(ConsistOf(expectedBuf...))
			}
		})

		AfterEach(func() {
			for i := 0; i < nodesLen; i++ {
				nodes[i].Stop()
				close(stops[i])
			}
		})

		When("one node has an message", func() {
			BeforeEach(func() {
				msg := "another-awesome-message"
				expectedBuf = append(expectedBuf, msg)

				randomNode := rand.Intn(nodesLen) //nolint: gosec
				err := nodes[randomNode].AddMessage(msg, bmmc.NOCALLBACK)
				Expect(err).ToNot(HaveOccurred())
				Eventually(getBufferFn(nodes[randomNode]), time.Second).Should(ConsistOf(expectedBuf...))
			})

			It("sync all nodes with the message", func() {
				for i := 0; i < nodesLen; i++ {
					Eventually(getBufferFn(nodes[i]), time.Second).Should(ConsistOf(expectedBuf...))
				}
			})
		})

		When("one node has more messages", func() {
			BeforeEach(func() {
				randomNode := rand.Intn(nodesLen) //nolint: gosec
				for i := 0; i < 3; i++ {
					msg := i
					expectedBuf = append(expectedBuf, msg)

					err := nodes[randomNode].AddMessage(msg, bmmc.NOCALLBACK)
					Expect(err).ToNot(HaveOccurred())
				}

				Eventually(getBufferFn(nodes[randomNode]), time.Second).Should(
					ConsistOf(anySliceToAnyString(expectedBuf)))
			})

			It("sync all nodes with all messages", func() {
				for i := 0; i < nodesLen; i++ {
					Eventually(getBufferFn(nodes[i]), time.Second).Should(
						ConsistOf(anySliceToAnyString(expectedBuf)))
				}
			})
		})

		When("three nodes have three different messages", func() {
			BeforeEach(func() {
				randomNodes := [3]int{2, 4, 6}

				for i := 0; i < 3; i++ {
					msg := i
					expectedBuf = append(expectedBuf, msg)

					err := nodes[randomNodes[i]].AddMessage(msg, bmmc.NOCALLBACK)
					Expect(err).ToNot(HaveOccurred())
				}
			})

			It("sync all nodes with all messages", func() {
				for i := 0; i < nodesLen; i++ {
					Eventually(getBufferFn(nodes[i]), time.Second).Should(
						ConsistOf(anySliceToAnyString(expectedBuf)))
				}
			})
		})
	})
})
