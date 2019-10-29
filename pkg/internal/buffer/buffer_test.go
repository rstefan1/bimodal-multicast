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

	Describe("shiftElements function", func() {
		When("buffer is full", func() {
			var buf *Buffer

			BeforeEach(func() {
				buf = &Buffer{
					Elements: make([]Element, 4),
					Len:      4,
					Mux:      &sync.Mutex{},
				}
				buf.Elements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
				buf.Elements[1] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
				buf.Elements[2] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}
				buf.Elements[3] = Element{Timestamp: time.Date(2012, time.October, 29, 0, 0, 0, 0, time.UTC)}
			})

			It("shift all elements if index is 0", func() {
				expectedElements := make([]Element, 4)
				expectedElements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[1] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[2] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[3] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}

				buf.shiftElements(0)

				Expect(buf.Elements).To(Equal(expectedElements))
			})

			It("shift elements if index is in the middle of array", func() {
				expectedElements := make([]Element, 4)
				expectedElements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[1] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[2] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[3] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}

				buf.shiftElements(1)

				Expect(buf.Elements).To(Equal(expectedElements))
			})

			It("doesn't shift any element if index is index of last element", func() {
				expectedElements := make([]Element, 4)
				expectedElements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[1] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[2] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[3] = Element{Timestamp: time.Date(2012, time.October, 29, 0, 0, 0, 0, time.UTC)}

				buf.shiftElements(3)

				Expect(buf.Elements).To(Equal(expectedElements))
			})

			It("doesn't shift any element if index is lower than 0", func() {
				// TODO: not implemented
			})

			It("doesn't shift any elements if index is greater then buffer size", func() {
				// TODO: not implemented
			})
		})

		When("buffer is not full", func() {
			var buf *Buffer

			BeforeEach(func() {
				buf = &Buffer{
					Elements: make([]Element, 4),
					Len:      3,
					Mux:      &sync.Mutex{},
				}
				buf.Elements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
				buf.Elements[1] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
				buf.Elements[2] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}
			})

			It("shift all elements if index is 0", func() {
				expectedElements := make([]Element, 4)
				expectedElements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[1] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[2] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[3] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}

				buf.shiftElements(0)

				Expect(buf.Elements).To(Equal(expectedElements))
			})

			It("shift elements if index is in the middle of array", func() {
				expectedElements := make([]Element, 4)
				expectedElements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[1] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[2] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[3] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}

				buf.shiftElements(1)

				Expect(buf.Elements).To(Equal(expectedElements))
			})

			It("doesn't shift any element if index is the index of last element", func() {
				expectedElements := make([]Element, 4)
				expectedElements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[1] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[2] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[3] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}

				buf.shiftElements(2)

				Expect(buf.Elements).To(Equal(expectedElements))
			})

			It("doesn't shift any element if index is lower than 0", func() {
				// TODO: not implemented
			})

			It("doesn't shift any elements if index is greater then buffer size", func() {
				// TODO: not implemented
			})

			It("doesn't shift any element if index is greater then buffer len but lower then buffer size", func() {
				buf = &Buffer{
					Elements: make([]Element, 4),
					Len:      1,
					Mux:      &sync.Mutex{},
				}
				buf.Elements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}

				expectedElements := make([]Element, 4)
				expectedElements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}

				buf.shiftElements(2)

				Expect(buf.Elements).To(Equal(expectedElements))
			})
		})
	})

	Describe("Add function", func() {
		It("adds the new element in buffer", func() {
			buf := &Buffer{
				Elements: make([]Element, 4),
				Len:      4,
				Mux:      &sync.Mutex{},
			}
			buf.Elements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
			buf.Elements[1] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
			buf.Elements[2] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}
			buf.Elements[3] = Element{Timestamp: time.Date(2012, time.October, 29, 0, 0, 0, 0, time.UTC)}

			expectedElements := make([]Element, 4)
			expectedElements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
			expectedElements[1] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
			expectedElements[2] = Element{Timestamp: time.Date(2015, time.October, 29, 0, 0, 0, 0, time.UTC)}
			expectedElements[3] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}

			el := Element{Timestamp: time.Date(2015, time.October, 29, 0, 0, 0, 0, time.UTC)}

			buf.Add(el)

			Expect(buf.Elements).To(Equal(expectedElements))
		})
	})
})
