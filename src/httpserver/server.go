package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rstefan1/bimodal-multicast/src/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/src/internal/config"
	"github.com/rstefan1/bimodal-multicast/src/internal/httpmessage"
	"github.com/rstefan1/bimodal-multicast/src/internal/peer"
)

var (
	peerBuffer *[]peer.Peer
	msgBuffer  *buffer.MessageBuffer
)

var gossipHandler = func(w http.ResponseWriter, r *http.Request) {
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
		if err != nil {
			panic(err)
		}
		path := fmt.Sprintf("http://%s:%s/solicitation", t.Addr, t.Port)
		_, err = http.Post(path, "json", bytes.NewBuffer(jsonSolicitation))
		if err != nil {
			panic(err)
		}
	}
}

var solicitationHandler = func(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t httpmessage.HTTPSolicitation
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	missingDigestBuffer := t.Digests
	missingMsgBuffer := missingDigestBuffer.GetMissingMessageBuffer(*msgBuffer)

	host := strings.Split(r.Host, ":")
	hostAddr := host[0]
	hostPort := host[1]

	// send synchronization message
	synchronization := httpmessage.HTTPSynchronization{
		Addr:     hostAddr,
		Port:     hostPort,
		Messages: missingMsgBuffer,
	}

	jsonSynchronization, err := json.Marshal(synchronization)
	if err != nil {
		panic(err)
	}
	path := fmt.Sprintf("http://%s:%s/synchronization", t.Addr, t.Port)
	_, err = http.Post(path, "json", bytes.NewBuffer(jsonSynchronization))
	if err != nil {
		panic(err)
	}
}

var synchronizationHandler = func(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t httpmessage.HTTPSynchronization
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	rcvMsgBuf := t.Messages
	for _, m := range rcvMsgBuf.Messages {
		*msgBuffer = msgBuffer.AddMessage(m)
	}
}

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

func Start(peerBuf *[]peer.Peer, msgBuf *buffer.MessageBuffer, stop <-chan struct{}, cfg config.Config) error {
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
