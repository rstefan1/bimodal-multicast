package buffer

import (
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
func NewMessageBuffer() MessageBuffer {
	return MessageBuffer{
		Mux: &sync.Mutex{},
	}
}

// Digest returns a slice with ID of messages from given buffer
func (msgBuffer MessageBuffer) DigestBuffer() DigestBuffer {
	var digestBuffer DigestBuffer

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
