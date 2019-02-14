package internal

type Message struct {
	msg         string
	id          string
	GossipCount int
}

type Node struct {
	Addr string
}

type GossipMessage struct {
	RoundNumber int
	Digest      []string
}
