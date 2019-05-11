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
	"log"
	"math/rand"
	"os"
	"sort"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/rstefan1/bimodal-multicast/pkg/bmmc"
	"github.com/rstefan1/bimodal-multicast/pkg/callback"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/testutil"
	"github.com/rstefan1/bimodal-multicast/pkg/peer"
)

// getSortedBuffer is a helper func that returns a sorted buffer
func getSortedBuffer(node *bmmc.Bmmc) func() []string {
	return func() []string {
		buf := node.GetMessages()
		sort.Strings(buf)
		return buf
	}
}

func getFakeCallbackRegistry(cbType string, b bool, err error) *callback.Registry {
	r := callback.NewRegistry()
	_ = r.Register(
		cbType,
		func(msg string) (bool, error) {
			return b, err
		})
	return r
}

var _ = Describe("BMMC", func() {
	When("creates new protocol instance with broken config", func() {
		var (
			peers = []peer.Peer{
				{
					Addr: "localhost",
					Port: "1999",
				},
			}
		)

		It("returns error when address is empty", func() {
			_, err := bmmc.New(&bmmc.Config{
				Port:      "1999",
				Peers:     peers,
				Beta:      0.5,
				Logger:    log.New(os.Stdout, "", 0),
				Callbacks: callback.NewRegistry(),
			})
			Expect(err).To(Not(Succeed()))
		})

		It("returns error when port is empty", func() {
			_, err := bmmc.New(&bmmc.Config{
				Addr:      "localhost",
				Peers:     peers,
				Beta:      0.5,
				Logger:    log.New(os.Stdout, "", 0),
				Callbacks: callback.NewRegistry(),
			})
			Expect(err).To(Not(Succeed()))
		})

		It("returns error when callback registry is nil", func() {
			cfg := bmmc.Config{
				Addr:   "localhost",
				Port:   "1999",
				Beta:   0.5,
				Peers:  peers,
				Logger: log.New(os.Stdout, "", 0),
			}
			_, err := bmmc.New(&cfg)
			Expect(err).To(Not(Succeed()))
		})

		It("set default value for beta if it is empty", func() {
			cfg := bmmc.Config{
				Addr:      "localhost",
				Port:      "1999",
				Peers:     peers,
				Logger:    log.New(os.Stdout, "", 0),
				Callbacks: callback.NewRegistry(),
			}
			_, err := bmmc.New(&cfg)
			Expect(err).To(Succeed())
			Expect(cfg.Beta).To(Equal(0.3))
		})

		It("set default value for logger if it is empty", func() {
			cfg := bmmc.Config{
				Addr:      "localhost",
				Port:      "1999",
				Peers:     peers,
				Beta:      0.5,
				Callbacks: callback.NewRegistry(),
			}
			_, err := bmmc.New(&cfg)
			Expect(err).To(Succeed())
			Expect(cfg.Logger).To(Not(BeNil()))
		})
	})

	DescribeTable("when system has two nodes and one node has a message in buffer", func(cbRegistry *callback.Registry, msg, callbackType string, expectedBuf []string) {
		port1 := testutil.SuggestPort()
		port2 := testutil.SuggestPort()

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

		node1, err := bmmc.New(&bmmc.Config{
			Addr:      "localhost",
			Port:      port1,
			Peers:     peers,
			Beta:      0.5,
			Callbacks: cbRegistry,
		})
		Expect(err).To(Succeed())

		node2, err := bmmc.New(&bmmc.Config{
			Addr:      "localhost",
			Port:      port2,
			Peers:     peers,
			Beta:      0.5,
			Callbacks: cbRegistry,
		})
		Expect(err).To(Succeed())

		Expect(node1.Start())
		Expect(node2.Start())

		// add a message in first node
		node1.AddMessage(msg, callbackType)
		Eventually(getSortedBuffer(node1), time.Second).Should(Equal([]string{msg}))

		// Wait 1 second before checking the buffer.
		// In this second buffer needs to be updated.
		time.Sleep(time.Millisecond * 500)

		Expect(getSortedBuffer(node2)()).To(Equal(expectedBuf))
	},
		Entry("sync buffers with the message", callback.NewRegistry(), "awesome-message", callback.NOCALLBACK, []string{"awesome-message"}),
		Entry("doesn't sync buffers if callback returned error",
			getFakeCallbackRegistry("my-callback", true, fmt.Errorf("invalid-callback")),
			"awesome-message",
			"my-callback",
			[]string{}),
		Entry("doesn't sync buffers if callback returned false",
			getFakeCallbackRegistry("my-callback", false, nil),
			"awesome-message",
			"my-callback",
			[]string{}),
	)

	When("system has ten nodes", func() {
		const len = 10
		var nodes [len]*bmmc.Bmmc

		BeforeEach(func() {
			var (
				ports [len]string
				peers []peer.Peer
				err   error
			)

			for i := 0; i < len; i++ {
				ports[i] = testutil.SuggestPort()
				peers = append(peers, peer.Peer{
					Addr: "localhost",
					Port: ports[i],
				})
			}

			for i := 0; i < len; i++ {
				nodes[i], err = bmmc.New(&bmmc.Config{
					Addr:      "localhost",
					Port:      ports[i],
					Peers:     peers,
					Callbacks: callback.NewRegistry(),
				})
				Expect(err).To(Succeed())
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
				nodes[randomNode].AddMessage(msg, callback.NOCALLBACK)
				Eventually(getSortedBuffer(nodes[randomNode]), time.Second).Should(Equal(expectedBuf))

				// start protocol for all nodes
				for i := 0; i < len; i++ {
					Expect(nodes[i].Start())
				}
			})

			It("sync all nodes with the message", func() {
				for i := range nodes {
					Eventually(getSortedBuffer(nodes[i]), time.Second).Should(Equal(expectedBuf))
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
					nodes[randomNode].AddMessage(msg, callback.NOCALLBACK)
				}

				Eventually(getSortedBuffer(nodes[randomNode]), time.Second).Should(Equal(expectedBuf))

				// start protocol for all nodes
				for i := 0; i < len; i++ {
					Expect(nodes[i].Start())
				}
			})

			It("sync all nodes with all messages", func() {
				for i := range nodes {
					Eventually(getSortedBuffer(nodes[i]), time.Second*2).Should(Equal(expectedBuf))
				}
			})
		})

		When("three nodes have three different messages", func() {
			var expectedBuf []string

			BeforeEach(func() {
				randomNodes := [3]int{2, 4, 6}

				for i := 0; i < 3; i++ {
					msg := fmt.Sprintf("awesome-message-%d", i)
					expectedBuf = append(expectedBuf, msg)
					nodes[randomNodes[i]].AddMessage(msg, callback.NOCALLBACK)

					Eventually(getSortedBuffer(nodes[randomNodes[i]]), time.Second).Should(Equal([]string{msg}))
				}

				// start protocol for all nodes
				for i := 0; i < len; i++ {
					Expect(nodes[i].Start())
				}
			})

			It("sync all nodes with all messages", func() {
				for i := range nodes {
					Eventually(getSortedBuffer(nodes[i]), time.Second*3).Should(Equal(expectedBuf))
				}
			})
		})
	})
})
