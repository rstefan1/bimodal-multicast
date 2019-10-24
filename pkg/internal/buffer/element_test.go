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

package buffer

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Buffer interface", func() {
	Describe("NewElement function", func() {
		It("creates new element", func() {
			el, err := NewElement("message", "callback type")
			Expect(err).To(BeNil())

			Expect(el.ID).NotTo(BeEmpty())
			Expect(el.Timestamp).NotTo(BeNil())
			Expect(el.Msg).To(Equal("message"))
			Expect(el.CallbackType).To(Equal("callback type"))
			Expect(el.GossipCount).To(Equal(int64(0)))
		})
	})
})
