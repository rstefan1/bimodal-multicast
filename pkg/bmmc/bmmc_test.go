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
	When("System has two nodes and one node has a message in buffer", func() {
		var (
			node1 *bmmc.Bmmc
			node2 *bmmc.Bmmc
			msg   string
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
			node1.AddMessage(msg)
			Eventually(getSortedBuffer(node1), time.Second).Should(SatisfyAll(
				HaveLen(1),
				Equal([]string{msg}),
			))
		})

		AfterEach(func() {
			node1.Stop()
			node2.Stop()
		})

		It("sync buffers with the message", func() {
			Eventually(getSortedBuffer(node2), time.Second).Should(SatisfyAll(
				HaveLen(1),
				Equal([]string{msg}),
			))
		})
	})
})
