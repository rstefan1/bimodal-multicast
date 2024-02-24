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
	"math"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
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
			Expect(err).ToNot(HaveOccurred())
			Expect(pos).To(Equal(0))
		})

		It("returns position when element should be in the middle", func() {
			el := Element{
				Timestamp: time.Date(2015, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			pos, err := fullBuf.elementPosition(el)
			Expect(err).ToNot(HaveOccurred())
			Expect(pos).To(Equal(2))
		})

		It("returns position when element should be last and buffer is full", func() {
			el := Element{
				Timestamp: time.Date(2013, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			pos, err := fullBuf.elementPosition(el)
			Expect(err).ToNot(HaveOccurred())
			Expect(pos).To(Equal(3))
		})

		It("returns position when element should be last and buffer is not full", func() {
			el := Element{
				Timestamp: time.Date(2010, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			pos, err := halfBuf.elementPosition(el)
			Expect(err).ToNot(HaveOccurred())
			Expect(pos).To(Equal(2))
		})

		It("returns error when element is too old and buffer is full", func() {
			el := Element{
				Timestamp: time.Date(2010, time.October, 29, 0, 0, 0, 0, time.UTC),
			}

			_, err := fullBuf.elementPosition(el)
			Expect(err).To(MatchError(errTooOldElement))
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
				Expect(buf.shiftElements(-1)).To(MatchError(errIndexOutOfRange))
			})

			It("doesn't shift any elements if index is greater then buffer size", func() {
				Expect(buf.shiftElements(len(buf.Elements))).To(MatchError(errIndexOutOfRange))
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
		var buf *Buffer

		BeforeEach(func() {
			buf = &Buffer{
				Elements: make([]Element, 4),
				Len:      4,
				Mux:      &sync.Mutex{},
			}
			buf.Elements[0] = Element{
				Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC),
				ID:        "2018",
			}
			buf.Elements[1] = Element{
				Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC),
				ID:        "2016",
			}
			buf.Elements[2] = Element{
				Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC),
				ID:        "2014",
			}
			buf.Elements[3] = Element{
				Timestamp: time.Date(2012, time.October, 29, 0, 0, 0, 0, time.UTC),
				ID:        "2012",
			}
		})

		It("adds the new element in buffer", func() {
			expectedElements := make([]Element, 4)
			expectedElements[0] = Element{
				Timestamp: time.Date(2018, time.October, 29, 0, 0, 0, 0, time.UTC),
				ID:        "2018",
			}
			expectedElements[1] = Element{
				Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC),
				ID:        "2016",
			}
			expectedElements[2] = Element{
				Timestamp: time.Date(2015, time.October, 29, 0, 0, 0, 0, time.UTC),
				ID:        "2015",
			}
			expectedElements[3] = Element{
				Timestamp: time.Date(2014, time.October, 29, 0, 0, 0, 0, time.UTC),
				ID:        "2014",
			}

			el := Element{
				Timestamp: time.Date(2015, time.October, 29, 0, 0, 0, 0, time.UTC),
				ID:        "2015",
			}

			Expect(buf.Add(el)).To(Succeed())

			Expect(buf.Elements).To(Equal(expectedElements))
		})

		It("returns error when buffer already contains given element", func() {
			el := Element{
				Timestamp: time.Date(2016, time.October, 29, 0, 0, 0, 0, time.UTC),
				ID:        "2016",
			}

			Expect(buf.Add(el)).To(MatchError(errAlreadyExists))
		})
	})

	Describe("Digest function", func() {
		It("returns proper digest when buffer is full", func() {
			fullBuf := &Buffer{
				Elements: make([]Element, 4),
				Len:      4,
				Mux:      &sync.Mutex{},
			}
			fullBuf.Elements[0] = Element{ID: "100"}
			fullBuf.Elements[1] = Element{ID: "110"}
			fullBuf.Elements[2] = Element{ID: "107"}
			fullBuf.Elements[3] = Element{ID: "104"}

			expectedDigest := []string{"100", "110", "107", "104"}

			Expect(fullBuf.Digest()).To(Equal(expectedDigest))
		})

		It("returns proper digest when buffer is not full", func() {
			halfBuf := &Buffer{
				Elements: make([]Element, 4),
				Len:      2,
				Mux:      &sync.Mutex{},
			}
			halfBuf.Elements[0] = Element{ID: "204"}
			halfBuf.Elements[1] = Element{ID: "201"}

			expectedDigest := [2]string{"204", "201"}

			Expect(halfBuf.Digest()).To(ConsistOf(expectedDigest))
		})
	})

	Describe("exists function", func() {
		buf := &Buffer{
			Elements: make([]Element, 4),
			Len:      4,
			Mux:      &sync.Mutex{},
		}
		buf.Elements[0] = Element{ID: "100"}
		buf.Elements[1] = Element{ID: "110"}
		buf.Elements[2] = Element{ID: "107"}
		buf.Elements[3] = Element{ID: "104"}

		It("return false if buffer doesn't contain the element", func() {
			el := Element{ID: "90"}

			e, _ := buf.contains(el)
			Expect(e).To(BeFalse())
		})

		It("return true if buffer contains the element and it is first", func() {
			el := Element{ID: "100"}

			e, pos := buf.contains(el)
			Expect(e).To(BeTrue())
			Expect(pos).To(Equal(0))
		})

		It("return true if buffer contains the element and it is in the middle", func() {
			el := Element{ID: "107"}

			e, pos := buf.contains(el)
			Expect(e).To(BeTrue())
			Expect(pos).To(Equal(2))
		})

		It("return true if buffer contains the element and it is last", func() {
			el := Element{ID: "104"}

			e, pos := buf.contains(el)
			Expect(e).To(BeTrue())
			Expect(pos).To(Equal(3))
		})
	})

	Describe("IncrementGoosipCount function", func() {
		It("increments gossip count for all elements from buffer", func() {
			buf := &Buffer{
				Elements: make([]Element, 4),
				Len:      3,
				Mux:      &sync.Mutex{},
			}
			buf.Elements[0] = Element{GossipCount: int64(100)}
			buf.Elements[1] = Element{GossipCount: int64(200)}
			buf.Elements[2] = Element{GossipCount: int64(300)}

			expectedElements := make([]Element, 4)
			expectedElements[0] = Element{GossipCount: int64(101)}
			expectedElements[1] = Element{GossipCount: int64(201)}
			expectedElements[2] = Element{GossipCount: int64(301)}

			buf.IncrementGossipCount()
			Expect(buf.Elements).To(Equal(expectedElements))
		})

		It("doesn't increment gossip count when it is equal with MAX_INT_64", func() {
			buf := &Buffer{
				Elements: []Element{
					{GossipCount: int64(math.MaxInt64 - 2)},
					{GossipCount: int64(math.MaxInt64 - 1)},
					{GossipCount: int64(math.MaxInt64)},
				},
				Len: 3,
				Mux: &sync.Mutex{},
			}

			expectedElements := []Element{
				{GossipCount: int64(math.MaxInt64 - 1)},
				{GossipCount: int64(math.MaxInt64)},
				{GossipCount: int64(0)},
			}

			buf.IncrementGossipCount()
			Expect(buf.Elements).To(Equal(expectedElements))
		})
	})

	Describe("Messages function", func() {
		type testType struct {
			String  string
			Int     int
			Boolean bool
		}

		var buf *Buffer

		BeforeEach(func() {
			buf = &Buffer{
				Elements: []Element{
					{
						Msg:      "string",
						Internal: false,
					},
					{
						Msg:      100,
						Internal: false,
					},
					{
						Msg:      true,
						Internal: false,
					},
					{
						Msg:      "internal-element",
						Internal: true,
					},
					{
						Msg: testType{
							String:  "another-string",
							Int:     200,
							Boolean: false,
						},
						Internal: false,
					},
				},
				Len: 5,
				Mux: &sync.Mutex{},
			}
		})

		It("returns all messages (internal + not internal) from each element from buffer", func() {
			expectedMsg := []any{
				"string",
				100,
				true,
				"internal-element",
				testType{
					String:  "another-string",
					Int:     200,
					Boolean: false,
				},
			}

			Expect(buf.Messages(true)).To(Equal(expectedMsg))
		})

		It("doesn't return messages from internal elements", func() {
			expectedMsg := []any{
				"string",
				100,
				true,
				testType{
					String:  "another-string",
					Int:     200,
					Boolean: false,
				},
			}

			Expect(buf.Messages(false)).To(Equal(expectedMsg))
		})
	})

	Describe("Length function", func() {
		It("returns number of elements in buffer", func() {
			buf := &Buffer{
				Elements: make([]Element, 4),
				Len:      2,
				Mux:      &sync.Mutex{},
			}

			Expect(buf.Length()).To(Equal(2))
		})
	})

	Describe("ElementsFromIDs function", func() {
		It("return elements from buffer", func() {
			buf := &Buffer{
				Elements: make([]Element, 10),
				Len:      10,
				Mux:      &sync.Mutex{},
			}
			buf.Elements[0] = Element{ID: "100"}
			buf.Elements[1] = Element{ID: "101"}
			buf.Elements[2] = Element{ID: "102"}
			buf.Elements[3] = Element{ID: "103"}
			buf.Elements[4] = Element{ID: "104"}
			buf.Elements[5] = Element{ID: "105"}
			buf.Elements[6] = Element{ID: "106"}
			buf.Elements[7] = Element{ID: "107"}
			buf.Elements[8] = Element{ID: "108"}
			buf.Elements[9] = Element{ID: "109"}

			digest := []string{"100", "109", "105", "200", "106"}

			expectedElements := []Element{
				{ID: "100"},
				{ID: "105"},
				{ID: "106"},
				{ID: "109"},
			}

			Expect(buf.ElementsFromIDs(digest)).To(Equal(expectedElements))
		})
	})
})
