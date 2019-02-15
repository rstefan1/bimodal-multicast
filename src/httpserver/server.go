package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	. "github.com/rstefan1/bimodal-multicast/src/internal"
	"github.com/rstefan1/bimodal-multicast/src/options"
)

var (
	nodeBuffer *[]Node
	msgBuffer  *[]Message
)

var gossipHandler = func(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error at parse form: %v", err)
		return
	}

	u, err := url.Parse(r.URL.Path)
	if err != nil {
		panic(err)
		return
	}

	msgDigest := Digest(*msgBuffer)
	gossipDigest := u.Query()["msg_id[]"]
	missingDigest := GetMissingDigest(msgDigest, gossipDigest)

	if len(missingDigest) > 0 {
		path := fmt.Sprintf("%s/solicitation", r.Host)
		v := url.Values{}
		for _, m := range missingDigest {
			v.Add("msg_id", m)
		}

		_, err := http.PostForm(path, v)
		if err != nil {
			panic(err)
			return
		}
	}
}

var solicitationHandler = func(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error at parse form: %v", err)
		return
	}

	u, err := url.Parse(r.URL.Path)
	if err != nil {
		panic(err)
		return
	}

	missingDigest := u.Query()["msg_id[]"]

	missingMsg := GetMissingMessages(*msgBuffer, missingDigest)
	path := fmt.Sprintf("%s/synchronization", r.Host)
	v := url.Values{}
	for _, msg := range missingMsg {
		m, err := json.Marshal(msg)
		if err != nil {
			panic(err)
			return
		}
		v.Add("msg_id", string(m))
	}

	_, err = http.PostForm(path, v)
	if err != nil {
		panic(err)
		return
	}
}

var synchronizationHandler = func(w http.ResponseWriter, r *http.Request) {

}

func Start(nodeBuf *[]Node, msgBuf *[]Message) {
	nodeBuffer = nodeBuf
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
