package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	. "github.com/rstefan1/bimodal-multicast/src/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/src/internal/config"
	"github.com/rstefan1/bimodal-multicast/src/internal/httpmessage"
	. "github.com/rstefan1/bimodal-multicast/src/internal/peer"
)

var (
	peerBuffer *[]Peer
	msgBuffer  *MessageBuffer
)

var handler = func(w http.ResponseWriter, r *http.Request) {
	switch path := r.URL.Path; path {
	case "/gossip":
		gossipHandler(w, r)
	case "/solicitation":
		solicitationHandler(w, r)
	case "/synchronization":
		synchronizationHandler(w, r)
	}
}

var gossipHandler = func(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error at parse form: %v", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var t httpmessage.HTTPGossip
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	gossipDigestBuffer := t.Digests
	msgDigestBuffer := msgBuffer.DigestBuffer()
	missingDigestBuffer := gossipDigestBuffer.GetMissingDigests(msgDigestBuffer)

	host := strings.Split(r.Host, ":")
	hostAddr := host[0]
	hostPort := host[1]
	if missingDigestBuffer.Length() > 0 {
		// send solicitation request
		solicitation := httpmessage.HTTPSolicitation{
			Addr:        hostAddr,
			Port:        hostPort,
			RoundNumber: t.RoundNumber,
			Digests:     missingDigestBuffer,
		}
		jsonSolicitation, err := json.Marshal(solicitation)
		path := fmt.Sprintf("http://%s:%s/solicitation", t.Addr, t.Port)
		_, err = http.Post(path, "json", bytes.NewBuffer(jsonSolicitation))
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

func startHTTPServer(s *http.Server) error {
	log.Print("HTTP Server listening at ", s.Addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func gracefullShutdown(s *http.Server) {
	if err := s.Shutdown(context.TODO()); err != nil {
		log.Fatal(err, "unable to shutdown HTTP server properly")
	}
}

func Start(peerBuf *[]Peer, msgBuf *MessageBuffer, stop <-chan struct{}, cfg config.Config) error {
	peerBuffer = peerBuf
	msgBuffer = msgBuf
	errChan := make(chan error)

	s := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: http.HandlerFunc(handler),
	}

	go func() {
		err := startHTTPServer(s)
		errChan <- err
	}()

	go func() {
		<-stop
		gracefullShutdown(s)
	}()

	// return <-errChan
	return nil
}
