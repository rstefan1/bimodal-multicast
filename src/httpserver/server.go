package httpserver

import (
	"fmt"
	"net/http"
	"net/url"

	. "github.com/rstefan1/bimodal-multicast/src/internal/buffer"
	. "github.com/rstefan1/bimodal-multicast/src/internal/peer"
	"github.com/rstefan1/bimodal-multicast/src/options"
)

var (
	peerBuffer *[]Peer
	msgBuffer  *MessageBuffer
)

var gossipHandler = func(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error at parse form: %v", err)
		return
	}

	u, err := url.Parse(r.URL.Path)
	if err != nil {
		panic(err)
	}

	gossipDigestBuffer := WrapDigestBuffer(u.Query()["msg_id[]"])
	msgDigestBuffer := (*msgBuffer).DigestBuffer()
	missingDigestBuffer := msgDigestBuffer.GetMissingDigests(gossipDigestBuffer)

	if missingDigestBuffer.Length() > 0 {
		path := fmt.Sprintf("%s/solicitation", r.Host)
		v := url.Values{}
		for _, m := range missingDigestBuffer.UnwrapDigestBuffer() {
			v.Add("msg_id", m)
		}

		_, err := http.PostForm(path, v)
		if err != nil {
			panic(err)
		}
	}
}

var solicitationHandler = func(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error at parse form: %v", err)
		return
	}

	// 	u, err := url.Parse(r.URL.Path)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// missingDigestBuffer := WrapDigestBuffer(u.Query()["msg_id[]"])
	// missingMsgBuffer := missingDigestBuffer.GetMissingMessageBuffer(*msgBuffer)
	// missingMsgList := missingMsgBuffer.UnwrapMessageBuffer()
	// TODO serialize only missingMsgList

	// path := fmt.Sprintf("%s/synchronization", r.Host)
	// for _, msg := range missingMsgList {
	// _, err = http.Post(path, msg.Marhsal())
	// }

	// _, err = http.PostForm(path, v)
	// if err != nil {
	// panic(err)
	// }
}

var synchronizationHandler = func(w http.ResponseWriter, r *http.Request) {

}

func Start(peerBuf *[]Peer, msgBuf *MessageBuffer) {
	peerBuffer = peerBuf
	msgBuffer = msgBuf

	go func() {
		http.HandleFunc("/gossip", gossipHandler)
		if err := http.ListenAndServe(options.HTTPAddr, nil); err != nil {
			panic(err)
		}
	}()
	go func() {
		http.HandleFunc("/solicitation", solicitationHandler)
		if err := http.ListenAndServe(options.HTTPAddr, nil); err != nil {
			panic(err)
		}
	}()
	go func() {
		http.HandleFunc("/synchronization", synchronizationHandler)
		if err := http.ListenAndServe(options.HTTPAddr, nil); err != nil {
			panic(err)
		}
	}()
}
