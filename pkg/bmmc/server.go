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
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
)

const (
	startServerLogFmt       = "Starting server for  %s"
	stopServerLogFmt        = "End of server round from %s"
	unableStartServerLogFmt = "Unable to start  server: %s"
	unableStopServerLogFmt  = "Unable to shutdown server properly: %s"

	gossipHandlerErrLogFmt          = "Error in gossip handler: %s"
	solicitationHandlerErrLogFmt    = "Error in solicitation handler: %s"
	synchronizationHandlerErrLogFmt = "Error in synchronization handler: %s"

	syncBufferLogErrFmt = "BMMC %s error at syncing buffer with message %s in round %d: %s"
	bufferSyncedLogFmt  = "BMMC %s synced buffer with message %s in round %d"

	gossipRoute          = "/gossip"
	solicitationRoute    = "/solicitation"
	synchronizationRoute = "/synchronization"
)

func (b *BMMC) gossipHandler(body []byte) {
	gossipDigest, p, roundNumber, err := b.receiveGossip(body)
	if err != nil {
		b.config.Logger.Printf("%s", err)

		return
	}

	digest := b.messageBuffer.Digest()
	missingDigest := buffer.MissingStrings(gossipDigest, digest)

	if len(missingDigest) > 0 {
		solicitationMsg := Solicitation{
			Host:        b.config.Host.String(),
			RoundNumber: roundNumber,
			Digest:      missingDigest,
		}

		if err = b.sendSolicitation(solicitationMsg, p); err != nil {
			b.config.Logger.Printf(gossipHandlerErrLogFmt, err)

			return
		}
	}
}

func (b *BMMC) gossipHTTPHandler(_ http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		b.config.Logger.Printf("%s", err)

		return
	}

	b.gossipHandler(body)
}

func (b *BMMC) solicitationHandler(body []byte) {
	missingDigest, p, _, err := b.receiveSolicitation(body)
	if err != nil {
		b.config.Logger.Printf(solicitationHandlerErrLogFmt, err)

		return
	}

	missingElements := b.messageBuffer.ElementsFromIDs(missingDigest)

	synchronizationMsg := Synchronization{
		Host:     b.config.Host.String(),
		Elements: missingElements,
	}

	if err = b.sendSynchronization(synchronizationMsg, p); err != nil {
		b.config.Logger.Printf(solicitationHandlerErrLogFmt, err)

		return
	}
}

func (b *BMMC) solicitationHTTPHandler(_ http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		b.config.Logger.Printf("%s", err)

		return
	}

	b.solicitationHandler(body)
}

func (b *BMMC) synchronizationHandler(body []byte) {
	rcvElements, _, err := b.receiveSynchronization(body)
	if err != nil {
		b.config.Logger.Printf(synchronizationHandlerErrLogFmt, err)

		return
	}

	for _, m := range rcvElements {
		err = b.messageBuffer.Add(m)
		if err != nil {
			b.config.Logger.Printf(syncBufferLogErrFmt, b.config.Host.String(), m.ID, b.gossipRound.GetNumber(), err)
		} else {
			b.config.Logger.Printf(bufferSyncedLogFmt, b.config.Host.String(), m.ID, b.gossipRound.GetNumber())
			b.runCallbacks(m)
		}
	}
}

func (b *BMMC) synchronizationHTTPHandler(_ http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		b.config.Logger.Printf("%s", err)

		return
	}

	b.synchronizationHandler(body)
}

func (b *BMMC) gracefullyShutdown() {
	if err := b.server.Shutdown(context.TODO()); err != nil {
		b.config.Logger.Printf(unableStopServerLogFmt, err)
	}
}

// TODO: implement a generic server.
func (b *BMMC) newServer() *http.Server {
	host, _ := b.config.Host.(peer.HTTPPeer)

	return &http.Server{
		Addr: fmt.Sprintf("%s:%s", host.Addr(), host.Port()),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch path := r.URL.Path; path {
			case gossipRoute:
				b.gossipHTTPHandler(w, r)
			case solicitationRoute:
				b.solicitationHTTPHandler(w, r)
			case synchronizationRoute:
				b.synchronizationHTTPHandler(w, r)
			}
		}),
		ReadHeaderTimeout: b.config.ReadHeaderTimeout,
	}
}

func (b *BMMC) startServer(stop <-chan struct{}) error {
	errChan := make(chan error)

	go func() {
		b.config.Logger.Printf(startServerLogFmt, b.server.Addr)

		if err := b.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			b.config.Logger.Printf(unableStartServerLogFmt, err)
			errChan <- err

			return
		}

		errChan <- nil
	}()

	go func() {
		<-stop
		b.gracefullyShutdown()
		b.config.Logger.Printf(stopServerLogFmt, b.server.Addr)
	}()

	// TODO return err chan
	// return <-errChan

	return nil
}
