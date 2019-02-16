package buffer

import (
	"fmt"
	"math/rand"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func expectProperMessage(msg Message, msgID, msgMsg string, msgGossipCount int) {
	Expect(msg.id).To(Equal(msgID))
	Expect(msg.msg).To(Equal(msgMsg))
	Expect(msg.GossipCount).To(Equal(msgGossipCount))
}

var _ = Describe("MessageBuffer interface", func() {
	var (
		msgBuffer      MessageBuffer
		msgCount       = 3
		msgID          string
		msgMsg         string
		msgGossipCount int
	)

	BeforeEach(func() {
		msgBuffer = MessageBuffer{
			mux: &sync.Mutex{},
		}

		msgID = fmt.Sprintf("%d", rand.Int31())
		msgMsg = fmt.Sprintf("%d", rand.Int31())
		msgGossipCount = int(rand.Int31())

		for i := 0; i < msgCount; i++ {
			_msgID := fmt.Sprintf("%s-%02d", msgID, i)
			_msgMsg := fmt.Sprintf("%s-%02d", msgMsg, i)
			msgBuffer = msgBuffer.AddMessage(Message{
				id:          _msgID,
				msg:         _msgMsg,
				GossipCount: msgGossipCount,
			})
		}
	})

	Describe("at AddMessage function call", func() {
		It("adds a new message into message buffer", func() {
			msg := Message{
				id:          msgID,
				msg:         msgMsg,
				GossipCount: msgGossipCount,
			}
			msgBuffer := MessageBuffer{
				mux: &sync.Mutex{},
			}

			msgBuffer = msgBuffer.AddMessage(msg)
			Expect(msgBuffer.listMessages).To(HaveLen(1))
			expectProperMessage(msgBuffer.listMessages[0], msgID, msgMsg, msgGossipCount)
		})
	})

	Describe("at UnwrapMessageBuffer function call", func() {
		It("unwrap message buffer", func() {
			var expectedUnwrap []Message
			for i := 0; i < msgCount; i++ {
				msg := Message{
					id:          fmt.Sprintf("%s-%02d", msgID, i),
					msg:         fmt.Sprintf("%s-%02d", msgMsg, i),
					GossipCount: msgGossipCount,
				}
				expectedUnwrap = append(expectedUnwrap, msg)
			}
			Expect(msgBuffer.UnwrapMessageBuffer()).To(Equal(expectedUnwrap))
		})
	})

	Describe("at DigestBuffer fuction call", func() {
		It("transform message buffer into digest buffer", func() {
			var expectedDigestBuffer DigestBuffer
			for i := 0; i < msgCount; i++ {
				_id := fmt.Sprintf("%s-%02d", msgID, i)
				expectedDigestBuffer.listDigests = append(expectedDigestBuffer.listDigests, Digest{id: _id})
			}
			digestBuffer := msgBuffer.DigestBuffer()
			Expect(digestBuffer).To(Equal(expectedDigestBuffer))
		})
	})

	Describe("at IncrementGossipCount function call", func() {
		It("increment gossip count for each message from message buffer", func() {
			msgBuffer = msgBuffer.IncrementGossipCount()
			newGossipCount := msgGossipCount + 1
			for i := 0; i < msgCount; i++ {
				_msgID := fmt.Sprintf("%s-%02d", msgID, i)
				_msgMsg := fmt.Sprintf("%s-%02d", msgMsg, i)
				expectProperMessage(msgBuffer.listMessages[i], _msgID, _msgMsg, newGossipCount)
			}
		})
	})
})
