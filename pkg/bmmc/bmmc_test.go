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

// import (
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"os"
// 	"time"

// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/ginkgo/extensions/table"
// 	. "github.com/onsi/gomega"

// 	"github.com/rstefan1/bimodal-multicast/pkg/bmmc"
// 	"github.com/rstefan1/bimodal-multicast/pkg/internal/callback"
// 	"github.com/rstefan1/bimodal-multicast/pkg/internal/testutil"
// )

// func getBuffer(node *bmmc.Bmmc) []string {
// 	buf := node.GetMessages()

// 	sbuf := make([]string, len(buf))
// 	for i, v := range buf {
// 		sbuf[i] = fmt.Sprint(v)
// 	}

// 	return sbuf
// }

// func getBufferFn(node *bmmc.Bmmc) func() []string {
// 	return func() []string {
// 		return getBuffer(node)
// 	}
// }

// func interfaceToString(b []interface{}) []string {
// 	s := make([]string, len(b))
// 	for i, v := range b {
// 		s[i] = fmt.Sprint(v)
// 	}
// 	return s
// }

// func fakeRegistry(cbType string, e error) map[string]func(interface{}, *log.Logger) error {
// 	return map[string]func(interface{}, *log.Logger) error{
// 		cbType: func(_ interface{}, _ *log.Logger) error {
// 			return e
// 		},
// 	}
// }

// var _ = Describe("BMMC", func() {
// 	When("creates new protocol instance with broken config", func() {
// 		It("returns error when address is empty", func() {
// 			_, err := bmmc.New(&bmmc.Config{
// 				Port:      "1999",
// 				Beta:      0.5,
// 				Logger:    log.New(os.Stdout, "", 0),
// 				Callbacks: map[string]func(interface{}, *log.Logger) error{},
// 			})
// 			Expect(err).To(Not(Succeed()))
// 		})

// 		It("returns error when port is empty", func() {
// 			_, err := bmmc.New(&bmmc.Config{
// 				Addr:      "localhost",
// 				Beta:      0.5,
// 				Logger:    log.New(os.Stdout, "", 0),
// 				Callbacks: map[string]func(interface{}, *log.Logger) error{},
// 			})
// 			Expect(err).To(Not(Succeed()))
// 		})

// 		It("doens't returns error when callback CustomRegistry is nil", func() {
// 			cfg := bmmc.Config{
// 				Addr:   "localhost",
// 				Port:   "1999",
// 				Beta:   0.5,
// 				Logger: log.New(os.Stdout, "", 0),
// 			}
// 			_, err := bmmc.New(&cfg)
// 			Expect(err).To(Succeed())
// 		})

// 		It("set default value for beta if it is empty", func() {
// 			cfg := bmmc.Config{
// 				Addr:      "localhost",
// 				Port:      "1999",
// 				Logger:    log.New(os.Stdout, "", 0),
// 				Callbacks: map[string]func(interface{}, *log.Logger) error{},
// 			}
// 			_, err := bmmc.New(&cfg)
// 			Expect(err).To(Succeed())
// 			Expect(cfg.Beta).To(Equal(0.3))
// 		})

// 		It("set default value for gossip round duration if it is empty", func() {
// 			cfg := bmmc.Config{
// 				Addr:      "localhost",
// 				Port:      "1999",
// 				Beta:      0.5,
// 				Logger:    log.New(os.Stdout, "", 0),
// 				Callbacks: map[string]func(interface{}, *log.Logger) error{},
// 			}
// 			_, err := bmmc.New(&cfg)
// 			Expect(err).To(Succeed())
// 			Expect(cfg.RoundDuration).NotTo(Equal(0))
// 		})

// 		It("set default value for logger if it is empty", func() {
// 			cfg := bmmc.Config{
// 				Addr:      "localhost",
// 				Port:      "1999",
// 				Beta:      0.5,
// 				Callbacks: map[string]func(interface{}, *log.Logger) error{},
// 			}
// 			_, err := bmmc.New(&cfg)
// 			Expect(err).To(Succeed())
// 			Expect(cfg.Logger).To(Not(BeNil()))
// 		})
// 	})

// 	DescribeTable("when system has two nodes and one node has a message in buffer",
// 		func(cbCustomRegistry map[string]func(interface{}, *log.Logger) error,
// 			msg, callbackType string,
// 			expectedBuf []string) {

// 			port1 := testutil.SuggestPort()
// 			port2 := testutil.SuggestPort()

// 			node1, err := bmmc.New(&bmmc.Config{
// 				Addr:      "localhost",
// 				Port:      port1,
// 				Beta:      0.5,
// 				Callbacks: cbCustomRegistry,
// 			})
// 			Expect(err).To(Succeed())

// 			node2, err := bmmc.New(&bmmc.Config{
// 				Addr:      "localhost",
// 				Port:      port2,
// 				Beta:      0.5,
// 				Callbacks: cbCustomRegistry,
// 			})
// 			Expect(err).To(Succeed())

// 			Expect(node1.Start())
// 			Expect(node2.Start())

// 			Expect(node1.AddPeer("localhost", port2)).To(Succeed())
// 			Expect(node2.AddPeer("localhost", port1)).To(Succeed())

// 			extraMsgBuffer := []string{
// 				fmt.Sprintf("localhost/%s", port1),
// 				fmt.Sprintf("localhost/%s", port2),
// 			}

// 			// add a message in first node
// 			Expect(node1.AddMessage(msg, callbackType)).To(Succeed())

// 			Eventually(getBufferFn(node1), time.Second).Should(
// 				ConsistOf(append(extraMsgBuffer, expectedBuf...)))

// 			Eventually(getBufferFn(node2), time.Second).Should(
// 				ConsistOf(append(extraMsgBuffer, expectedBuf...)))
// 		},
// 		Entry("sync buffers if callback return error",
// 			fakeRegistry("my-callback", fmt.Errorf("invalid-callback")),
// 			"awesome-message",
// 			"my-callback",
// 			[]string{"awesome-message"}),
// 		Entry("sync buffers if callback not return false",
// 			fakeRegistry("my-callback", nil),
// 			"awesome-message",
// 			"my-callback",
// 			[]string{"awesome-message"}),
// 	)

// 	When("system has ten nodes", func() {
// 		const len = 10
// 		var (
// 			nodes          [len]*bmmc.Bmmc
// 			ports          [len]string
// 			addrs          [len]string
// 			extraMsgBuffer []interface{}
// 			expectedBuf    []interface{}
// 		)

// 		BeforeEach(func() {
// 			var err error
// 			extraMsgBuffer = make([]interface{}, len)
// 			expectedBuf = []interface{}{}

// 			for i := 0; i < len; i++ {
// 				ports[i] = testutil.SuggestPort()
// 				addrs[i] = "localhost"
// 				extraMsgBuffer[i] = fmt.Sprintf("%s/%s", addrs[i], ports[i])
// 			}

// 			for i := 0; i < len; i++ {
// 				nodes[i], err = bmmc.New(&bmmc.Config{
// 					Addr:      addrs[i],
// 					Port:      ports[i],
// 					Callbacks: map[string]func(interface{}, *log.Logger) error{},
// 				})
// 				Expect(err).To(Succeed())
// 			}

// 			// start protocol for all nodes
// 			for p := 0; p < len; p++ {
// 				for i := 0; i < len; i++ {
// 					_ = nodes[p].AddPeer(addrs[i], ports[i])
// 				}
// 				Expect(nodes[p].Start())
// 			}
// 		})

// 		AfterEach(func() {
// 			for i := range nodes {
// 				nodes[i].Stop()
// 			}
// 		})

// 		When("one node has an message", func() {
// 			BeforeEach(func() {
// 				msg := "another-awesome-message"
// 				expectedBuf = append(expectedBuf, msg)

// 				randomNode := rand.Intn(len)
// 				err := nodes[randomNode].AddMessage(msg, callback.NOCALLBACK)
// 				Expect(err).To(BeNil())
// 				Eventually(getBufferFn(nodes[randomNode]), time.Second).Should(ConsistOf(append(expectedBuf, extraMsgBuffer...)))
// 			})

// 			It("sync all nodes with the message", func() {
// 				for i := range nodes {
// 					Eventually(getBufferFn(nodes[i]), time.Second*3).Should(ConsistOf(append(expectedBuf, extraMsgBuffer...)))
// 				}
// 			})
// 		})

// 		When("one node has more messages", func() {
// 			BeforeEach(func() {
// 				randomNode := rand.Intn(len)
// 				for i := 0; i < 3; i++ {
// 					msg := i
// 					expectedBuf = append(expectedBuf, msg)

// 					err := nodes[randomNode].AddMessage(msg, callback.NOCALLBACK)
// 					Expect(err).To(BeNil())
// 				}

// 				Eventually(getBufferFn(nodes[randomNode]), time.Second).Should(
// 					ConsistOf(interfaceToString(append(expectedBuf, extraMsgBuffer...))))
// 			})

// 			It("sync all nodes with all messages", func() {
// 				for i := range nodes {
// 					Eventually(getBufferFn(nodes[i]), time.Second*5).Should(
// 						ConsistOf(interfaceToString(append(expectedBuf, extraMsgBuffer...))))
// 				}
// 			})
// 		})

// 		When("three nodes have three different messages", func() {
// 			BeforeEach(func() {
// 				randomNodes := [3]int{2, 4, 6}

// 				for i := 0; i < 3; i++ {
// 					msg := i
// 					expectedBuf = append(expectedBuf, msg)

// 					err := nodes[randomNodes[i]].AddMessage(msg, callback.NOCALLBACK)
// 					Expect(err).To(BeNil())
// 				}
// 			})

// 			It("sync all nodes with all messages", func() {
// 				for i := range nodes {
// 					Eventually(getBufferFn(nodes[i]), time.Second*5).Should(
// 						ConsistOf(interfaceToString(append(expectedBuf, extraMsgBuffer...))))
// 				}
// 			})
// 		})
// 	})
// })
