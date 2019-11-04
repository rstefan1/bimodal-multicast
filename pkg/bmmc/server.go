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
	"net/http"
	"strings"

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

	syncBufferLogErrFmt = "BMMC %s:%s error at syncing buffer with message %s in round %d: %s"
	bufferSyncedLogFmt  = "BMMC %s:%s synced buffer with message %s in round %d"

	invalidHostErr = "invalid host"

	gossipRoute          = "/gossip"
	solicitationRoute    = "/solicitation"
	synchronizationRoute = "/synchronization"
)

func addrPort(s string) (string, string, error) {
	host := strings.Split(s, ":")
	if len(host) != 2 {
		return "", "", errors.New(invalidHostErr)
	}

	addr := host[0]
	port := host[1]

	return addr, port, nil
}

func fullHost(addr, port string) string {
	return fmt.Sprintf("%s:%s", addr, port)
}

func (b *BMMC) gossipHandler(_ http.ResponseWriter, r *http.Request) {
	gossipDigest, tAddr, tPort, tRoundNumber, err := b.receiveGossip(r)
	if err != nil {
		b.config.Logger.Printf("%s", err)
		return
	}

	digest := b.messageBuffer.Digest()
	missingDigest := buffer.MissingStrings(gossipDigest, digest)

	hostAddr, hostPort, err := addrPort(r.Host)
	if err != nil {
		b.config.Logger.Printf(gossipHandlerErrLogFmt, err)
		return
	}

	if len(missingDigest) > 0 {
		solicitationMsg := HTTPSolicitation{
			Addr:        hostAddr,
			Port:        hostPort,
			RoundNumber: tRoundNumber,
			Digest:      missingDigest,
		}

		if err = b.sendSolicitation(solicitationMsg, tAddr, tPort); err != nil {
			b.config.Logger.Printf(gossipHandlerErrLogFmt, err)
			return
		}
	}
}

func (b *BMMC) solicitationHandler(_ http.ResponseWriter, r *http.Request) {
	missingDigest, tAddr, tPort, _, err := b.receiveSolicitation(r)
	if err != nil {
		b.config.Logger.Printf(solicitationHandlerErrLogFmt, err)
		return
	}

	missingElements := b.messageBuffer.ElementsFromIDs(missingDigest)

	hostAddr, hostPort, err := addrPort(r.Host)
	if err != nil {
		b.config.Logger.Printf(solicitationHandlerErrLogFmt, err)
		return
	}

	synchronizationMsg := HTTPSynchronization{
		Addr:     hostAddr,
		Port:     hostPort,
		Elements: missingElements,
	}

	if err = b.sendSynchronization(synchronizationMsg, tAddr, tPort); err != nil {
		b.config.Logger.Printf(solicitationHandlerErrLogFmt, err)
		return
	}
}

func (b *BMMC) synchronizationHandler(_ http.ResponseWriter, r *http.Request) {
	hostAddr, hostPort, err := addrPort(r.Host)
	if err != nil {
		b.config.Logger.Printf(synchronizationHandlerErrLogFmt, err)
		return
	}

	rcvElements, _, _, err := b.receiveSynchronization(r)
	if err != nil {
		b.config.Logger.Printf(synchronizationHandlerErrLogFmt, err)
		return
	}

	for _, m := range rcvElements {
		err = b.messageBuffer.Add(m)
		if err != nil {
			b.config.Logger.Printf(syncBufferLogErrFmt, hostAddr, hostPort, m.ID, b.gossipRound.GetNumber(), err)
		} else {
			b.config.Logger.Printf(bufferSyncedLogFmt, hostAddr, hostPort, m.ID, b.gossipRound.GetNumber())
			b.runCallbacks(m, hostAddr, hostPort)
		}
	}
}

func (b *BMMC) gracefullyShutdown() {
	if err := b.server.Shutdown(context.TODO()); err != nil {
		b.config.Logger.Printf(unableStopServerLogFmt, err)
	}
}

func (b *BMMC) newServer() *http.Server {
	return &http.Server{
		Addr: fullHost("0.0.0.0", b.config.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch path := r.URL.Path; path {
			case gossipRoute:
				b.gossipHandler(w, r)
			case solicitationRoute:
				b.solicitationHandler(w, r)
			case synchronizationRoute:
				b.synchronizationHandler(w, r)
			}
		}),
	}
}

func (b *BMMC) startServer(stop <-chan struct{}) error {
	errChan := make(chan error)

	go func() {
		b.config.Logger.Printf(startServerLogFmt, b.server.Addr)

		if err := b.server.ListenAndServe(); err != http.ErrServerClosed {
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
