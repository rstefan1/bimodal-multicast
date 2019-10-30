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
	"fmt"
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

			pos, err := fullBuf.elementPosition(el)
			Expect(err).To(BeNil())
			Expect(pos).To(Equal(0))
		})

		It("returns position when element should be in the middle", func() {
			el := Element{
				Timestamp: time.Date(2015, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			pos, err := fullBuf.elementPosition(el)
			Expect(err).To(BeNil())
			Expect(pos).To(Equal(2))
		})

		It("returns position when element should be last and buffer is full", func() {
			el := Element{
				Timestamp: time.Date(2013, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			pos, err := fullBuf.elementPosition(el)
			Expect(err).To(BeNil())
			Expect(pos).To(Equal(3))
		})

		It("returns position when element should be last and buffer is not full", func() {
			el := Element{
				Timestamp: time.Date(2010, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			pos, err := halfBuf.elementPosition(el)
			Expect(err).To(BeNil())
			Expect(pos).To(Equal(2))
		})

		It("returns error when element already exists and is first", func() {
			el := Element{
				Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			_, err := fullBuf.elementPosition(el)
			Expect(err).To(Equal(fmt.Errorf(alreadyExistsErrFmt)))
		})

		It("returns error when element already exists and is in the middle", func() {
			el := Element{
				Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			_, err := fullBuf.elementPosition(el)
			Expect(err).To(Equal(fmt.Errorf(alreadyExistsErrFmt)))
		})

		It("returns error when element already exists and is last", func() {
			el := Element{
				Timestamp: time.Date(2012, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			_, err := fullBuf.elementPosition(el)
			Expect(err).To(Equal(fmt.Errorf(alreadyExistsErrFmt)))
		})

		It("returns error when element is too old and buffer is full", func() {
			el := Element{
				Timestamp: time.Date(2010, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			_, err := fullBuf.elementPosition(el)
			Expect(err).To(Equal(fmt.Errorf(tooOldElementErrFmt)))
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

				Expect(buf.shiftElements(0)).To(Succeed())

				Expect(buf.Elements).To(Equal(expectedElements))
			})

			It("shift elements if index is in the middle of array", func() {
				expectedElements := make([]Element, 4)
				expectedElements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[1] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[2] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[3] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}

				Expect(buf.shiftElements(1)).To(Succeed())

				Expect(buf.Elements).To(Equal(expectedElements))
			})

			It("doesn't shift any element if index is index of last element", func() {
				expectedElements := make([]Element, 4)
				expectedElements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[1] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[2] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[3] = Element{Timestamp: time.Date(2012, time.October, 29, 0, 0, 0, 0, time.UTC)}

				Expect(buf.shiftElements(3)).To(Succeed())

				Expect(buf.Elements).To(Equal(expectedElements))
			})

			It("doesn't shift any element if index is lower than 0", func() {
				Expect(buf.shiftElements(-1)).To(Equal(fmt.Errorf(indexOutOfRangeErrFmt)))
			})

			It("doesn't shift any elements if index is greater then buffer size", func() {
				Expect(buf.shiftElements(len(buf.Elements))).To(Equal(fmt.Errorf(indexOutOfRangeErrFmt)))
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

				Expect(buf.shiftElements(0)).To(Succeed())

				Expect(buf.Elements).To(Equal(expectedElements))
			})

			It("shift elements if index is in the middle of array", func() {
				expectedElements := make([]Element, 4)
				expectedElements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[1] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[2] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[3] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}

				Expect(buf.shiftElements(1)).To(Succeed())

				Expect(buf.Elements).To(Equal(expectedElements))
			})

			It("doesn't shift any element if index is the index of last element", func() {
				expectedElements := make([]Element, 4)
				expectedElements[0] = Element{Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[1] = Element{Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[2] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}
				expectedElements[3] = Element{Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC)}

				Expect(buf.shiftElements(2)).To(Succeed())

				Expect(buf.Elements).To(Equal(expectedElements))
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

				Expect(buf.shiftElements(2)).To(Succeed())

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

			Expect(buf.Add(el)).To(Succeed())

			Expect(buf.Elements).To(Equal(expectedElements))
		})
	})
})