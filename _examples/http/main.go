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
	"bufio"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rstefan1/bimodal-multicast/pkg/bmmc"
)

func main() { //nolint: funlen, cyclop, gocyclo, gocognit
	addr := os.Getenv("ADDR")
	port := os.Getenv("PORT")

	fmt.Println()
	fmt.Println("*** Address:", addr)
	fmt.Println("*** Port:", port)

	// // create file for logs
	// fName := fmt.Sprintf("log-%s-%s-%d", addr, port, rand.Int31())
	// logFile, err := os.OpenFile(fName, os.O_RDWR|os.O_CREATE, 0600)
	// if err != nil {
	// 	fmt.Println("Error at creating log file with name", fName)
	// 	return
	// }

	callbacks := map[string]func(any, *slog.Logger) error{
		"first-callback": func(msg any, _ *slog.Logger) error {
			fmt.Printf("*** First callback called for message `%s`. ***\n", msg)

			return nil
		},
		"second-callback": func(msg any, _ *slog.Logger) error {
			fmt.Printf("### Second callback called for message `%s`. ###\n", msg)

			return nil
		},
	}

	host, err := NewPeer(addr, port, &http.Client{})
	if err != nil {
		fmt.Println("Error at creating BMMC Host", err)

		return
	}

	cfg := &bmmc.Config{
		Host:       host,
		Callbacks:  callbacks,
		Beta:       0.25, //nolint: gomnd
		BufferSize: 1024, //nolint: gomnd
		// Logger:     slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
		// RoundDuration:  time.Second * 2,
	}

	node, err := bmmc.New(cfg)
	if err != nil {
		fmt.Println("Error at creating BMMC instance", err)

		return
	}

	err = node.Start()
	if err != nil {
		fmt.Println("Error at starting BMMC instance", err)

		return
	}

	time.Sleep(time.Millisecond * 150) //nolint: gomnd

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error at reading the line")

			continue
		}

		args := strings.Fields(line)
		if len(args) < 1 {
			// empty line
			continue
		}

		switch args[0] {
		case "add-peer":
			if len(args) != 3 { //nolint: gomnd
				fmt.Println("Invalid command. The `add-peer` command must be in form: " +
					"add-peer 127.100.1.4 19999")

				break
			}

			addr := args[1]
			port := args[2]

			peer, err := NewPeer(addr, port, nil)
			if err != nil {
				fmt.Println("Error at creating peer:", err)
			}

			err = node.AddPeer(peer.String())
			if err != nil {
				fmt.Println("Error at adding peer in buffer:", err)
			}

		case "delete-peer":
			if len(args) != 3 { //nolint: gomnd
				fmt.Println("Invalid command. The `delete-peer` command must be in form: " +
					"delete-peer 127.100.1.4 19999")

				break
			}

			addr := args[1]
			port := args[2]

			peer, err := NewPeer(addr, port, nil)
			if err != nil {
				fmt.Println("Error at creating peer:", err)
			}

			err = node.RemovePeer(peer.String())
			if err != nil {
				fmt.Println("Error at removing peer from buffer:", err)
			}

		case "add-message":
			if len(args) != 3 { //nolint: gomnd
				fmt.Println("Invalid command. The `add-message` command must be in form: " +
					"add-message my-message first-callback")

				break
			}

			message := args[1]
			cbType := args[2]

			err := node.AddMessage(message, cbType)
			if err != nil {
				fmt.Println("Error at adding message in buffer:", err)
			}

		case "get-messages":
			fmt.Println("Messages:\n", node.GetMessages())

		case "get-peers":
			fmt.Println("Peers:\n", node.GetPeers())

		case "stop":
			node.Stop()

			return

		case "exit":
			node.Stop()

			return

		default:
			fmt.Println("Wrong command. Available commands are:")
			fmt.Println("   -> add-peer")
			fmt.Println("   -> delete-peer")
			fmt.Println("   -> add-message")
			fmt.Println("   -> stop / exit")

			continue
		}
	}
}
