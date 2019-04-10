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
	"sort"
	"sync"
)

type Digest struct {
	ID string `json:"digest_id"`
}

type DigestBuffer struct {
	Digests []Digest `json:"digest_buffer_list"`
}

// WrapDigestBuffer wraps []string into DigestBuffer
func WrapDigestBuffer(digestSlice []string) *DigestBuffer {
	digestBuffer := &DigestBuffer{}
	for _, d := range digestSlice {
		digestBuffer.Digests = append(digestBuffer.Digests, Digest{ID: d})
	}
	return digestBuffer
}

func compareDigestsFn(digest []Digest) func(int, int) bool {
	return func(i, j int) bool {
		return digest[i].ID <= digest[j].ID
	}
}

// SameDigests returns true if given digest are same
func (a *DigestBuffer) SameDigests(b *DigestBuffer) bool {
	if len(a.Digests) != len(b.Digests) {
		return false
	}

	sort.Slice(a.Digests, compareDigestsFn(a.Digests))
	sort.Slice(b.Digests, compareDigestsFn(b.Digests))

	for i := range a.Digests {
		if a.Digests[i].ID != b.Digests[i].ID {
			return false
		}
	}

	return true
}

// GetMissingDigests returns the disjunction between the digest buffers
// digestBufferA - digestBufferB
func (a *DigestBuffer) GetMissingDigests(b *DigestBuffer) *DigestBuffer {
	missingDigestBuffer := &DigestBuffer{
		Digests: []Digest{},
	}

	sort.Slice(a.Digests, compareDigestsFn(a.Digests))
	sort.Slice(b.Digests, compareDigestsFn(b.Digests))

	mapB := make(map[string]bool)
	for i := range b.Digests {
		mapB[b.Digests[i].ID] = true
	}

	for i := range a.Digests {
		if _, e := mapB[a.Digests[i].ID]; !e {
			missingDigestBuffer.Digests = append(missingDigestBuffer.Digests, a.Digests[i])
		}
	}

	return missingDigestBuffer
}

// ContainsDigest check if digest buffer contains the given digest
func (digestBuffer *DigestBuffer) ContainsDigest(digest Digest) bool {
	for _, d := range digestBuffer.Digests {
		if d.ID == digest.ID {
			return true
		}
	}
	return false
}

// Length returns the length of digest buffer
func (digestBuffer *DigestBuffer) Length() int {
	return len(digestBuffer.Digests)
}

// GetMissingMessageBuffer returns messages buffer from given digest buffer
func (digestBuffer *DigestBuffer) GetMissingMessageBuffer(msgBuffer *MessageBuffer) *MessageBuffer {
	missingMsgBuffer := &MessageBuffer{
		Mux: &sync.Mutex{},
	}

	msgBuffer.Mux.Lock()
	for _, msg := range msgBuffer.Messages {
		if digestBuffer.ContainsDigest(Digest{ID: msg.ID}) {
			missingMsgBuffer.Messages = append(missingMsgBuffer.Messages, msg)
		}
	}
	msgBuffer.Mux.Unlock()

	// TODO what happens if the buffer does not have digests anymore
	return missingMsgBuffer
}
