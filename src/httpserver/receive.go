package httpserver

import "net/http"

func rcvGossipMsg(r *http.Request) {
	// TODO compare with contents of local message buffer
	// foreach missing_message, most recent first
	//      if this solicitation won't exceed limit on
	//      retransmitions per round
	//          send solicit_retransmition(roundNumber,
	//          msg.id) tp gmsg.sender

}

func rcvSolicitRetransmition(r *http.Request) {
	// TODO if I am no longer in msg.round, or if have exceeded limits for
	// this round
	//      ignore
	// else send(makeCopy(msg.solicitedMsgId), msg.Sender)

}
