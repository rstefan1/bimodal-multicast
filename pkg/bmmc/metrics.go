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
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/testutil"
)

func RunWithSpec(retries int,
	noPeers int,
	loss float64,
	beta float64,
	cbRegistry map[string]func(string) (bool, error),
	cbType string,
	timeout time.Duration) error {

	// create path for logs
	logPath := "metrics/logs"
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if err := os.Mkdir(logPath, 0777); err != nil {
			return nil
		}
	}

	// run the protocol `retries` times
	for r := 0; r < retries; r++ {
		peers := []Peer{}
		nodes := []*Bmmc{}

		// create file for logs
		fName := fmt.Sprintf("%s/log_r_%d_l_%d_b_%d.log", logPath, r, int(loss*100), int(beta*100))
		f, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return err
		}

		// create logger
		logger := log.New(f, "", 0)

		// set random ports for peers
		for i := 0; i < noPeers; i++ {
			suggestedPort := testutil.SuggestPort()
			peers = append(peers, Peer{Addr: "localhost", Port: suggestedPort})
		}

		// create nodes
		for i := 0; i < noPeers; i++ {
			n, err := New(
				&Config{
					Addr:      peers[i].Addr,
					Port:      peers[i].Port,
					Peers:     peers,
					Beta:      beta,
					Logger:    logger,
					Callbacks: cbRegistry,
				},
			)
			if err != nil {
				return err
			}

			nodes = append(nodes, n)
		}

		// add a message
		msg := "another-awesome-message"
		randomNode := rand.Intn(noPeers)
		nodes[randomNode].AddMessage(msg, cbType)

		// start nodes
		for i := 0; i < noPeers; i++ {
			if err := nodes[i].Start(); err != nil {
				return nil
			}
		}

		time.Sleep(timeout)

		// stop nodes
		for i := 0; i < noPeers; i++ {
			nodes[i].Stop()
		}

		// wait for stopping all nodes
		time.Sleep(time.Millisecond * 100)

		// close log file
		f.Close()
	}

	return nil
}
