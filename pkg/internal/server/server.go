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
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/httputil"
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
	gossipDigestBuffer, tAddr, tPort, tRoundNumber, err := httputil.ReceiveGossip(r)
	if err != nil {
		logger.Printf("%s", err)
		return
	}

	msgDigestBuffer := msgBuffer.DigestBuffer()
	missingDigestBuffer := gossipDigestBuffer.GetMissingDigests(msgDigestBuffer)

	host := strings.Split(r.Host, ":")
	hostAddr := host[0]
	hostPort := host[1]

	if missingDigestBuffer.Length() > 0 {
		solicitationMsg := httputil.HTTPSolicitation{
			Addr:        hostAddr,
			Port:        hostPort,
			RoundNumber: tRoundNumber,
			Digests:     *missingDigestBuffer,
		}

		err = httputil.SendSolicitation(solicitationMsg, tAddr, tPort)
		if err != nil {
			logger.Printf("%s", err)
			return
		}
	}
}

func getSolicitationHandler(w http.ResponseWriter, r *http.Request, msgBuffer *buffer.MessageBuffer, gossipRound *round.GossipRound, logger *log.Logger) {
	missingDigestBuffer, tAddr, tPort, _, err := httputil.ReceiveSolicitation(r)
	if err != nil {
		logger.Printf("%s", err)
		return
	}
	missingMsgBuffer := missingDigestBuffer.GetMissingMessageBuffer(msgBuffer)

	host := strings.Split(r.Host, ":")
	hostAddr := host[0]
	hostPort := host[1]

	synchronizationMsg := httputil.HTTPSynchronization{
		Addr:     hostAddr,
		Port:     hostPort,
		Messages: *missingMsgBuffer,
	}

	err = httputil.SendSynchronization(synchronizationMsg, tAddr, tPort)
	if err != nil {
		logger.Printf("%s", err)
		return
	}
}

func getSynchronizationHandler(w http.ResponseWriter, r *http.Request, msgBuffer *buffer.MessageBuffer, gossipRound *round.GossipRound, logger *log.Logger) {
	rcvMsgBuf, _, _, err := httputil.ReceiveSynchronization(r)
	if err != nil {
		logger.Printf("%s", err)
		return
	}

	host := strings.Split(r.Host, ":")
	hostAddr := host[0]
	hostPort := host[1]

	for _, m := range rcvMsgBuf.Messages {
		msgBuffer.AddMessage(m)
		logger.Printf("BMMC %s:%s synced buffer with message %s in round %d", hostAddr, hostPort, m.ID, gossipRound.GetNumber())
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
					// in production, cfg.Loss = 0
					if rand.Float64() >= cfg.Loss {
						getGossipHandler(w, r, cfg.MsgBuf, cfg.GossipRound, cfg.Logger)
					}
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
