package buffer

import "sync"

type Message struct {
	id          string
	msg         string
	GossipCount int
}

type MessageBuffer struct {
	listMessages []Message
	mux          *sync.Mutex
}

func (msgBuffer MessageBuffer) AddMutex(mx *sync.Mutex) {
	msgBuffer.mux = mx
}

func (msgBuffer MessageBuffer) UnwrapMessageBuffer() []Message {
	return msgBuffer.listMessages
}

// Digest returns a slice with id of messages from given buffer
func (msgBuffer MessageBuffer) DigestBuffer() DigestBuffer {
	// TODO write unit tests for this func

	var digestBuffer DigestBuffer

	msgBuffer.mux.Lock()
	for _, b := range msgBuffer.listMessages {
		digestBuffer.listDigests = append(digestBuffer.listDigests, Digest{id: b.id})
	}
	msgBuffer.mux.Unlock()

	return digestBuffer
}

// AddMessage adds message in message buffer
func (msgBuffer MessageBuffer) AddMessage(msg Message) MessageBuffer {
	msgBuffer.mux.Lock()
	msgBuffer.listMessages = append(msgBuffer.listMessages, msg)
	msgBuffer.mux.Unlock()
	return msgBuffer
}

// IncrementGossipCount increments gossip countfor each message from message
// buffer
func (msgBuffer MessageBuffer) IncrementGossipCount() {
	msgBuffer.mux.Lock()
	for _, m := range msgBuffer.listMessages {
		m.GossipCount++
	}
	msgBuffer.mux.Unlock()
}
