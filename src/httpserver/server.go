package httpserver

import (
	"fmt"
	"net/http"

	"github.com/rstefan1/bimodal-multicast/src/options"
)

var httpHandler = func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// Call ParseForm() to parse the raw query
		// and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		msgType := r.FormValue("value")
		switch msgType {
		case "gossip":
			rcvGossipMsg(r)
		case "retransmit":
			rcvSolicitRetransmition(r)
		default:
			// TODO
			// return
			// error
			return
		}
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}

func Start() {
	go func() {
		http.HandleFunc("/", httpHandler)
		if err := http.ListenAndServe(options.HTTPAddr, nil); err != nil {
			panic(err)
		}
	}()
}
