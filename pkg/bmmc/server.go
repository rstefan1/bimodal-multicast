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
	gossipDigestBuffer, tAddr, tPort, tRoundNumber, err := receiveGossip(r)
	if err != nil {
		b.config.Logger.Printf("%s", err)
		return
	}

	msgDigestBuffer := b.messageBuffer.DigestBuffer()
	missingDigestBuffer := gossipDigestBuffer.GetMissingDigests(msgDigestBuffer)

	hostAddr, hostPort, err := addrPort(r.Host)
	if err != nil {
		b.config.Logger.Printf(gossipHandlerErrLogFmt, err)
		return
	}

	if missingDigestBuffer.Length() > 0 {
		solicitationMsg := HTTPSolicitation{
			Addr:        hostAddr,
			Port:        hostPort,
			RoundNumber: tRoundNumber,
			Digests:     *missingDigestBuffer,
		}

		err = sendSolicitation(solicitationMsg, tAddr, tPort)
		if err != nil {
			b.config.Logger.Printf(gossipHandlerErrLogFmt, err)
			return
		}
	}
}

func (b *BMMC) solicitationHandler(_ http.ResponseWriter, r *http.Request) {
	missingDigestBuffer, tAddr, tPort, _, err := receiveSolicitation(r)
	if err != nil {
		b.config.Logger.Printf(solicitationHandlerErrLogFmt, err)
		return
	}
	missingMsgBuffer := missingDigestBuffer.GetMissingMessageBuffer(b.messageBuffer)

	hostAddr, hostPort, err := addrPort(r.Host)
	if err != nil {
		b.config.Logger.Printf(solicitationHandlerErrLogFmt, err)
		return
	}

	synchronizationMsg := HTTPSynchronization{
		Addr:     hostAddr,
		Port:     hostPort,
		Messages: *missingMsgBuffer,
	}

	err = sendSynchronization(synchronizationMsg, tAddr, tPort)
	if err != nil {
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

	rcvMsgBuf, _, _, err := receiveSynchronization(r)
	if err != nil {
		b.config.Logger.Printf(synchronizationHandlerErrLogFmt, err)
		return
	}

	for _, m := range rcvMsgBuf.Messages {
		err = b.messageBuffer.AddMessage(m)
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
			case "/gossip":
				b.gossipHandler(w, r)
			case "/solicitation":
				b.solicitationHandler(w, r)
			case "/synchronization":
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

	// return <-errChan
	return nil
}
