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
	"sort"
	"sync"
	"time"
)

// Message is a message
type Message struct {
	ID           string      `json:"message_ID"`
	Msg          interface{} `json:"message_msg"`
	CallbackType string      `json:"message_callback_type"`
	GossipCount  int         `json:"message_gossip_count"`
}

// MessageBuffer is the buffer with messages
type MessageBuffer struct {
	Messages []Message   `json:"message_buffer_list"`
	Mux      *sync.Mutex `json:"message_buffer_mux"`
}

// NewMessage creates new message
func NewMessage(m interface{}, callbackType string) Message {
	return Message{
		ID:           fmt.Sprintf("%s%d", time.Now().Format("20060102150405"), rand.Int31()),
		Msg:          m,
		CallbackType: callbackType,
		GossipCount:  0,
	}
}

// NewMessageBuffer creates new MessageBuffer
func NewMessageBuffer() *MessageBuffer {
	return &MessageBuffer{
		Mux: &sync.Mutex{},
	}
}

// Length return the length of message buffer
func (msgBuffer *MessageBuffer) Length() int {
	msgBuffer.Mux.Lock()
	defer msgBuffer.Mux.Unlock()

	l := len(msgBuffer.Messages)
	return l
}

// DigestBuffer returns a DigestBuffer from given buffer
func (msgBuffer *MessageBuffer) DigestBuffer() *DigestBuffer {
	digestBuffer := &DigestBuffer{}

	msgBuffer.Mux.Lock()
	defer msgBuffer.Mux.Unlock()

	for _, b := range msgBuffer.Messages {
		digestBuffer.Digests = append(digestBuffer.Digests, Digest{ID: b.ID})
	}

	return digestBuffer
}

// alreadyExists return true if the message already exists in message buffer
func (msgBuffer *MessageBuffer) alreadyExists(msg Message) bool {
	// TODO: Solve the following situatiion:
	// Important! Whoever calls this function must LOCK the buffer
	for i := range msgBuffer.Messages {
		if msgBuffer.Messages[i].ID == msg.ID {
			return true
		}
	}
	return false
}

// AddMessage adds message in message buffer
func (msgBuffer *MessageBuffer) AddMessage(msg Message) error {
	msgBuffer.Mux.Lock()
	defer msgBuffer.Mux.Unlock()

	if msgBuffer.alreadyExists(msg) {
		return fmt.Errorf("message %s already exists in buffer message", msg.Msg)
	}

	msgBuffer.Messages = append(msgBuffer.Messages, msg)
	return nil
}

// UnwrapMessageBuffer wraps a message buffer
// TODO write a test for this function
func (msgBuffer *MessageBuffer) UnwrapMessageBuffer() []interface{} {
	msgBuffer.Mux.Lock()
	defer msgBuffer.Mux.Unlock()

	messages := make([]interface{}, len(msgBuffer.Messages))
	for i := range msgBuffer.Messages {
		messages[i] = msgBuffer.Messages[i].Msg
	}

	return messages
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

// SameMessages returns true if given message buffer contains same messages
func (msgBuffer *MessageBuffer) SameMessages(b *MessageBuffer) bool {
	msgBuffer.Mux.Lock()
	b.Mux.Lock()
	defer msgBuffer.Mux.Unlock()
	defer b.Mux.Unlock()

	if len(msgBuffer.Messages) != len(b.Messages) {
		return false
	}

	sort.Slice(msgBuffer.Messages, compareMessagesFn(msgBuffer.Messages))
	sort.Slice(b.Messages, compareMessagesFn(b.Messages))

	for i := range msgBuffer.Messages {
		if msgBuffer.Messages[i].ID != b.Messages[i].ID {
			return false
		}
	}

	return true
}
