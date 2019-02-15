package buffer

import "sync"

type Message struct {
	id          string
	msg         string
	GossipCount int
}

type MessageBuffer struct {
	listMessages []Message
}

var mux = &sync.Mutex{}

func (msgBuffer MessageBuffer) UnwrapMessageBuffer() []Message {
	return msgBuffer.listMessages
}

// Digest returns a slice with id of messages from given buffer
func (msgBuffer MessageBuffer) DigestBuffer() DigestBuffer {
	// TODO write unit tests for this func

	var digestBuffer DigestBuffer

	mux.Lock()
	for _, b := range msgBuffer.listMessages {
		digestBuffer.listDigests = append(digestBuffer.listDigests, Digest{id: b.id})
	}
	mux.Unlock()

	return digestBuffer
}

// AddMessage adds message in message buffer
func (msgBuffer MessageBuffer) AddMessage(msg Message) MessageBuffer {
	mux.Lock()
	msgBuffer.listMessages = append(msgBuffer.listMessages, msg)
	mux.Unlock()
	return msgBuffer
}
