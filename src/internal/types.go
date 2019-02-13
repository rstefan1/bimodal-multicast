package internal

type Message struct {
	msg         string
	timestamp   string
	GossipCount int
}

type Node struct {
	Addr string
}

type GossipMessage struct {
	roundNumber int
	msg         []Message
}
