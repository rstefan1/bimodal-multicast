package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/config"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/httpmessage"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
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
		msgBuffer.AddMessage(m)
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
	log.Printf("HTTP Server listening at %s", s.Addr)
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

func Start(cfg config.HTTPConfig, stop <-chan struct{}) error {
	peerBuffer = cfg.PeerBuf
	msgBuffer = cfg.MsgBuf
	errChan := make(chan error)

	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Addr, cfg.Port),
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
