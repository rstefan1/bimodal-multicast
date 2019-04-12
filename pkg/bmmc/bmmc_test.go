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
	"fmt"
	"math/rand"
	"net"
	"sort"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rstefan1/bimodal-multicast/pkg/bmmc"
	"github.com/rstefan1/bimodal-multicast/pkg/peer"
)

func suggestPort() string {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
}

// getSortedBuffer is a helper func that returns a sorted buffer
func getSortedBuffer(node *bmmc.Bmmc) func() []string {
	return func() []string {
		buf := node.GetMessages()
		sort.Strings(buf)
		return buf
	}
}

var _ = Describe("BMMC", func() {
	When("system has two nodes and one node has a message in buffer", func() {
		var (
			node1       *bmmc.Bmmc
			node2       *bmmc.Bmmc
			msg         string
			expectedBuf []string
		)

		BeforeEach(func() {
			port1 := suggestPort()
			port2 := suggestPort()

			peers := []peer.Peer{
				{
					Addr: "localhost",
					Port: port1,
				},
				{
					Addr: "localhost",
					Port: port2,
				},
			}

			node1 = bmmc.New(bmmc.Config{
				Addr:  "localhost",
				Port:  port1,
				Peers: peers,
				Beta:  0.5,
			})
			node2 = bmmc.New(bmmc.Config{
				Addr:  "localhost",
				Port:  port2,
				Peers: peers,
				Beta:  0.5,
			})

			Expect(node1.Start())
			Expect(node2.Start())

			// add a message in first node
			msg = "awesome-message"
			expectedBuf = append(expectedBuf, msg)
			node1.AddMessage(msg)
			Eventually(getSortedBuffer(node1), time.Second).Should(SatisfyAll(
				HaveLen(1),
				Equal(expectedBuf),
			))
		})

		AfterEach(func() {
			node1.Stop()
			node2.Stop()
		})

		It("sync buffers with the message", func() {
			Eventually(getSortedBuffer(node2), time.Second).Should(SatisfyAll(
				HaveLen(1),
				Equal(expectedBuf),
			))
		})
	})

	When("system has ten nodes", func() {
		const len = 10
		var nodes [len]*bmmc.Bmmc

		BeforeEach(func() {
			var (
				ports [len]string
				peers []peer.Peer
			)

			for i := 0; i < len; i++ {
				ports[i] = suggestPort()
				peers = append(peers, peer.Peer{
					Addr: "localhost",
					Port: ports[i],
				})
			}

			for i := 0; i < len; i++ {
				nodes[i] = bmmc.New(bmmc.Config{
					Addr:  "localhost",
					Port:  ports[i],
					Peers: peers,
				})
			}
		})

		AfterEach(func() {
			for i := range nodes {
				nodes[i].Stop()
			}
		})

		When("one node has an message", func() {
			var expectedBuf []string

			BeforeEach(func() {
				msg := "another-awesome-message"
				expectedBuf = append(expectedBuf, msg)
				randomNode := rand.Intn(len)
				nodes[randomNode].AddMessage(msg)
				Eventually(getSortedBuffer(nodes[randomNode]), time.Second).Should(SatisfyAll(
					HaveLen(1),
					Equal(expectedBuf),
				))

				// start protocol for all nodes
				for i := 0; i < len; i++ {
					Expect(nodes[i].Start())
				}
			})

			It("sync all nodes with the message", func() {
				for i := range nodes {
					Eventually(getSortedBuffer(nodes[i]), time.Second).Should(SatisfyAll(
						HaveLen(1),
						Equal(expectedBuf),
					))
				}
			})
		})

		When("one node has more messages", func() {
			var expectedBuf []string

			BeforeEach(func() {
				randomNode := rand.Intn(len)
				for i := 0; i < 3; i++ {
					msg := fmt.Sprintf("awesome-message-%d", i)
					expectedBuf = append(expectedBuf, msg)
					nodes[randomNode].AddMessage(msg)
				}

				Eventually(getSortedBuffer(nodes[randomNode]), time.Second).Should(SatisfyAll(
					HaveLen(3),
					Equal(expectedBuf),
				))

				// start protocol for all nodes
				for i := 0; i < len; i++ {
					Expect(nodes[i].Start())
				}
			})

			It("sync all nodes with all messages", func() {
				for i := range nodes {
					Eventually(getSortedBuffer(nodes[i]), time.Second*2).Should(SatisfyAll(
						HaveLen(3),
						Equal(expectedBuf),
					))
				}
			})
		})

		When("three nodes have three different messages", func() {
			var expectedBuf []string

			BeforeEach(func() {
				for i := 0; i < 3; i++ {
					randomNode := rand.Intn(len)
					msg := fmt.Sprintf("awesome-message-%d", i)
					expectedBuf = append(expectedBuf, msg)
					nodes[randomNode].AddMessage(msg)

					Eventually(getSortedBuffer(nodes[randomNode]), time.Second).Should(SatisfyAll(
						HaveLen(1),
						Equal([]string{msg}),
					))
				}

				// start protocol for all nodes
				for i := 0; i < len; i++ {
					Expect(nodes[i].Start())
				}
			})

			It("sync all nodes with all messages", func() {
				for i := range nodes {
					Eventually(getSortedBuffer(nodes[i]), time.Second*3).Should(SatisfyAll(
						HaveLen(3),
						Equal(expectedBuf),
					))
				}
			})
		})
	})
})
