package buffer

import (
	"fmt"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func expectProperMessage(msg Message, msgID, msgMsg string, msgGossipCount int) {
	Expect(msg.ID).To(Equal(msgID))
	Expect(msg.Msg).To(Equal(msgMsg))
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
		msgBuffer = NewMessageBuffer()

		msgID = fmt.Sprintf("%d", rand.Int31())
		msgMsg = fmt.Sprintf("%d", rand.Int31())
		msgGossipCount = int(rand.Int31())

		for i := 0; i < msgCount; i++ {
			_msgID := fmt.Sprintf("%s-%02d", msgID, i)
			_msgMsg := fmt.Sprintf("%s-%02d", msgMsg, i)
			msgBuffer = msgBuffer.AddMessage(Message{
				ID:          _msgID,
				Msg:         _msgMsg,
				GossipCount: msgGossipCount,
			})
		}
	})

	Describe("at AddMessage function call", func() {
		It("adds a new message into message buffer", func() {
			msg := Message{
				ID:          msgID,
				Msg:         msgMsg,
				GossipCount: msgGossipCount,
			}
			msgBuffer := NewMessageBuffer()

			msgBuffer = msgBuffer.AddMessage(msg)
			Expect(msgBuffer.Messages).To(HaveLen(1))
			expectProperMessage(msgBuffer.Messages[0], msgID, msgMsg, msgGossipCount)
		})
	})

	Describe("at UnwrapMessageBuffer function call", func() {
		It("unwrap message buffer", func() {
			var expectedUnwrap []Message
			for i := 0; i < msgCount; i++ {
				msg := Message{
					ID:          fmt.Sprintf("%s-%02d", msgID, i),
					Msg:         fmt.Sprintf("%s-%02d", msgMsg, i),
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
				_ID := fmt.Sprintf("%s-%02d", msgID, i)
				expectedDigestBuffer.Digests = append(expectedDigestBuffer.Digests, Digest{ID: _ID})
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
				expectProperMessage(msgBuffer.Messages[i], _msgID, _msgMsg, newGossipCount)
			}
		})
	})
})
