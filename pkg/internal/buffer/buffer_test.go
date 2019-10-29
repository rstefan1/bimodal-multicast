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
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Buffer interface", func() {
	Describe("NewBuffer function", func() {
		It("creates new buffer", func() {
			buf := NewBuffer(5)

			Expect(buf.Elements).To(HaveLen(5))
			Expect(buf.Len).To(BeZero())
			Expect(buf.Mux).NotTo(BeNil())
		})
	})

	Describe("elementPosition function", func() {
		fullBuf := &Buffer{
			Elements: make([]Element, 4),
			Len:      4,
			Mux:      &sync.Mutex{},
		}
		fullBuf.Elements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
		fullBuf.Elements[1] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
		fullBuf.Elements[2] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}
		fullBuf.Elements[3] = Element{Timestamp: time.Date(2012, time.October, 29, 0, 0, 0, 0, time.UTC)}

		halfBuf := &Buffer{
			Elements: make([]Element, 4),
			Len:      2,
			Mux:      &sync.Mutex{},
		}
		halfBuf.Elements[0] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
		halfBuf.Elements[1] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}

		It("returns position when element should be first", func() {
			el := Element{
				Timestamp: time.Date(2019, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			Expect(fullBuf.elementPosition(el)).To(Equal(0))
		})

		It("returns position when element should be in the middle", func() {
			el := Element{
				Timestamp: time.Date(2015, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			Expect(fullBuf.elementPosition(el)).To(Equal(2))
		})

		It("returns position when element should be last and buffer is full", func() {
			el := Element{
				Timestamp: time.Date(2013, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			Expect(fullBuf.elementPosition(el)).To(Equal(3))
		})

		It("returns position when element should be last and buffer is not full", func() {
			el := Element{
				Timestamp: time.Date(2010, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			Expect(halfBuf.elementPosition(el)).To(Equal(2))
		})

		It("returns -1 when element already exists and is first", func() {
			el := Element{
				Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			Expect(fullBuf.elementPosition(el)).To(Equal(-1))
		})

		It("returns -1 when element already exists and is in the middle", func() {
			el := Element{
				Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			Expect(fullBuf.elementPosition(el)).To(Equal(-1))
		})

		It("returns -1 when element already exists and is last", func() {
			el := Element{
				Timestamp: time.Date(2012, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			Expect(fullBuf.elementPosition(el)).To(Equal(-1))
		})

		It("returns -1 when element is too old and buffer is full", func() {
			el := Element{
				Timestamp: time.Date(2010, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			Expect(fullBuf.elementPosition(el)).To(Equal(-1))
		})
	})
})
