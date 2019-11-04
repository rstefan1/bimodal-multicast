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
	Describe("MissingStrings function", func() {
		It("returns empty slice when slices are same", func() {
			a := []string{"a", "b", "c", "d", "e"}
			b := []string{"a", "d", "c", "b", "e"}

			Expect(MissingStrings(a, b)).To(HaveLen(0))
		})

		It("returns empty slice when slices contain same element, but in another order", func() {
			a := []string{"a", "b", "c", "d", "e"}
			b := []string{"b", "c", "a", "e", "d"}

			Expect(MissingStrings(a, b)).To(HaveLen(0))
		})

		It("returns extra string when first slice has an extra string at the beginning", func() {
			a := []string{"extra", "a", "b", "c", "d", "e"}
			b := []string{"b", "c", "a", "e", "d"}

			Expect(MissingStrings(a, b)).To(Equal([]string{"extra"}))
		})

		It("returns extra string when first slice has an extra string at the end", func() {
			a := []string{"a", "b", "c", "d", "e", "extra"}
			b := []string{"b", "c", "a", "e", "d"}

			Expect(MissingStrings(a, b)).To(Equal([]string{"extra"}))
		})

		It("returns extra string when first slice has an extra string in the middle", func() {
			a := []string{"a", "b", "c", "extra", "d", "e"}
			b := []string{"b", "c", "a", "e", "d"}

			Expect(MissingStrings(a, b)).To(Equal([]string{"extra"}))
		})

		It("returns extra strings when first slice has more extra strings", func() {
			a := []string{"extra-1", "a", "b", "c", "extra-2", "d", "e", "extra-3"}
			b := []string{"b", "c", "a", "e", "d"}

			Expect(MissingStrings(a, b)).To(Equal([]string{"extra-1", "extra-2", "extra-3"}))
		})
	})
})
