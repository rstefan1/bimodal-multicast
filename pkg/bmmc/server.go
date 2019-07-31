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

package bmmc

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/callback"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/httputil"
)

func gossipHandler(_ http.ResponseWriter, r *http.Request, cfg Config) {
	gossipDigestBuffer, tAddr, tPort, tRoundNumber, err := httputil.ReceiveGossip(r)
	if err != nil {
		cfg.Logger.Printf("%s", err)
		return
	}

	msgDigestBuffer := cfg.MsgBuf.DigestBuffer()
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
			cfg.Logger.Printf("%s", err)
			return
		}
	}
}

func solicitationHandler(_ http.ResponseWriter, r *http.Request, cfg Config) {
	missingDigestBuffer, tAddr, tPort, _, err := httputil.ReceiveSolicitation(r)
	if err != nil {
		cfg.Logger.Printf("%s", err)
		return
	}
	missingMsgBuffer := missingDigestBuffer.GetMissingMessageBuffer(cfg.MsgBuf)

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
		cfg.Logger.Printf("%s", err)
		return
	}
}

func synchronizationHandler(_ http.ResponseWriter, r *http.Request, cfg Config) {
	host := strings.Split(r.Host, ":")
	hostAddr := host[0]
	hostPort := host[1]

	rcvMsgBuf, _, _, err := httputil.ReceiveSynchronization(r)
	if err != nil {
		cfg.Logger.Printf("BMMC %s:%s Error at receiving synchronization error: %s", hostAddr, hostPort, err)
		return
	}

	for _, m := range rcvMsgBuf.Messages {
		err = cfg.MsgBuf.AddMessage(m)
		if err != nil {
			cfg.Logger.Printf("BMMC %s:%s error at syncing buffer with message %s in round %d: %s", hostAddr, hostPort, m.ID, cfg.GossipRound.GetNumber(), err)
		} else {
			cfg.Logger.Printf("BMMC %s:%s synced buffer with message %s in round %d", hostAddr, hostPort, m.ID, cfg.GossipRound.GetNumber())

			// run callback function for messages with a callback registered
			if m.CallbackType != callback.NOCALLBACK {
				err = cfg.DefaultCallbacks.RunDefaultCallbacks(m, cfg.PeerBuf, cfg.Logger)
				if err != nil {
					cfg.Logger.Printf("Error at calling default callback at %s:%s for message %s in round %d", hostAddr, hostPort, m.ID, cfg.GossipRound.GetNumber())
				}

				err = cfg.CustomCallbacks.RunCustomCallbacks(m, cfg.Logger)
				if err != nil {
					cfg.Logger.Printf("Error at calling custom callback at %s:%s for message %s in round %d", hostAddr, hostPort, m.ID, cfg.GossipRound.GetNumber())
				}
			}
		}
	}
}

func startServer(s *Server) error {
	s.logger.Printf("Server listening at %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Printf("Unable to start  server: %s", err)
		return err
	}
	return nil
}

func gracefullyShutdown(s *Server) {
	if err := s.server.Shutdown(context.TODO()); err != nil {
		s.logger.Printf("Unable to shutdown server properly: %s", err)
	}
}

func (s *BMMC) startServer(stop <-chan struct{}) error {
	errChan := make(chan error)

	srv := &Server{
		server: &http.Server{
			Addr: fmt.Sprintf("0.0.0.0:%s", cfg.Port),
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch path := r.URL.Path; path {
				case "/gossip":
					gossipHandler(w, r, cfg)
				case "/solicitation":
					solicitationHandler(w, r, cfg)
				case "/synchronization":
					synchronizationHandler(w, r, cfg)
				}
			}),
		},
		peerBuffer:        cfg.PeerBuf,
		msgBuffer:         cfg.MsgBuf,
		gossipRoundNumber: cfg.GossipRound,
		logger:            cfg.Logger,
	}

	go func() {
		err := startServer(srv)
		errChan <- err
	}()

	go func() {
		<-stop
		gracefullyShutdown(s)
	}()

	// return <-errChan
	return nil
}
