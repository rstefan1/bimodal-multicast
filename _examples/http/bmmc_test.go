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

package main

import (
	"errors"
	"fmt"
	"log/slog"
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

func fakeRegistry(cbType string, e error) map[string]func(any, *slog.Logger) error {
	return map[string]func(any, *slog.Logger) error{
		cbType: func(_ any, _ *slog.Logger) error {
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

func newBMMC(p Peer, callbacksRegistry map[string]func(any, *slog.Logger) error, logger *slog.Logger) *bmmc.BMMC {
	cbRegistry := map[string]func(any, *slog.Logger) error{}
	maps.Copy(cbRegistry, callbacksRegistry)

	b, err := bmmc.New(&bmmc.Config{
		Host:       p,
		Callbacks:  cbRegistry,
		BufferSize: 32,
		Logger:     logger,
	})
	Expect(err).ToNot(HaveOccurred())

	return b
}

var _ = Describe("BMMC with HTTP Server", func() {
	var bmmcLog, srvLog *slog.Logger

	BeforeEach(func() {
		bmmcLog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})).With("component", "bmmc_server")
		srvLog = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})).With("component", "http_server")
	})

	When("system has 2 nodes", func() {
		var (
			addr1, addr2 string
			port1, port2 string
			host1, host2 Peer
		)

		BeforeEach(func() {
			var err error

			addr1 = "localhost"
			port1 = suggestPort()

			addr2 = "localhost"
			port2 = suggestPort()

			host1, err = NewPeer(addr1, port1, &http.Client{})
			Expect(err).ToNot(HaveOccurred())

			host2, err = NewPeer(addr2, port2, &http.Client{})
			Expect(err).ToNot(HaveOccurred())
		})

		When("callback returns error", func() {
			var (
				bmmc1, bmmc2       *bmmc.BMMC
				srv1, srv2         *Server
				stopSrv1, stopSrv2 chan struct{}
			)

			BeforeEach(func() {
				stopSrv1 = make(chan struct{})
				stopSrv2 = make(chan struct{})

				cbRegistry := fakeRegistry("my-callback", errors.New("invalid-callback")) //nolint: goerr113

				bmmc1 = newBMMC(host1, cbRegistry, bmmcLog)
				bmmc2 = newBMMC(host2, cbRegistry, bmmcLog)

				srv1 = NewServer(bmmc1, addr1, port1, srvLog)
				srv2 = NewServer(bmmc2, addr2, port2, srvLog)

				Expect(bmmc1.Start()).To(Succeed())
				Expect(bmmc2.Start()).To(Succeed())

				Expect(srv1.Start(stopSrv1, srvLog)).To(Succeed())
				Expect(srv2.Start(stopSrv2, srvLog)).To(Succeed())

				Expect(bmmc1.AddPeer(host2.String())).To(Succeed())
				Expect(bmmc2.AddPeer(host1.String())).To(Succeed())
			})

			AfterEach(func() {
				// stop bmmcs
				bmmc1.Stop()
				bmmc2.Stop()

				// stop http servers
				close(stopSrv1)
				close(stopSrv2)
			})

			It("sync buffers", func() {
				// Add a message in first node.
				// Both nodes must have this message.
				Expect(bmmc1.AddMessage("awesome-first-message", "my-callback")).To(Succeed())

				Eventually(getBufferFn(bmmc1)).Should(ConsistOf(
					[]string{
						"awesome-first-message",
					},
				))
				Eventually(getBufferFn(bmmc2)).Should(ConsistOf(
					[]string{
						"awesome-first-message",
					},
				))

				// Add a message in second node.
				// Both nodes must have this message.
				Expect(bmmc2.AddMessage("awesome-second-message", "my-callback")).To(Succeed())

				Eventually(getBufferFn(bmmc1)).Should(ConsistOf(
					[]string{
						"awesome-first-message",
						"awesome-second-message",
					},
				))
				Eventually(getBufferFn(bmmc2)).Should(ConsistOf(
					[]string{
						"awesome-first-message",
						"awesome-second-message",
					},
				))
			})
		})

		When("callback doesn't return error", func() {
			var (
				bmmc1, bmmc2       *bmmc.BMMC
				srv1, srv2         *Server
				stopSrv1, stopSrv2 chan struct{}
			)

			BeforeEach(func() {
				stopSrv1 = make(chan struct{})
				stopSrv2 = make(chan struct{})

				cbRegistry := fakeRegistry("my-callback", nil)

				bmmc1 = newBMMC(host1, cbRegistry, bmmcLog)
				bmmc2 = newBMMC(host2, cbRegistry, bmmcLog)

				srv1 = NewServer(bmmc1, addr1, port1, srvLog)
				srv2 = NewServer(bmmc2, addr2, port2, srvLog)

				Expect(bmmc1.Start()).To(Succeed())
				Expect(bmmc2.Start()).To(Succeed())

				Expect(srv1.Start(stopSrv1, srvLog)).To(Succeed())
				Expect(srv2.Start(stopSrv2, srvLog)).To(Succeed())

				Expect(bmmc1.AddPeer(host2.String())).To(Succeed())
				Expect(bmmc2.AddPeer(host1.String())).To(Succeed())
			})

			AfterEach(func() {
				// stop bmmcs
				bmmc1.Stop()
				bmmc2.Stop()

				// stop http servers
				close(stopSrv1)
				close(stopSrv2)
			})

			It("sync buffers", func() {
				// Add a message in first node.
				// Both nodes must have this message.
				Expect(bmmc1.AddMessage("first-message", "my-callback")).To(Succeed())

				Eventually(getBufferFn(bmmc1)).Should(ConsistOf(
					[]string{
						"first-message",
					},
				))
				Eventually(getBufferFn(bmmc2)).Should(ConsistOf(
					[]string{
						"first-message",
					},
				))

				// Add a message in second node.
				// Both nodes must have this message.
				Expect(bmmc2.AddMessage("second-message", "my-callback")).To(Succeed())

				Eventually(getBufferFn(bmmc1)).Should(ConsistOf(
					[]string{
						"first-message",
						"second-message",
					},
				))
				Eventually(getBufferFn(bmmc2)).Should(ConsistOf(
					[]string{
						"first-message",
						"second-message",
					},
				))
			})
		})
	})

	When("system has ten nodes", func() {
		const nodesLen = 10

		var (
			bmmcs    [nodesLen]*bmmc.BMMC
			srvs     [nodesLen]*Server
			hosts    [nodesLen]Peer
			stopSrvs [nodesLen]chan struct{}

			err error
		)

		BeforeEach(func() {
			for i := 0; i < nodesLen; i++ {
				hosts[i], err = NewPeer("localhost", suggestPort(), &http.Client{})
				Expect(err).ToNot(HaveOccurred())

				bmmcs[i] = newBMMC(hosts[i], map[string]func(any, *slog.Logger) error{}, bmmcLog)
				Expect(bmmcs[i].Start()).To(Succeed())

				srvs[i] = NewServer(bmmcs[i], hosts[i].Addr, hosts[i].Port, srvLog)
				stopSrvs[i] = make(chan struct{})
				Expect(srvs[i].Start(stopSrvs[i], srvLog)).To(Succeed())
			}

			// add peers to the first bmmc node
			for i := 1; i < nodesLen; i++ {
				Expect(bmmcs[0].AddPeer(hosts[i].String())).To(Succeed())
			}

			// add the first bmmc node to the second bmmc node
			Expect(bmmcs[1].AddPeer(hosts[0].String())).To(Succeed())
		})

		AfterEach(func() {
			for i := 0; i < nodesLen; i++ {
				bmmcs[i].Stop()
				close(stopSrvs[i])
			}
		})

		When("one node has a message", func() {
			var expectedBuf []any

			BeforeEach(func() {
				msg := "another-message"
				expectedBuf = []any{msg}

				// select a random node to send the new message
				randomNode := rand.Intn(nodesLen) //nolint: gosec
				Expect(bmmcs[randomNode].AddMessage(msg, bmmc.NOCALLBACK)).To(Succeed())
			})

			It("sync all nodes with the new message", func() {
				for i := 0; i < nodesLen; i++ {
					Eventually(getBufferFn(bmmcs[i]), time.Second).Should(ConsistOf(expectedBuf...))
				}
			})
		})

		When("one node has more messages", func() {
			var expectedBuf []any

			BeforeEach(func() {
				expectedBuf = []any{}

				randomNode := rand.Intn(nodesLen) //nolint: gosec

				for i := 0; i < 3; i++ {
					msg := fmt.Sprint("new-message-$d", rand.Int31()) //nolint: gosec
					expectedBuf = append(expectedBuf, msg)

					Expect(bmmcs[randomNode].AddMessage(msg, bmmc.NOCALLBACK)).To(Succeed())
				}
			})

			It("sync all nodes with all messages", func() {
				for i := 0; i < nodesLen; i++ {
					Eventually(getBufferFn(bmmcs[i]), time.Second).Should(
						ConsistOf(anySliceToAnyString(expectedBuf)))
				}
			})
		})

		When("three nodes have three different messages", func() {
			var expectedBuf []any

			BeforeEach(func() {
				expectedBuf = []any{}

				randomNodes := []int{2, 4, 6}

				for _, randomNode := range randomNodes {
					msg := fmt.Sprintf("new-message-%d", rand.Int31()) //nolint: gosec
					expectedBuf = append(expectedBuf, msg)

					Expect(bmmcs[randomNode].AddMessage(msg, bmmc.NOCALLBACK)).To(Succeed())

				}
			})

			It("sync all nodes with all messages", func() {
				for i := 0; i < nodesLen; i++ {
					Eventually(getBufferFn(bmmcs[i]), time.Second).Should(
						ConsistOf(anySliceToAnyString(expectedBuf)))
				}
			})
		})
	})
})
