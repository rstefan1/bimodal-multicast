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

// https://aods.cryingpotato.com/

import (
	"log"
	"log/slog"
	"os"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/rstefan1/bimodal-multicast/pkg/bmmc"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))

	host := Peer{
		Node: maelstrom.NewNode(),
	}

	cfg := &bmmc.Config{
		Host:          &host,
		Beta:          0.3,                    //nolint: gomnd
		RoundDuration: time.Millisecond * 500, //nolint: gomnd
		BufferSize:    8192,                   //nolint: gomnd
		Logger:        logger.With("component", "bmmc"),
	}

	bmmcNode, err := bmmc.New(cfg)
	if err != nil {
		log.Fatal(err)

		return
	}

	if err := bmmcNode.Start(); err != nil {
		log.Fatal(err)

		return
	}

	createAndRunServer(bmmcNode, host.Node, logger)
}
