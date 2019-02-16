package buffer

import (
	"fmt"
	"math/rand"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func createDigest(id string, cnt int) ([]string, DigestBuffer) {
	var digestBuffer DigestBuffer
	var digestSlice []string
	for i := 0; i < cnt; i++ {
		_id := fmt.Sprintf("%s-%02d", id, i)
		digestBuffer.listDigests = append(digestBuffer.listDigests, Digest{id: _id})
		digestSlice = append(digestSlice, _id)
	}
	return digestSlice, digestBuffer
}

func createDigestFromInitialBuffer(id string, beginIndex, endIndex int, digestBuffer DigestBuffer) DigestBuffer {
	for i := beginIndex; i < endIndex; i++ {
		_id := fmt.Sprintf("%s-%02d", id, i)
		digestBuffer.listDigests = append(digestBuffer.listDigests, Digest{id: _id})
	}
	return digestBuffer
}

var _ = Describe("DigestBuffer interface", func() {
	var (
		digestBuffer DigestBuffer
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

	Describe("at UnwrapDigestBuffer function call", func() {
		It("unwraps digest buffer", func() {
			unwrappedDigest := digestBuffer.UnwrapDigestBuffer()
			Expect(unwrappedDigest).To(Equal(digestSlice))
		})
	})

	Describe("at SameDigests function call", func() {
		var (
			digestBuffer1 DigestBuffer
			extraDigest   Digest
		)

		BeforeEach(func() {
			_, digestBuffer1 = createDigest(digestID, digestCount)
		})

		It("returns false when a digest has an extra digest at the beginning of buffer", func() {
			var digestBuffer2 DigestBuffer
			digestBuffer2.listDigests = append(digestBuffer2.listDigests, extraDigest)
			digestBuffer2 = createDigestFromInitialBuffer(digestID, 0, digestCount, digestBuffer2)

			Expect(digestBuffer1.SameDigests(digestBuffer2)).To(Equal(false))
			Expect(digestBuffer2.SameDigests(digestBuffer1)).To(Equal(false))
		})

		It("returns false when a digest has an extra digest at the end of buffer", func() {
			_, digestBuffer2 := createDigest(digestID, digestCount)
			digestBuffer2.listDigests = append(digestBuffer2.listDigests, extraDigest)

			Expect(digestBuffer1.SameDigests(digestBuffer2)).To(Equal(false))
			Expect(digestBuffer2.SameDigests(digestBuffer1)).To(Equal(false))
		})

		It("returns false when a digest has an extra digest in the middle of buffer", func() {
			_, digestBuffer2 := createDigest(digestID, digestCount/2)
			digestBuffer2.listDigests = append(digestBuffer2.listDigests, extraDigest)
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
			digestBuffer1     DigestBuffer
			extraDigest       Digest
			emptyDigestBuffer = DigestBuffer{listDigests: []Digest{}}
		)

		BeforeEach(func() {
			_, digestBuffer1 = createDigest(digestID, digestCount)
			extraDigest = Digest{id: digestID}
		})

		It("returns empty digest buffer when digests are same", func() {
			_, digestBuffer2 := createDigest(digestID, digestCount)
			Expect(digestBuffer2.GetMissingDigests(digestBuffer1)).To(Equal(emptyDigestBuffer))
		})

		It("returns extra digest when a buffer has an extra digest at the beginning", func() {
			var digestBuffer2 DigestBuffer
			digestBuffer2.listDigests = append(digestBuffer2.listDigests, extraDigest)
			digestBuffer2 = createDigestFromInitialBuffer(digestID, 0, digestCount, digestBuffer2)

			Expect(digestBuffer1.GetMissingDigests(digestBuffer2)).To(Equal(emptyDigestBuffer))
			Expect(digestBuffer2.GetMissingDigests(digestBuffer1)).To(Equal(
				DigestBuffer{listDigests: []Digest{extraDigest}}))
		})

		It("returns extra digest when a buffer has an extra digest at the end", func() {
			_, digestBuffer2 := createDigest(digestID, digestCount)
			digestBuffer2.listDigests = append(digestBuffer2.listDigests, extraDigest)

			Expect(digestBuffer1.GetMissingDigests(digestBuffer2)).To(Equal(emptyDigestBuffer))
			Expect(digestBuffer2.GetMissingDigests(digestBuffer1)).To(Equal(
				DigestBuffer{listDigests: []Digest{extraDigest}}))
		})

		It("returns extra digest when a buffer  has an extra digest in the middle of  buffer", func() {
			_, digestBuffer2 := createDigest(digestID, digestCount/2)
			digestBuffer2.listDigests = append(digestBuffer2.listDigests, Digest{id: digestID})
			digestBuffer2 = createDigestFromInitialBuffer(digestID, digestCount/2, digestCount, digestBuffer2)

			Expect(digestBuffer1.GetMissingDigests(digestBuffer2)).To(Equal(emptyDigestBuffer))
			Expect(digestBuffer2.GetMissingDigests(digestBuffer1)).To(Equal(
				DigestBuffer{listDigests: []Digest{extraDigest}}))
		})
	})

	Describe("at ContainsDigest function call", func() {
		var digest = Digest{id: digestID}

		It("returns false when buffer does not contains the given digest", func() {
			Expect(digestBuffer.ContainsDigest(digest)).To(Equal(false))
		})

		It("returns true when the buffer contains the given digest at the beginning", func() {
			digestBuffer := DigestBuffer{}
			digestBuffer.listDigests = append(digestBuffer.listDigests, digest)
			digestBuffer = createDigestFromInitialBuffer(digestID, 0, digestCount, digestBuffer)

			Expect(digestBuffer.ContainsDigest(digest)).To(Equal(true))
		})

		It("returns true when the buffer contains the given digest at the end", func() {
			digestBuffer.listDigests = append(digestBuffer.listDigests, digest)
			Expect(digestBuffer.ContainsDigest(digest)).To(Equal(true))
		})

		It("returs true when the buffer contains the given digest in the middle", func() {
			_, digestBuffer := createDigest(digestID, digestCount/2)
			digestBuffer.listDigests = append(digestBuffer.listDigests, digest)
			digestBuffer = createDigestFromInitialBuffer(digestID, digestCount/2, digestCount, digestBuffer)

			Expect(digestBuffer.ContainsDigest(digest)).To(Equal(true))
		})
	})

	Describe("at GetMissingMessageBuffer function call", func() {
		var (
			msgBuffer         MessageBuffer
			expectedMsgBuffer MessageBuffer
			msgMsg            string
			msgGossipCount    int
		)

		BeforeEach(func() {
			msgMsg = fmt.Sprintf("%d", rand.Int31())
			msgGossipCount = int(rand.Int31())

			msgBuffer = MessageBuffer{mux: &sync.Mutex{}}
			for i := 0; i < digestCount; i++ {
				_id := fmt.Sprintf("%s-%02d", digestID, i)
				_msgMsg := fmt.Sprintf("%s-%02d", msgMsg, i)
				msgBuffer.listMessages = append(msgBuffer.listMessages,
					Message{
						id:          _id,
						msg:         _msgMsg,
						GossipCount: msgGossipCount,
					},
				)
			}
		})

		It("returns the message buffer in concordance with the given digest buffer", func() {
			expectedMsgBuffer = msgBuffer
			msgBuffer.listMessages = append(msgBuffer.listMessages,
				Message{
					id:          digestID,
					msg:         msgMsg,
					GossipCount: msgGossipCount,
				},
			)

			newMsgBuffer := digestBuffer.GetMissingMessageBuffer(msgBuffer)
			mux1 := fmt.Sprintf("%p", newMsgBuffer.mux)
			mux2 := fmt.Sprintf("%p", msgBuffer.mux)
			Expect(mux1).NotTo(Equal(mux2))
			Expect(newMsgBuffer.listMessages).To(Equal(expectedMsgBuffer.listMessages))
		})

		It("returns an empty message buffer if the given message buffer does not have digests anymore", func() {
			msgBuffer := MessageBuffer{mux: &sync.Mutex{}}
			newMsgBuffer := digestBuffer.GetMissingMessageBuffer(msgBuffer)
			mux1 := fmt.Sprintf("%p", newMsgBuffer.mux)
			mux2 := fmt.Sprintf("%p", msgBuffer.mux)

			Expect(mux1).NotTo(Equal(mux2))
			Expect(newMsgBuffer.listMessages).To(HaveLen(0))
		})

		It("returns an empty message buffer if the digest buffer is empty", func() {
			digestBuffer := DigestBuffer{}
			newMsgBuffer := digestBuffer.GetMissingMessageBuffer(msgBuffer)
			Expect(newMsgBuffer.listMessages).To(HaveLen(0))
		})
	})
})
