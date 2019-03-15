package buffer

import (
	"fmt"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func createDigest(ID string, cnt int) ([]string, *DigestBuffer) {
	digestBuffer := &DigestBuffer{}
	var digestSlice []string
	for i := 0; i < cnt; i++ {
		_ID := fmt.Sprintf("%s-%02d", ID, i)
		digestBuffer.Digests = append(digestBuffer.Digests, Digest{ID: _ID})
		digestSlice = append(digestSlice, _ID)
	}
	return digestSlice, digestBuffer
}

func createDigestFromInitialBuffer(ID string, beginIndex, endIndex int, digestBuffer *DigestBuffer) *DigestBuffer {
	for i := beginIndex; i < endIndex; i++ {
		_ID := fmt.Sprintf("%s-%02d", ID, i)
		digestBuffer.Digests = append(digestBuffer.Digests, Digest{ID: _ID})
	}
	return digestBuffer
}

var _ = Describe("DigestBuffer interface", func() {
	var (
		digestBuffer *DigestBuffer
		digestSlice  []string
		digestID     string
		digestCount  = 3
	)

	BeforeEach(func() {
		digestID = fmt.Sprintf("%d", rand.Int31())
		digestSlice, digestBuffer = createDigest(digestID, digestCount)
	})

	Describe("at WrapDigestBuffer function call", func() {
		It("wraps digest buffer", func() {
			wrappedDigest := WrapDigestBuffer(digestSlice)
			Expect(wrappedDigest).To(Equal(digestBuffer))
		})
	})

	Describe("at SameDigests function call", func() {
		var (
			digestBuffer1 *DigestBuffer
			extraDigest   Digest
		)

		BeforeEach(func() {
			_, digestBuffer1 = createDigest(digestID, digestCount)
			extraDigest = Digest{ID: "extra"}
		})

		It("returns false when a digest has an extra digest at the beginning of buffer", func() {
			digestBuffer2 := &DigestBuffer{}
			digestBuffer2.Digests = append(digestBuffer2.Digests, extraDigest)
			digestBuffer2 = createDigestFromInitialBuffer(digestID, 0, digestCount, digestBuffer2)

			Expect(digestBuffer1.SameDigests(digestBuffer2)).To(Equal(false))
			Expect(digestBuffer2.SameDigests(digestBuffer1)).To(Equal(false))
		})

		It("returns false when a digest has an extra digest at the end of buffer", func() {
			_, digestBuffer2 := createDigest(digestID, digestCount)
			digestBuffer2.Digests = append(digestBuffer2.Digests, extraDigest)

			Expect(digestBuffer1.SameDigests(digestBuffer2)).To(Equal(false))
			Expect(digestBuffer2.SameDigests(digestBuffer1)).To(Equal(false))
		})

		It("returns false when a digest has an extra digest in the mIDdle of buffer", func() {
			_, digestBuffer2 := createDigest(digestID, digestCount/2)
			digestBuffer2.Digests = append(digestBuffer2.Digests, extraDigest)
			digestBuffer2 = createDigestFromInitialBuffer(digestID, digestCount/2, digestCount, digestBuffer2)

			Expect(digestBuffer1.SameDigests(digestBuffer2)).To(Equal(false))
			Expect(digestBuffer2.SameDigests(digestBuffer1)).To(Equal(false))
		})

		It("returns true when the digests are same", func() {
			_, digestBuffer2 := createDigest(digestID, digestCount)
			Expect(digestBuffer1.SameDigests(digestBuffer2)).To(Equal(true))
		})
	})

	Describe("at GetMissignDigests function call", func() {
		var (
			digestBuffer1     *DigestBuffer
			extraDigest       Digest
			emptyDigestBuffer *DigestBuffer
			extraDigestBuffer *DigestBuffer
		)

		BeforeEach(func() {
			_, digestBuffer1 = createDigest(digestID, digestCount)
			extraDigest = Digest{ID: digestID}
			emptyDigestBuffer = &DigestBuffer{Digests: []Digest{}}
			extraDigestBuffer = &DigestBuffer{Digests: []Digest{extraDigest}}
		})

		It("returns empty digest buffer when digests are same", func() {
			_, digestBuffer2 := createDigest(digestID, digestCount)
			Expect(digestBuffer2.GetMissingDigests(digestBuffer1)).To(Equal(emptyDigestBuffer))
		})

		It("returns extra digest when a buffer has an extra digest at the beginning", func() {
			digestBuffer2 := &DigestBuffer{}
			digestBuffer2.Digests = append(digestBuffer2.Digests, extraDigest)
			digestBuffer2 = createDigestFromInitialBuffer(digestID, 0, digestCount, digestBuffer2)

			Expect(digestBuffer1.GetMissingDigests(digestBuffer2)).To(Equal(emptyDigestBuffer))
			Expect(digestBuffer2.GetMissingDigests(digestBuffer1)).To(Equal(extraDigestBuffer))
		})

		It("returns extra digest when a buffer has an extra digest at the end", func() {
			_, digestBuffer2 := createDigest(digestID, digestCount)
			digestBuffer2.Digests = append(digestBuffer2.Digests, extraDigest)

			Expect(digestBuffer1.GetMissingDigests(digestBuffer2)).To(Equal(emptyDigestBuffer))
			Expect(digestBuffer2.GetMissingDigests(digestBuffer1)).To(Equal(extraDigestBuffer))
		})

		It("returns extra digest when a buffer  has an extra digest in the middle of  buffer", func() {
			_, digestBuffer2 := createDigest(digestID, digestCount/2)
			digestBuffer2.Digests = append(digestBuffer2.Digests, Digest{ID: digestID})
			digestBuffer2 = createDigestFromInitialBuffer(digestID, digestCount/2, digestCount, digestBuffer2)

			Expect(digestBuffer1.GetMissingDigests(digestBuffer2)).To(Equal(emptyDigestBuffer))
			Expect(digestBuffer2.GetMissingDigests(digestBuffer1)).To(Equal(extraDigestBuffer))
		})
	})

	Describe("at ContainsDigest function call", func() {
		var digest = Digest{ID: digestID}

		It("returns false when buffer does not contains the given digest", func() {
			Expect(digestBuffer.ContainsDigest(digest)).To(Equal(false))
		})

		It("returns true when the buffer contains the given digest at the beginning", func() {
			digestBuffer := &DigestBuffer{}
			digestBuffer.Digests = append(digestBuffer.Digests, digest)
			digestBuffer = createDigestFromInitialBuffer(digestID, 0, digestCount, digestBuffer)

			Expect(digestBuffer.ContainsDigest(digest)).To(Equal(true))
		})

		It("returns true when the buffer contains the given digest at the end", func() {
			digestBuffer.Digests = append(digestBuffer.Digests, digest)
			Expect(digestBuffer.ContainsDigest(digest)).To(Equal(true))
		})

		It("returs true when the buffer contains the given digest in the mIDdle", func() {
			_, digestBuffer := createDigest(digestID, digestCount/2)
			digestBuffer.Digests = append(digestBuffer.Digests, digest)
			digestBuffer = createDigestFromInitialBuffer(digestID, digestCount/2, digestCount, digestBuffer)

			Expect(digestBuffer.ContainsDigest(digest)).To(Equal(true))
		})
	})

	Describe("at GetMissingMessageBuffer function call", func() {
		var (
			msgBuffer         *MessageBuffer
			expectedMsgBuffer *MessageBuffer
			msgMsg            string
			msgGossipCount    int
		)

		BeforeEach(func() {
			msgMsg = fmt.Sprintf("%d", rand.Int31())
			msgGossipCount = int(rand.Int31())

			msgBuffer = NewMessageBuffer()
			for i := 0; i < digestCount; i++ {
				_id := fmt.Sprintf("%s-%02d", digestID, i)
				_msgMsg := fmt.Sprintf("%s-%02d", msgMsg, i)
				msgBuffer.Messages = append(msgBuffer.Messages,
					Message{
						ID:          _id,
						Msg:         _msgMsg,
						GossipCount: msgGossipCount,
					},
				)
			}
		})

		It("returns the message buffer in concordance with the given digest buffer", func() {
			expectedMsgBuffer = NewMessageBuffer()
			expectedMsgBuffer.Messages = append(expectedMsgBuffer.Messages, msgBuffer.Messages...)

			// add extra message in initial message buffer
			msgBuffer.Messages = append(msgBuffer.Messages,
				Message{
					ID:          digestID,
					Msg:         msgMsg,
					GossipCount: msgGossipCount,
				},
			)

			newMsgBuffer := digestBuffer.GetMissingMessageBuffer(msgBuffer)
			mux1 := fmt.Sprintf("%p", newMsgBuffer.Mux)
			mux2 := fmt.Sprintf("%p", msgBuffer.Mux)
			Expect(mux1).NotTo(Equal(mux2))
			Expect(newMsgBuffer.Messages).To(Equal(expectedMsgBuffer.Messages))
		})

		It("returns an empty message buffer if the given message buffer does not have digests anymore", func() {
			msgBuffer := NewMessageBuffer()
			newMsgBuffer := digestBuffer.GetMissingMessageBuffer(msgBuffer)
			mux1 := fmt.Sprintf("%p", newMsgBuffer.Mux)
			mux2 := fmt.Sprintf("%p", msgBuffer.Mux)

			Expect(mux1).NotTo(Equal(mux2))
			Expect(newMsgBuffer.Messages).To(HaveLen(0))
		})

		It("returns an empty message buffer if the digest buffer is empty", func() {
			digestBuffer := DigestBuffer{}
			newMsgBuffer := digestBuffer.GetMissingMessageBuffer(msgBuffer)
			Expect(newMsgBuffer.Messages).To(HaveLen(0))
		})
	})
})
