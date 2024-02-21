package main

import (
	"io/ioutil"
	"log"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/rstefan1/bimodal-multicast/pkg/bmmc"
)

func main() {
	// https://aods.cryingpotato.com/

	logger := log.Default()
	logger.SetOutput(ioutil.Discard)

	host := Peer{
		Node: maelstrom.NewNode(),
	}

	cfg := &bmmc.Config{
		Host:          &host,
		Beta:          0.75,
		RoundDuration: time.Millisecond * 200,
		BufferSize:    8192,
		Logger:        logger,
	}

	bmmcNode, err := bmmc.New(cfg)
	if err != nil {
		logger.Fatal(err)

		return
	}

	if err := bmmcNode.Start(); err != nil {
		logger.Fatal(err)

		return
	}

	CreateAndRunServer(bmmcNode, host.Node, logger)
}
