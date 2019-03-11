package protocol

import (
	. "github.com/rstefan1/bimodal-multicast/src/internal/buffer"
	. "github.com/rstefan1/bimodal-multicast/src/internal/peer"
)

//var (
//	// buffer with addresses of nodes in system
//	peerBuffer *[]Peer
//	// buffer with gossip messages
//	msgBuffer *MessageBuffer
//	// round numer
//	roundNumber = 1
//	// beta is the expected fanout for gossip
//	beta = 0.5
//)

// TODO implement a func that read a yml file with all node (addr & port)

//func randomlySelectedPeer() Peer {
//	// TODO check if node wasn't selected before
//	r := rand.Intn(len(*peerBuffer))
//	return (*peerBuffer)[r]
//}

//// GossipRound is the gossip round that runs every 100ms in out implementation
//func gossipRound() {
//	for {
//		// increment the round number
//		roundNumber++
//
//		gMsg := GossipMessage{
//			RoundNumber: roundNumber,
//			Digest:      (*msgBuffer).DigestBuffer(),
//		}
//
//		length := int(beta * float64(len(*peerBuffer)) / float64(roundNumber))
//		for i := 0; i < length; i++ {
//			dest := randomlySelectedPeer()
//			send.Gossip(dest, gMsg)
//		}
//
//		(*msgBuffer).IncrementGossipCount()
//
//		time.Sleep(100 * time.Millisecond)
//	}
//
//	// TODO discard messages for which gossip_count
//	// exceeds G, the garbage-collection limit
//}

func Start(peerBuf *[]Peer, msgBuf *MessageBuffer) {
	//peerBuffer = peerBuf
	//msgBuffer = msgBuf
	//gossipRound()
}
