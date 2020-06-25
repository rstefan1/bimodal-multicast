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

package bmmc

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	DescribeTable("fullHost helper function", func(addr, port, expectedFullHost string) {
		Expect(fullHost(addr, port)).To(Equal(expectedFullHost))
	},
		Entry("returns proper address and port", "127.168.0.100", "8080", "127.168.0.100:8080"),
		Entry("returns proper localhost address and port", "localhost", "7070", "localhost:7070"),
	)

	DescribeTable("addrPort helper function", func(host, expectedAddr, expectedPort string) {
		addr, port, err := addrPort(host)
		Expect(err).To(BeNil())
		Expect(addr).To(Equal(expectedAddr))
		Expect(port).To(Equal(expectedPort))
	},
		Entry("returns proper address and port", "127.168.0.100:8080", "127.168.0.100", "8080"),
		Entry("returns proper localhost address and port", "localhost:7070", "localhost", "7070"),
	)

	DescribeTable("addrPort helper function", func(host string, expectedErr error) {
		_, _, err := addrPort(host)
		Expect(err).To(MatchError(expectedErr))
	},
		Entry("returns error when full host contains only addr or only port", "127.168.0.100", errInvalidHost),
		Entry("returns error when full host contains to much elements", "localhost:127.168.0.100:7070", errInvalidHost),
	)
})
