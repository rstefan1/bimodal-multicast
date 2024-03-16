/*
Copyright 2024 Robert Andrei STEFAN

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

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/rstefan1/bimodal-multicast/pkg/bmmc"
)

// Server decorates an HTTP Server.
type Server struct {
	*http.Server
}

// NewServer creates a new http server.
func NewServer(b *bmmc.BMMC, addr, port string, log *slog.Logger) *Server {
	return &Server{
		&http.Server{
			Addr: fmt.Sprintf("%s:%s", addr, port),
			Handler: http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				switch path := r.URL.Path; path {
				case bmmc.GossipRoute:
					body, err := io.ReadAll(r.Body)
					if err != nil {
						log.Error("unable to read gossip message body", err)

						return
					}

					b.GossipHandler(body)

				case bmmc.SolicitationRoute:
					body, err := io.ReadAll(r.Body)
					if err != nil {
						log.Error("unable to read solicitation message body", err)

						return
					}

					b.SolicitationHandler(body)

				case bmmc.SynchronizationRoute:
					body, err := io.ReadAll(r.Body)
					if err != nil {
						log.Error("unable to read synchronization message body", err)

						return
					}

					b.SynchronizationHandler(body)
				}
			}),
			ReadHeaderTimeout: 30 * time.Second, //nolint: gomnd
		},
	}
}

// Start starts the http server.
func (s *Server) Start(stop <-chan struct{}, log *slog.Logger) error { //nolint: unparam
	errChan := make(chan error)

	go func() {
		log.Info("starting server", "addr", s.Addr)

		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Error("unable to start server", "addr", s.Addr, "err", err)

			errChan <- err

			return
		}

		errChan <- nil
	}()

	go func() {
		<-stop

		err := s.gracefullyShutdown()
		if err != nil {
			log.Error("unable to stop server", "err", err, "addr", s.Addr)
		}

		log.Info("server stopped", "addr", s.Addr)
	}()

	// TODO return err chan
	// return <-errChan

	return nil
}

func (s *Server) gracefullyShutdown() error {
	return s.Shutdown(context.TODO()) //nolint: wrapcheck
}
