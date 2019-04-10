/*
Copyright 2019 Robert Andrei STEFAN

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/httpmessage"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

type HTTP struct {
	server     *http.Server
	peerBuffer []peer.Peer
	msgBuffer  *buffer.MessageBuffer
}

func getGossipHandler(w http.ResponseWriter, r *http.Request, msgBuffer *buffer.MessageBuffer) {
	decoder := json.NewDecoder(r.Body)
	var t httpmessage.HTTPGossip
	err := decoder.Decode(&t)
	if err != nil {
		log.Fatal("Error at decoding HTTPGossip message: ", err)
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
			Digests:     *missingDigestBuffer,
		}
		jsonSolicitation, err := json.Marshal(solicitation)
		if err != nil {
			log.Fatal("Error at marshal HTTPSolicitation: ", err)
			panic(err)
		}
		path := fmt.Sprintf("http://%s:%s/solicitation", t.Addr, t.Port)
		_, err = http.Post(path, "json", bytes.NewBuffer(jsonSolicitation))
		if err != nil {
			log.Fatal("Error at sending HTTPSoliccitation message: ", err)
			panic(err)
		}
	}
}

func getSolicitationHandler(w http.ResponseWriter, r *http.Request, msgBuffer *buffer.MessageBuffer) {
	decoder := json.NewDecoder(r.Body)
	var t httpmessage.HTTPSolicitation
	err := decoder.Decode(&t)
	if err != nil {
		log.Fatal("Error at decoding HTTPSolicitation message: ", err)
		panic(err)
	}

	missingDigestBuffer := t.Digests
	missingMsgBuffer := missingDigestBuffer.GetMissingMessageBuffer(msgBuffer)

	host := strings.Split(r.Host, ":")
	hostAddr := host[0]
	hostPort := host[1]

	// send synchronization message
	synchronization := httpmessage.HTTPSynchronization{
		Addr:     hostAddr,
		Port:     hostPort,
		Messages: *missingMsgBuffer,
	}

	jsonSynchronization, err := json.Marshal(synchronization)
	if err != nil {
		log.Fatal("Error at marshal HTTPSynchronization message: ", err)
		panic(err)
	}
	path := fmt.Sprintf("http://%s:%s/synchronization", t.Addr, t.Port)
	_, err = http.Post(path, "json", bytes.NewBuffer(jsonSynchronization))
	if err != nil {
		log.Fatal("Error at sending HTTPSynchronization message: ", err)
		panic(err)
	}
}

func getSynchronizationHandler(w http.ResponseWriter, r *http.Request, msgBuffer *buffer.MessageBuffer) {
	decoder := json.NewDecoder(r.Body)
	var t httpmessage.HTTPSynchronization
	err := decoder.Decode(&t)
	if err != nil {
		log.Fatal("Error at decoding HTTPSynchronization message: ", err)
		panic(err)
	}

	rcvMsgBuf := t.Messages
	for _, m := range rcvMsgBuf.Messages {
		msgBuffer.AddMessage(m)
	}
}

func startHTTPServer(s *http.Server) error {
	log.Printf("HTTP Server listening at %s", s.Addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err, "unable to start HTTP server")
		return err
	}
	return nil
}

func gracefullyShutdown(s *http.Server) {
	if err := s.Shutdown(context.TODO()); err != nil {
		log.Fatal(err, "unable to shutdown HTTP server properly")
	}
}

func New(cfg Config) *HTTP {
	return &HTTP{
		server: &http.Server{
			Addr: fmt.Sprintf("%s:%s", cfg.Addr, cfg.Port),
			//Handler:getSynchronizationHandler(cfg.MsgBuf),
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch path := r.URL.Path; path {
				case "/gossip":
					getGossipHandler(w, r, cfg.MsgBuf)
				case "/solicitation":
					getSolicitationHandler(w, r, cfg.MsgBuf)
				case "/synchronization":
					getSynchronizationHandler(w, r, cfg.MsgBuf)
				}
			}),
		},
		peerBuffer: cfg.PeerBuf,
		msgBuffer:  cfg.MsgBuf,
	}
}

func (h *HTTP) Start(stop <-chan struct{}) error {
	errChan := make(chan error)

	go func() {
		err := startHTTPServer(h.server)
		errChan <- err
	}()

	go func() {
		<-stop
		gracefullyShutdown(h.server)
	}()

	// return <-errChan
	return nil
}
