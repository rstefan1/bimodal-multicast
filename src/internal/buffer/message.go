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

func (msgBuffer MessageBuffer) AddMutex(mx *sync.Mutex) MessageBuffer {
	msgBuffer.Mux = mx
	return msgBuffer
}

func (msgBuffer MessageBuffer) UnwrapMessageBuffer() []Message {
	return msgBuffer.Messages
}

// Digest returns a slice with ID of messages from given buffer
func (msgBuffer MessageBuffer) DigestBuffer() DigestBuffer {
	var digestBuffer DigestBuffer

	msgBuffer.Mux.Lock()
	for _, b := range msgBuffer.Messages {
		digestBuffer.Digests = append(digestBuffer.Digests, Digest{ID: b.ID})
	}
	msgBuffer.Mux.Unlock()

	return digestBuffer
}

// AddMessage adds message in message buffer
func (msgBuffer MessageBuffer) AddMessage(msg Message) MessageBuffer {
	msgBuffer.Mux.Lock()
	msgBuffer.Messages = append(msgBuffer.Messages, msg)
	msgBuffer.Mux.Unlock()
	return msgBuffer
}

// IncrementGossipCount increments gossip countfor each message from message
// buffer
func (msgBuffer MessageBuffer) IncrementGossipCount() MessageBuffer {
	msgBuffer.Mux.Lock()
	for i := range msgBuffer.Messages {
		msgBuffer.Messages[i].GossipCount++
	}
	msgBuffer.Mux.Unlock()
	return msgBuffer
}
