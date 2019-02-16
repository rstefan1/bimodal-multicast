package buffer

import (
	"sort"
	"sync"
)

type Digest struct {
	id string
}

type DigestBuffer struct {
	listDigests []Digest
}

func WrapDigestBuffer(digestSlice []string) DigestBuffer {
	var digestBuffer DigestBuffer
	for _, d := range digestSlice {
		digestBuffer.listDigests = append(digestBuffer.listDigests, Digest{id: d})
	}
	return digestBuffer
}

func (digestBuffer DigestBuffer) UnwrapDigestBuffer() []string {
	var digestSlice []string
	for _, d := range digestBuffer.listDigests {
		digestSlice = append(digestSlice, d.id)
	}
	return digestSlice
}

func compareFn(digest []Digest) func(int, int) bool {
	return func(i, j int) bool {
		return digest[i].id <= digest[j].id
	}
}

// HasSameDigests returns true if given digest are same
func (a DigestBuffer) SameDigests(b DigestBuffer) bool {
	if len(a.listDigests) != len(b.listDigests) {
		return false
	}

	sort.Slice(a.listDigests, compareFn(a.listDigests))
	sort.Slice(b.listDigests, compareFn(b.listDigests))

	for i := range a.listDigests {
		if a.listDigests[i].id != b.listDigests[i].id {
			return false
		}
	}

	return true
}

// GetMissingDigests returns the disjunction between the digest buffers
// digestBufferA - digestBufferB
func (a DigestBuffer) GetMissingDigests(b DigestBuffer) DigestBuffer {
	missingDigestBuffer := DigestBuffer{
		listDigests: []Digest{},
	}

	sort.Slice(a.listDigests, compareFn(a.listDigests))
	sort.Slice(b.listDigests, compareFn(b.listDigests))

	mapB := make(map[string]bool)
	for i := range b.listDigests {
		mapB[b.listDigests[i].id] = true
	}

	for i := range a.listDigests {
		if _, e := mapB[a.listDigests[i].id]; !e {
			missingDigestBuffer.listDigests = append(missingDigestBuffer.listDigests, a.listDigests[i])
		}
	}

	return missingDigestBuffer
}

// ContainsDigest check if digest buffer contains the given digest
func (digestBuffer DigestBuffer) ContainsDigest(digest Digest) bool {
	for _, d := range digestBuffer.listDigests {
		if d.id == digest.id {
			return true
		}
	}
	return false
}

// Length returns the length of digest buffer
func (digestBuffer DigestBuffer) Length() int {
	return len(digestBuffer.listDigests)
}

// GetMissingMessageBuffer returns messages buffer from given digest buffer
func (digestBuffer DigestBuffer) GetMissingMessageBuffer(msgBuffer MessageBuffer) MessageBuffer {
	missingMsgBuffer := MessageBuffer{
		mux: &sync.Mutex{},
	}

	msgBuffer.mux.Lock()
	for _, msg := range msgBuffer.listMessages {
		if digestBuffer.ContainsDigest(Digest{id: msg.id}) {
			missingMsgBuffer.listMessages = append(missingMsgBuffer.listMessages, msg)
		}
	}
	msgBuffer.mux.Unlock()

	// TODO what happens if the buffer does not have digests anymore
	return missingMsgBuffer
}
