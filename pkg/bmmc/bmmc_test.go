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
	"github.com/rstefan1/bimodal-multicast/pkg/internal/callback"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/testutil"
)

// getSortedBuffer is a helper func that returns a sorted buffer
func getSortedBuffer(node *bmmc.Bmmc) func() []string {
	return func() []string {
		buf := node.GetMessages()
		sort.Strings(buf)
		return buf
	}
}

func fakeRegistry(cbType string, b bool, e error) map[string]func(string, *log.Logger) (bool, error) {
	return map[string]func(string, *log.Logger) (bool, error){
		cbType: func(msg string, logger *log.Logger) (bool, error) {
			return b, e
		},
	}
}

var _ = Describe("BMMC", func() {
	When("creates new protocol instance with broken config", func() {
		It("returns error when address is empty", func() {
			_, err := bmmc.New(&bmmc.Config{
				Port:      "1999",
				Beta:      0.5,
				Logger:    log.New(os.Stdout, "", 0),
				Callbacks: map[string]func(string, *log.Logger) (bool, error){},
			})
			Expect(err).To(Not(Succeed()))
		})

		It("returns error when port is empty", func() {
			_, err := bmmc.New(&bmmc.Config{
				Addr:      "localhost",
				Beta:      0.5,
				Logger:    log.New(os.Stdout, "", 0),
				Callbacks: map[string]func(string, *log.Logger) (bool, error){},
			})
			Expect(err).To(Not(Succeed()))
		})

		It("doens't returns error when callback CustomRegistry is nil", func() {
			cfg := bmmc.Config{
				Addr:   "localhost",
				Port:   "1999",
				Beta:   0.5,
				Logger: log.New(os.Stdout, "", 0),
			}
			_, err := bmmc.New(&cfg)
			Expect(err).To(Succeed())
		})

		It("set default value for beta if it is empty", func() {
			cfg := bmmc.Config{
				Addr:      "localhost",
				Port:      "1999",
				Logger:    log.New(os.Stdout, "", 0),
				Callbacks: map[string]func(string, *log.Logger) (bool, error){},
			}
			_, err := bmmc.New(&cfg)
			Expect(err).To(Succeed())
			Expect(cfg.Beta).To(Equal(0.3))
		})

		It("set default value for logger if it is empty", func() {
			cfg := bmmc.Config{
				Addr:      "localhost",
				Port:      "1999",
				Beta:      0.5,
				Callbacks: map[string]func(string, *log.Logger) (bool, error){},
			}
			_, err := bmmc.New(&cfg)
			Expect(err).To(Succeed())
			Expect(cfg.Logger).To(Not(BeNil()))
		})
	})

	DescribeTable("when system has two nodes and one node has a message in buffer",
		func(cbCustomRegistry map[string]func(string, *log.Logger) (bool, error),
			msg, callbackType string,
			expectedBuf []string) {

			port1 := testutil.SuggestPort()
			port2 := testutil.SuggestPort()

			node1, err := bmmc.New(&bmmc.Config{
				Addr:      "localhost",
				Port:      port1,
				Beta:      0.5,
				Callbacks: cbCustomRegistry,
			})
			Expect(err).To(Succeed())

			node2, err := bmmc.New(&bmmc.Config{
				Addr:      "localhost",
				Port:      port2,
				Beta:      0.5,
				Callbacks: cbCustomRegistry,
			})
			Expect(err).To(Succeed())

			Expect(node1.Start())
			Expect(node2.Start())

			Expect(node1.AddPeer("localhost", port2)).To(Succeed())
			Expect(node2.AddPeer("localhost", port1)).To(Succeed())

			extraMsgBuffer := []string{
				fmt.Sprintf("localhost/%s", port1),
				fmt.Sprintf("localhost/%s", port2),
			}

			// add a message in first node
			Expect(node1.AddMessage(msg, callbackType)).To(Succeed())
			Eventually(getSortedBuffer(node1), time.Second).Should(
				ConsistOf(append(extraMsgBuffer, msg)))

			Eventually(getSortedBuffer(node2), time.Second).Should(
				ConsistOf(append(extraMsgBuffer, expectedBuf...)))
		},
		Entry("sync buffers with the message",
			map[string]func(string, *log.Logger) (bool, error){},
			"awesome-message",
			callback.NOCALLBACK,
			[]string{"awesome-message"}),
		Entry("sync buffers if callback returned error",
			fakeRegistry("my-callback", true, fmt.Errorf("invalid-callback")),
			"awesome-message",
			"my-callback",
			[]string{"awesome-message"}),
		Entry("doesn't sync buffers if callback returned false",
			fakeRegistry("my-callback", false, nil),
			"awesome-message",
			"my-callback",
			[]string{}),
	)

	When("system has ten nodes", func() {
		const len = 10
		var (
			nodes          [len]*bmmc.Bmmc
			ports          [len]string
			addrs          [len]string
			extraMsgBuffer []string
		)

		BeforeEach(func() {
			var err error
			extraMsgBuffer = make([]string, len)

			for i := 0; i < len; i++ {
				ports[i] = testutil.SuggestPort()
				addrs[i] = "localhost"
				extraMsgBuffer[i] = fmt.Sprintf("%s/%s", addrs[i], ports[i])
			}

			for i := 0; i < len; i++ {
				nodes[i], err = bmmc.New(&bmmc.Config{
					Addr:      "localhost",
					Port:      ports[i],
					Callbacks: map[string]func(string, *log.Logger) (bool, error){},
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
				expectedBuf = []string{}

				msg := "another-awesome-message"
				expectedBuf = append(expectedBuf, msg)

				randomNode := rand.Intn(len)
				Expect(nodes[randomNode].AddMessage(msg, callback.NOCALLBACK)).To(Succeed())
				Eventually(getSortedBuffer(nodes[randomNode]), time.Second).Should(ConsistOf(expectedBuf))

				// start protocol for all nodes
				for p := 0; p < len; p++ {
					for i := 0; i < len; i++ {
						Expect(nodes[p].AddPeer(addrs[i], ports[i])).To(Succeed())
					}
					Expect(nodes[p].Start())
				}
			})

			It("sync all nodes with the message", func() {
				for i := range nodes {
					Eventually(getSortedBuffer(nodes[i]), time.Second*3).Should(ConsistOf(append(expectedBuf, extraMsgBuffer...)))
				}
			})
		})

		When("one node has more messages", func() {
			var expectedBuf []string

			BeforeEach(func() {
				expectedBuf = []string{}

				randomNode := rand.Intn(len)
				for i := 0; i < 3; i++ {
					msg := fmt.Sprintf("awesome-message-%d", i)
					expectedBuf = append(expectedBuf, msg)
					Expect(nodes[randomNode].AddMessage(msg, callback.NOCALLBACK)).To(Succeed())
				}

				Eventually(getSortedBuffer(nodes[randomNode]), time.Second).Should(ConsistOf(expectedBuf))

				// start protocol for all nodes
				for p := 0; p < len; p++ {
					for i := 0; i < len; i++ {
						Expect(nodes[p].AddPeer(addrs[i], ports[i])).To(Succeed())
					}
					Expect(nodes[p].Start())
				}
			})

			It("sync all nodes with all messages", func() {
				for i := range nodes {
					Eventually(getSortedBuffer(nodes[i]), time.Second*5).Should(ConsistOf(append(expectedBuf, extraMsgBuffer...)))
				}
			})
		})

		When("three nodes have three different messages", func() {
			var expectedBuf []string

			BeforeEach(func() {
				expectedBuf = []string{}
				randomNodes := [3]int{2, 4, 6}

				for i := 0; i < 3; i++ {
					msg := fmt.Sprintf("awesome-message-%d", i)
					expectedBuf = append(expectedBuf, msg)
					Expect(nodes[randomNodes[i]].AddMessage(msg, callback.NOCALLBACK)).To(Succeed())

					Eventually(getSortedBuffer(nodes[randomNodes[i]]), time.Second).Should(ConsistOf([]string{msg}))
				}

				// start protocol for all nodes
				for p := 0; p < len; p++ {
					for i := 0; i < len; i++ {
						Expect(nodes[p].AddPeer(addrs[i], ports[i])).To(Succeed())
					}
					Expect(nodes[p].Start())
				}
			})

			It("sync all nodes with all messages", func() {
				for i := range nodes {
					Eventually(getSortedBuffer(nodes[i]), time.Second*5).Should(ConsistOf(append(expectedBuf, extraMsgBuffer...)))
				}
			})
		})
	})
})
