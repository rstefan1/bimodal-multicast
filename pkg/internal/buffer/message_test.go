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
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func expectProperMessage(msg Message, msgID, msgMsg string, msgGossipCount int) {
	Expect(msg.ID).To(Equal(msgID))
	Expect(msg.Msg).To(Equal(msgMsg))
	Expect(msg.GossipCount).To(Equal(msgGossipCount))
}

func createMessageFromInitialBuffer(ID, msg string, gossipCount int, beginIndex, endIndex int, messageBuffer *MessageBuffer) *MessageBuffer {
	for i := beginIndex; i < endIndex; i++ {
		_ID := fmt.Sprintf("%s-%02d", ID, i)
		messageBuffer.Messages = append(messageBuffer.Messages, Message{
			ID:          _ID,
			Msg:         msg,
			GossipCount: gossipCount,
		})
	}
	return messageBuffer
}

func createMessage(ID, msg string, gossipCount, cnt int) ([]Message, *MessageBuffer) {
	var messageSlice []Message
	messageBuffer := NewMessageBuffer()

	for i := 0; i < cnt; i++ {
		_ID := fmt.Sprintf("%s-%02d", ID, i)
		newMsg := Message{
			ID:          _ID,
			Msg:         msg,
			GossipCount: gossipCount,
		}
		messageBuffer.Messages = append(messageBuffer.Messages, newMsg)
		messageSlice = append(messageSlice, newMsg)
	}
	return messageSlice, messageBuffer
}

var _ = Describe("MessageBuffer interface", func() {
	var (
		msgBuffer      *MessageBuffer
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
			msgBuffer.AddMessage(Message{
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

			msgBuffer.AddMessage(msg)
			Expect(msgBuffer.Messages).To(HaveLen(1))
			expectProperMessage(msgBuffer.Messages[0], msgID, msgMsg, msgGossipCount)
		})
	})

	Describe("at DigestBuffer fuction call", func() {
		It("transform message buffer into digest buffer", func() {
			expectedDigestBuffer := &DigestBuffer{}

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
			msgBuffer.IncrementGossipCount()
			newGossipCount := msgGossipCount + 1

			for i := 0; i < msgCount; i++ {
				_msgID := fmt.Sprintf("%s-%02d", msgID, i)
				_msgMsg := fmt.Sprintf("%s-%02d", msgMsg, i)
				expectProperMessage(msgBuffer.Messages[i], _msgID, _msgMsg, newGossipCount)
			}
		})
	})

	Describe("at SameMessages function call", func() {
		var (
			messageBuffer1     *MessageBuffer
			extraMessage       Message
			messageID          string
			messageMsg         string
			messageGossipCount int
			messageCount       int
		)

		BeforeEach(func() {
			messageID = fmt.Sprintf("%d", rand.Int31())
			messageMsg = fmt.Sprintf("%d", rand.Int31())
			messageCount = rand.Int()

			messageCount = 3

			_, messageBuffer1 = createMessage(messageID, messageMsg, messageGossipCount, messageCount)

			extraMessage = Message{
				ID:          "extra",
				Msg:         "extra",
				GossipCount: rand.Int(),
			}
		})

		It("returns false when the message has an extra message at the beginning of buffer", func() {
			messageBuffer2 := NewMessageBuffer()
			messageBuffer2.Messages = append(messageBuffer2.Messages, extraMessage)
			messageBuffer2 = createMessageFromInitialBuffer(messageID, msgMsg, msgGossipCount, 0, messageCount, messageBuffer2)

			Expect(messageBuffer1.SameMessages(messageBuffer2)).To(Equal(false))
			Expect(messageBuffer2.SameMessages(messageBuffer1)).To(Equal(false))
		})

		It("returns false when the message buffer has an extra message at the end of buffer", func() {
			_, messageBuffer2 := createMessage(messageID, messageMsg, messageGossipCount, messageCount)
			messageBuffer2.Messages = append(messageBuffer2.Messages, extraMessage)

			Expect(messageBuffer1.SameMessages(messageBuffer2)).To(Equal(false))
			Expect(messageBuffer2.SameMessages(messageBuffer1)).To(Equal(false))
		})

		It("returns false when the message buffer has an extra message in the middle of buffer", func() {
			_, messageBuffer2 := createMessage(messageID, messageMsg, messageGossipCount, messageCount/2)
			messageBuffer2.Messages = append(messageBuffer2.Messages, extraMessage)
			messageBuffer2 = createMessageFromInitialBuffer(messageID, messageMsg, messageGossipCount, messageCount/2, messageCount, messageBuffer2)

			Expect(messageBuffer1.SameMessages(messageBuffer2)).To(Equal(false))
			Expect(messageBuffer2.SameMessages(messageBuffer1)).To(Equal(false))
		})

		It("returns true when the messages are same", func() {
			_, messageBuffer2 := createMessage(messageID, messageMsg, messageGossipCount, messageCount)
			Expect(messageBuffer1.SameMessages(messageBuffer2)).To(Equal(true))
		})
	})
})
