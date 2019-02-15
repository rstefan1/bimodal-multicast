package buffer

import "sort"

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

// HasSameDigest returns true if given digest are same
func (a DigestBuffer) HasSameDigest(b DigestBuffer) bool {
	// TODO write unit tests for this func

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

// GetMissingDigest returns the disjunction between the digest buffers
func (a DigestBuffer) GetMissingDigest(b DigestBuffer) DigestBuffer {
	// TODO write unit tests for this func

	var missingDigestBuffer DigestBuffer
	var i = 0
	var j = 0

	sort.Slice(a.listDigests, compareFn(a.listDigests))
	sort.Slice(b.listDigests, compareFn(b.listDigests))

	for ; i < len(a.listDigests) && j < len(b.listDigests); i++ {
		if a.listDigests[i].id == b.listDigests[j].id {
			j++
			continue
		}
		missingDigestBuffer.listDigests = append(missingDigestBuffer.listDigests, a.listDigests[i])
	}
	missingDigestBuffer.listDigests = append(a.listDigests[i:], missingDigestBuffer.listDigests...)

	return missingDigestBuffer
}

// ContainsDigest check if digest buffer contains the given digest
func (digestBuffer DigestBuffer) ContainsDigest(digest Digest) bool {
	// TODO write unit tests for this func

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
	// TODO write unit tests for this func

	var missingMsgBuffer MessageBuffer

	mux.Lock()
	for _, msg := range msgBuffer.listMessages {
		if digestBuffer.ContainsDigest(Digest{id: msg.id}) {
			missingMsgBuffer.listMessages = append(missingMsgBuffer.listMessages, msg)
		}
	}
	mux.Unlock()

	return missingMsgBuffer
}
