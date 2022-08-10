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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GossipRound", func() {
	DescribeTable("Increment function", func(initial, expected int64) {
		r := NewGossipRound()
		r.Number = initial

		r.Increment()
		Expect(r.Number).To(Equal(expected))
	},
		Entry("counter is smaller then the max round number", int64(5), int64(6)),
		Entry("counter is equal to max round number minus 1", maxRoundNumber-1, maxRoundNumber),
		Entry("counter is equal to max round number", maxRoundNumber, int64(1)),
		Entry("counter is equal to max round number plus 1", maxRoundNumber+1, int64(1)),
	)
})
