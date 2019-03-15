package buffer

import (
	"sort"
	"sync"
)

type Message struct {
	ID          string `json:"message_ID"`
	Msg         string `json:"message_msg"`
	GossipCount int    `json:"message_gossip_count"`
}

type MessageBuffer struct {
	Messages []Message   `json:"message_buffer_list"`
	Mux      *sync.Mutex `json:"message_buffer_Mux"`
}

// NewMessageBuffer creates new MessageBuffer
func NewMessageBuffer() *MessageBuffer {
	return &MessageBuffer{
		Mux: &sync.Mutex{},
	}
}

// Digest returns a slice with ID of messages from given buffer
func (msgBuffer *MessageBuffer) DigestBuffer() *DigestBuffer {
	digestBuffer := &DigestBuffer{}

	msgBuffer.Mux.Lock()
	defer msgBuffer.Mux.Unlock()

	for _, b := range msgBuffer.Messages {
		digestBuffer.Digests = append(digestBuffer.Digests, Digest{ID: b.ID})
	}

	return digestBuffer
}

// AddMessage adds message in message buffer
func (msgBuffer *MessageBuffer) AddMessage(msg Message) {
	msgBuffer.Mux.Lock()
	defer msgBuffer.Mux.Unlock()

	msgBuffer.Messages = append(msgBuffer.Messages, msg)
}

// IncrementGossipCount increments gossip count for each message from message
// buffer
func (msgBuffer *MessageBuffer) IncrementGossipCount() {
	msgBuffer.Mux.Lock()
	defer msgBuffer.Mux.Unlock()

	for i := range msgBuffer.Messages {
		msgBuffer.Messages[i].GossipCount++
	}
}

func compareMessagesFn(msg []Message) func(int, int) bool {
	return func(i, j int) bool {
		return msg[i].ID <= msg[j].ID
	}
}

func (a *MessageBuffer) SameMessages(b *MessageBuffer) bool {
	a.Mux.Lock()
	b.Mux.Lock()
	defer a.Mux.Unlock()
	defer b.Mux.Unlock()

	if len(a.Messages) != len(b.Messages) {
		return false
	}

	sort.Slice(a.Messages, compareMessagesFn(a.Messages))
	sort.Slice(b.Messages, compareMessagesFn(b.Messages))

	for i := range a.Messages {
		if a.Messages[i].ID != b.Messages[i].ID {
			return false
		}
	}

	return true
}
