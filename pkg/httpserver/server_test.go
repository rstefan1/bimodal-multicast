package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/buffer"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/config"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/httpmessage"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

type receivedMessages struct {
	recv interface{}
	mux  sync.Mutex
}

func (r *receivedMessages) PutMessage(m interface{}) {
	r.mux.Lock()
	r.recv = m
	r.mux.Unlock()
}

func (r *receivedMessages) GetMessage() interface{} {
	r.mux.Lock()
	m := r.recv
	r.mux.Unlock()
	return m
}

func getPort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

var _ = Describe("HTTP Server", func() {
	var (
		httpServerPort      string
		mockServerPort      string
		httpServerCfg       config.Config
		peerBuffer          []peer.Peer
		httpServerMsgBuffer buffer.MessageBuffer
		httpServerStop      chan struct{}
		rcvMsg              receivedMessages
		mockServer          *http.Server
		requestPath         string
	)

	BeforeEach(func() {
		httpServerPort = strconv.Itoa(getPort())
		mockServerPort = strconv.Itoa(getPort())

		httpServerCfg = config.Config{
			HTTPAddr: fmt.Sprintf(":%s", httpServerPort),
		}

		peerBuffer = []peer.Peer{}

		httpServerMsgBuffer = buffer.NewMessageBuffer()

		httpServerStop = make(chan struct{})

		rcvMsg = receivedMessages{
			mux: sync.Mutex{},
		}

		// start http server
		err := Start(&peerBuffer, &httpServerMsgBuffer, httpServerStop, httpServerCfg)
		Expect(err).To(Succeed())
	})

	AfterEach(func() {
		close(httpServerStop)
	})

	Describe("at gossip request", func() {
		BeforeEach(func() {
			// create new mock server
			mockServer = &http.Server{
				Addr: fmt.Sprintf(":%s", mockServerPort),
				Handler: http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						if r.URL.Path == "/solicitation" {
							decoder := json.NewDecoder(r.Body)
							var t httpmessage.HTTPSolicitation
							err := decoder.Decode(&t)
							if err != nil {
								panic(err)
							}
							rcvMsg.PutMessage(t)
						}
					}),
			}

			// start mock server
			go func() {
				err := startHTTPServer(mockServer)
				Expect(err).To(Succeed())
			}()

			requestPath = fmt.Sprintf("http://localhost:%s/gossip", httpServerPort)
		})

		AfterEach(func() {
			gracefullShutdown(mockServer)
		})

		It("responds with solicitation request when nodes have different digests", func() {
			// create gossip request
			gossipDigest := buffer.DigestBuffer{
				Digests: []buffer.Digest{
					{ID: fmt.Sprintf("%d", rand.Int31())},
				},
			}
			gossipMessage := httpmessage.HTTPGossip{
				Addr:        "localhost",
				Port:        mockServerPort,
				RoundNumber: rand.Int63(),
				Digests:     gossipDigest,
			}
			jsonGossip, err := json.Marshal(gossipMessage)
			Expect(err).To(Succeed())

			// send gossip request
			_, err = http.Post(requestPath, "json", bytes.NewBuffer(jsonGossip))
			Expect(err).To(Succeed())

			// wait for receiving response from http server
			msg := rcvMsg.GetMessage()
			Expect(msg.(httpmessage.HTTPSolicitation).Digests).To(Equal(gossipDigest))
		})

		It("does not respond with solicitation request when nodes have same digests", func() {
			httpServerMsgBuffer.AddMessage(buffer.Message{
				ID:          fmt.Sprintf("%d", rand.Int31()),
				Msg:         fmt.Sprintf("%d", rand.Int31()),
				GossipCount: rand.Int(),
			})
			gossipMessage := httpmessage.HTTPGossip{
				Addr:        "localhost",
				Port:        mockServerPort,
				RoundNumber: rand.Int63(),
				Digests:     httpServerMsgBuffer.DigestBuffer(),
			}
			jsonDigest, err := json.Marshal(gossipMessage)
			Expect(err).To(Succeed())

			// send gossip request
			_, err = http.Post(requestPath, "json", bytes.NewBuffer(jsonDigest))
			Expect(err).To(Succeed())

			// wait 1 second for solicitation message
			time.Sleep(time.Second * 1)
			msg := rcvMsg.GetMessage()
			Expect(msg).Should(BeNil())
		})
	})

	Describe("at solicitation request", func() {
		BeforeEach(func() {
			// create new mock server
			mockServer = &http.Server{
				Addr: fmt.Sprintf(":%s", mockServerPort),
				Handler: http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						if r.URL.Path == "/synchronization" {
							decoder := json.NewDecoder(r.Body)
							var t httpmessage.HTTPSynchronization
							err := decoder.Decode(&t)
							if err != nil {
								panic(err)
							}
							rcvMsg.PutMessage(t)
						}
					}),
			}

			// start mock server
			go func() {
				err := startHTTPServer(mockServer)
				Expect(err).To(Succeed())
			}()

			requestPath = fmt.Sprintf("http://localhost:%s/solicitation", httpServerPort)
		})

		AfterEach(func() {
			gracefullShutdown(mockServer)
		})

		It("responds with synchronization message", func() {
			messageID := fmt.Sprintf("%d", rand.Int31())

			// populate buffer with a message
			httpServerMsgBuffer.AddMessage(buffer.Message{
				ID:          messageID,
				Msg:         fmt.Sprintf("%d", rand.Int31()),
				GossipCount: rand.Int(),
			})

			// create solicitation request
			solicitationDigest := buffer.DigestBuffer{
				Digests: []buffer.Digest{
					{ID: messageID},
				},
			}
			solicitationMessage := httpmessage.HTTPSolicitation{
				Addr:        "localhost",
				Port:        mockServerPort,
				RoundNumber: rand.Int63(),
				Digests:     solicitationDigest,
			}
			jsonSolicitation, err := json.Marshal(solicitationMessage)
			Expect(err).To(Succeed())

			// send solicitation request
			_, err = http.Post(requestPath, "json", bytes.NewBuffer(jsonSolicitation))
			Expect(err).To(Succeed())

			// wait for receiving response from http server
			msg := rcvMsg.GetMessage()
			Expect(msg.(httpmessage.HTTPSynchronization).Messages).To(Equal(httpServerMsgBuffer))
		})
	})

	Describe("at synchronization request", func() {
		BeforeEach(func() {
			requestPath = fmt.Sprintf("http://localhost:%s/synchronization", httpServerPort)
		})
		It("updates the message buffer", func() {
			syncMsgBuffer := buffer.NewMessageBuffer()
			syncMsgBuffer.AddMessage(buffer.Message{
				ID:          fmt.Sprintf("%d", rand.Int31()),
				Msg:         fmt.Sprintf("%d", rand.Int31()),
				GossipCount: rand.Int(),
			})

			// create synchronization request
			syncMessage := httpmessage.HTTPSynchronization{
				Addr:     "localhost",
				Port:     mockServerPort,
				Messages: syncMsgBuffer,
			}
			jsonSync, err := json.Marshal(syncMessage)
			Expect(err).To(Succeed())

			// send synchronization request to http server
			_, err = http.Post(requestPath, "json", bytes.NewBuffer(jsonSync))
			Expect(err).To(Succeed())

			// wait for synchronization of shared message buffer
			Expect(httpServerMsgBuffer.Messages).To(Equal(syncMsgBuffer.Messages))
		})
	})
})
