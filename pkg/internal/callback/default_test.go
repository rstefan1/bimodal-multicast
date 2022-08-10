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

package callback

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Default Callback interface", func() {
	Describe("ComposeAddPeerMessage helper function", func() {
		It("returns proper `add peer` message", func() {
			Expect(ComposeAddPeerMessage("127.120.100.0", "1999")).To(Equal("add/127.120.100.0/1999"))
		})
	})

	DescribeTable("DecomposeAddPeerMessage helper function return proper addr and port",
		func(msg, expectedAddr, expectedPort string) {
			addr, port, err := DecomposeAddPeerMessage(msg)
			Expect(err).To(BeNil())
			Expect(addr).To(Equal(expectedAddr))
			Expect(port).To(Equal(expectedPort))
		},
		Entry("message contains a host name", "add/localhost/9090", "localhost", "9090"),
		Entry("message contains an ip", "add/127.120.200.100/7070", "127.120.200.100", "7070"),
		Entry("message contains an empty host", "add//6060", "", "6060"),
		Entry("message contains an empty port", "add/localhost/", "localhost", ""),
		Entry("message contains empty host and empty port", "add//", "", ""),
	)

	Describe("ComposeRemovePeerMessage helepr function", func() {
		It("returns proper `remove peer` message", func() {
			Expect(ComposeRemovePeerMessage("localhost", "9080")).To(Equal("remove/localhost/9080"))
		})
	})

	DescribeTable("DecomposeAddPeerMessage helper function returns error", func(msg string) {
		_, _, err := DecomposeAddPeerMessage(msg)
		Expect(err).To(MatchError(errInvalidAddPeerMsg))
	},
		Entry("message is invalid", "add/localhost"),
		Entry("message is invalid", "add/127.100.120.0"),
		Entry("message is empty", ""),
		Entry("message doesn't contain `add` prefix", "localhost/19999"),
	)

	DescribeTable("DecomposeRemovePeerMessage helper function return proper addr and port",
		func(msg, expectedAddr, expectedPort string) {
			addr, port, err := DecomposeRemovePeerMessage(msg)
			Expect(err).To(BeNil())
			Expect(addr).To(Equal(expectedAddr))
			Expect(port).To(Equal(expectedPort))
		},
		Entry("message contains a host name", "remove/localhost/9090", "localhost", "9090"),
		Entry("message contains an ip", "remove/127.120.200.100/7070", "127.120.200.100", "7070"),
		Entry("message contains an empty host", "remove//6060", "", "6060"),
		Entry("message contains an empty port", "remove/localhost/", "localhost", ""),
		Entry("message contains empty host and empty port", "remove//", "", ""),
	)

	DescribeTable("DecomposeRemovePeerMessage helper function return error", func(msg string) {
		_, _, err := DecomposeRemovePeerMessage(msg)
		Expect(err).To(MatchError(errInvalidRemovePeerMsg))
	},
		Entry("message is invalid", "remove/localhost"),
		Entry("message is invalid", "remove/127.100.120.0"),
		Entry("message is empty", ""),
		Entry("message doesn't contain `remove` prefix", "localhost/19999"),
	)
})
