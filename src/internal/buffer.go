package internal

import "sort"

func compareFn(slice []string) func(int, int) bool {
	return func(i, j int) bool {
		return slice[i] <= slice[j]
	}
}

// HasSomeDigest returns true if given digest are same
func HasSameDigest(a, b []string) bool {
	// TODO write unit tests for this func

	// id one is nil, the other must also be nil
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	sort.Slice(a, compareFn(a))
	sort.Slice(b, compareFn(b))

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// GetMissingMessagesFromDigest returns missing messages id from buffer b
func GetMissingDigest(a, b []string) []string {
	// TODO write unit tests for this func

	var (
		missingMsg []string
		i          = 0
		j          = 0
	)

	sort.Slice(a, compareFn(a))
	sort.Slice(b, compareFn(b))

	for ; i < len(a) && j < len(b); i++ {
		if a[i] == b[j] {
			j++
			continue
		}
		missingMsg = append(missingMsg, a[i])
	}
	missingMsg = append(a[i:], missingMsg...)

	return missingMsg
}

// Digest returns a slice with id of messages from given buffer
func Digest(buffer []Message) []string {
	var digest []string
	for _, b := range buffer {
		digest = append(digest, b.id)
	}
	return digest
}

func contains(slice []string, elem string) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}

// GetMissingMessages returns missing messages from given digest
func GetMissingMessages(buffer []Message, missingDigest []string) []Message {
	var missingMsg []Message

	for _, msg := range buffer {
		if contains(missingDigest, msg.id) {
			missingMsg = append(missingMsg, msg)
		}
	}

	return missingMsg
}
