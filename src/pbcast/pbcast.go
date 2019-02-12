package pbcast

type Message struct {
	msg         string
	timestamp   string
	gossipCount int
}

var msgBuffer [100]Message
var bufferSize = 0
var myRoundNumber = 0

func receiveMessage() {

}

func addToMsgBuffer(msg Message) {
	// TODO implement this func
	bufferSize++
	msgBuffer[bufferSize] = msg
	msgBuffer[bufferSize].gossipCount = 0

}

func unreliablyMulticast(msg Message) {
	// TODO implement this func

}

func pbcast(msg Message) {
	// TODO implement this func
	addToMsgBuffer(msg)
	unreliablyMulticast(msg)

}

func firstReception(msg Message) {
	// TODO implement this func
	addToMsgBuffer(msg)
	// TODO deliver messages that are now in order; report gaps
	// after suitable delay;

}

func randomlySelectedNumber() int {
	// TODO implement this func
	return -1

}

func send(msg Message, dest int) {
	// TODO implement this func

}

// gossipRound runs every 100ms in out implementation
func gossipRound() {
	// TODO implement this func
	myRoundNumber++
	var gossipMsg Message
	// gossipMsg = <myRoundNumber, digest(msgBuffer)>

	var beta, N, R int

	for i := 0; i < beta*N/R; i = i + 1 {
		dest := randomlySelectedNumber()
		send(gossipMsg, dest)

	}

	for i := range msgBuffer {
		msgBuffer[i].gossipCount++

	}

	// TODO discard messages for which gossip_count
	// exceeds G, the garbage-collection limit

}
