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
)

func addrPort(s string) (string, string) {
	host := strings.Split(s, ":")
	if len(s) < 2 {
		panic("invalid host")
	}

	addr := host[0]
	port := host[1]
	return addr, port
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

	hostAddr, hostPort := addrPort(r.Host)

	if missingDigestBuffer.Length() > 0 {
		solicitationMsg := HTTPSolicitation{
			Addr:        hostAddr,
			Port:        hostPort,
			RoundNumber: tRoundNumber,
			Digests:     *missingDigestBuffer,
		}

		err = sendSolicitation(solicitationMsg, tAddr, tPort)
		if err != nil {
			b.config.Logger.Printf("%s", err)
			return
		}
	}
}

func (b *BMMC) solicitationHandler(_ http.ResponseWriter, r *http.Request) {
	missingDigestBuffer, tAddr, tPort, _, err := receiveSolicitation(r)
	if err != nil {
		b.config.Logger.Printf("%s", err)
		return
	}
	missingMsgBuffer := missingDigestBuffer.GetMissingMessageBuffer(b.messageBuffer)

	hostAddr, hostPort := addrPort(r.Host)

	synchronizationMsg := HTTPSynchronization{
		Addr:     hostAddr,
		Port:     hostPort,
		Messages: *missingMsgBuffer,
	}

	err = sendSynchronization(synchronizationMsg, tAddr, tPort)
	if err != nil {
		b.config.Logger.Printf("%s", err)
		return
	}
}

func (b *BMMC) synchronizationHandler(_ http.ResponseWriter, r *http.Request) {

	hostAddr, hostPort := addrPort(r.Host)

	rcvMsgBuf, _, _, err := receiveSynchronization(r)
	if err != nil {
		b.config.Logger.Printf("BMMC %s:%s Error at receiving synchronization error: %s", hostAddr, hostPort, err)
		return
	}

	for _, m := range rcvMsgBuf.Messages {
		err = b.messageBuffer.AddMessage(m)
		if err != nil {
			b.config.Logger.Printf("BMMC %s:%s error at syncing buffer with message %s in round %d: %s", hostAddr, hostPort, m.ID, b.gossipRound.GetNumber(), err)
		} else {
			b.config.Logger.Printf("BMMC %s:%s synced buffer with message %s in round %d", hostAddr, hostPort, m.ID, b.gossipRound.GetNumber())

			// run callback function for messages with a callback registered
			if m.CallbackType != callback.NOCALLBACK {
				err = b.defaultCallbacks.RunDefaultCallbacks(m, b.peerBuffer, b.config.Logger)
				if err != nil {
					b.config.Logger.Printf("Error at calling default callback at %s:%s for message %s in round %d", hostAddr, hostPort, m.ID, b.gossipRound.GetNumber())
				}

				err = b.customCallbacks.RunCustomCallbacks(m, b.config.Logger)
				if err != nil {
					b.config.Logger.Printf("Error at calling custom callback at %s:%s for message %s in round %d", hostAddr, hostPort, m.ID, b.gossipRound.GetNumber())
				}
			}
		}
	}
}

func (b *BMMC) gracefullyShutdown() {
	if err := b.server.Shutdown(context.TODO()); err != nil {
		b.config.Logger.Printf("Unable to shutdown server properly: %s", err)
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
		b.config.Logger.Printf("Server listening at %s", b.server.Addr)

		if err := b.server.ListenAndServe(); err != http.ErrServerClosed {
			b.config.Logger.Printf("Unable to start  server: %s", err)
			errChan <- err
			return
		}
		errChan <- nil
	}()

	go func() {
		<-stop
		b.gracefullyShutdown()
	}()

	// return <-errChan
	return nil
}
