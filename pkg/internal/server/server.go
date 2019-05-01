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

package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/httpmessage"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/round"
	"github.com/rstefan1/bimodal-multicast/pkg/peer"
)

type HTTP struct {
	server            *http.Server
	peerBuffer        []peer.Peer
	msgBuffer         *buffer.MessageBuffer
	gossipRoundNumber *round.GossipRound
	logger            *log.Logger
}

func getGossipHandler(w http.ResponseWriter, r *http.Request, msgBuffer *buffer.MessageBuffer, gossipRound *round.GossipRound, logger *log.Logger) {
	decoder := json.NewDecoder(r.Body)
	var t httpmessage.HTTPGossip
	err := decoder.Decode(&t)
	if err != nil {
		logger.Printf("Error at decoding HTTPGossip message in HTTP Server: %s", err)
		return
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
			logger.Printf("Error at marshal HTTPSolicitation in HTTP Server: %s", err)
			return
		}
		path := fmt.Sprintf("http://%s:%s/solicitation", t.Addr, t.Port)
		_, err = http.Post(path, "json", bytes.NewBuffer(jsonSolicitation))
		if err != nil {
			logger.Printf("Error at sending HTTPSoliccitation message in HTTP Server: %s", err)
			return
		}
	}
}

func getSolicitationHandler(w http.ResponseWriter, r *http.Request, msgBuffer *buffer.MessageBuffer, gossipRound *round.GossipRound, logger *log.Logger) {
	decoder := json.NewDecoder(r.Body)
	var t httpmessage.HTTPSolicitation
	err := decoder.Decode(&t)
	if err != nil {
		logger.Printf("Error at decoding HTTPSolicitation message in HTTP Server: %s", err)
		return
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
		logger.Printf("Error at marshal HTTPSynchronization message in HTTP Server: %s", err)
		return
	}
	path := fmt.Sprintf("http://%s:%s/synchronization", t.Addr, t.Port)
	_, err = http.Post(path, "json", bytes.NewBuffer(jsonSynchronization))
	if err != nil {
		logger.Printf("Error at sending HTTPSynchronization message in HTTP Server: %s", err)
		return
	}
}

func getSynchronizationHandler(w http.ResponseWriter, r *http.Request, msgBuffer *buffer.MessageBuffer, gossipRound *round.GossipRound, logger *log.Logger) {
	decoder := json.NewDecoder(r.Body)
	var t httpmessage.HTTPSynchronization
	err := decoder.Decode(&t)
	if err != nil {
		logger.Printf("Error at decoding HTTPSynchronization message in HTTP Server: %s", err)
		return
	}

	host := strings.Split(r.Host, ":")
	hostAddr := host[0]
	hostPort := host[1]

	rcvMsgBuf := t.Messages
	for _, m := range rcvMsgBuf.Messages {
		if !msgBuffer.AlreadyExists(m) {
			msgBuffer.AddMessage(m)
			logger.Printf("BMMC %s:%s synced buffer with message %s in round %d", hostAddr, hostPort, m.ID, gossipRound.GetNumber())
		}
	}
}

func startHTTPServer(s *HTTP) error {
	s.logger.Printf("HTTP Server listening at %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Printf("Unable to start HTTP server: %s", err)
		return err
	}
	return nil
}

func gracefullyShutdown(s *HTTP) {
	if err := s.server.Shutdown(context.TODO()); err != nil {
		s.logger.Printf("Unable to shutdown HTTP server properly: %s", err)
	}
}

func New(cfg Config) *HTTP {
	if cfg.Logger == nil {
		cfg.Logger = log.New(os.Stdout, "", 0)
	}

	return &HTTP{
		server: &http.Server{
			Addr: fmt.Sprintf("%s:%s", cfg.Addr, cfg.Port),
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch path := r.URL.Path; path {
				case "/gossip":
					getGossipHandler(w, r, cfg.MsgBuf, cfg.GossipRound, cfg.Logger)
				case "/solicitation":
					getSolicitationHandler(w, r, cfg.MsgBuf, cfg.GossipRound, cfg.Logger)
				case "/synchronization":
					getSynchronizationHandler(w, r, cfg.MsgBuf, cfg.GossipRound, cfg.Logger)
				}
			}),
		},
		peerBuffer:        cfg.PeerBuf,
		msgBuffer:         cfg.MsgBuf,
		gossipRoundNumber: cfg.GossipRound,
		logger:            cfg.Logger,
	}
}

func (h *HTTP) Start(stop <-chan struct{}) error {
	errChan := make(chan error)

	go func() {
		err := startHTTPServer(h)
		errChan <- err
	}()

	go func() {
		<-stop
		gracefullyShutdown(h)
	}()

	// return <-errChan
	return nil
}
