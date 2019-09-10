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
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Default Callback interface", func() {
	DescribeTable("DecomposeAddPeerMessage helper function return proper addr and port",
		func(msg, expectedAddr, expectedPort string) {
			addr, port, err := DecomposeAddPeerMessage(msg)
			Expect(err).To(BeNil())
			Expect(addr).To(Equal(expectedAddr))
			Expect(port).To(Equal(expectedPort))
		},
		Entry("message contains a host name", "localhost/9090", "localhost", "9090"),
		Entry("message contains an ip", "127.120.200.100/7070", "127.120.200.100", "7070"),
		Entry("message contains an empty host", "/6060", "", "6060"),
		Entry("message contains an empty port", "localhost/", "localhost", ""),
		Entry("message contains empty host and empty port", "/", "", ""),
	)

	DescribeTable("DecomposeAddPeerMessage helper function return error", func(msg string) {
		_, _, err := DecomposeAddPeerMessage(msg)
		Expect(err).To(Equal(errors.New(invalidAddPeerMsgErr)))
	},
		Entry("message is invalid", "localhost"),
		Entry("message is invalid", "127.100.120.0"),
		Entry("message is invalid", ""),
	)

	DescribeTable("DecomposeRemovePeerMessage helper function return proper addr and port",
		func(msg, expectedAddr, expectedPort string) {
			addr, port, err := DecomposeRemovePeerMessage(msg)
			Expect(err).To(BeNil())
			Expect(addr).To(Equal(expectedAddr))
			Expect(port).To(Equal(expectedPort))
		},
		Entry("message contains a host name", "localhost/9090", "localhost", "9090"),
		Entry("message contains an ip", "127.120.200.100/7070", "127.120.200.100", "7070"),
		Entry("message contains an empty host", "/6060", "", "6060"),
		Entry("message contains an empty port", "localhost/", "localhost", ""),
		Entry("message contains empty host and empty port", "/", "", ""),
	)

	DescribeTable("DecomposeRemovePeerMessage helper function return error", func(msg string) {
		_, _, err := DecomposeRemovePeerMessage(msg)
		Expect(err).To(Equal(errors.New(invalidRemovePeerMsgErr)))
	},
		Entry("message is invalid", "localhost"),
		Entry("message is invalid", "127.100.120.0"),
		Entry("message is invalid", ""),
	)
})
