package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"

	"github.com/rstefan1/bimodal-multicast/pkg/bmmc"
)

const (
	typeBodyKey    = "type"
	messageBodyKey = "message"
)

const cannotAddBMMCMessageErrFmt = "cannot add BMMC message %d: %w"

func CreateAndRunServer(b *bmmc.BMMC, n *maelstrom.Node, logger *log.Logger) {
	n.Handle(bmmc.GossipRoute, func(msg maelstrom.Message) error {
		var body map[string]string

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		b.GossipHandler([]byte(body[messageBodyKey]))

		return nil
	})

	n.Handle(bmmc.SolicitationRoute, func(msg maelstrom.Message) error {
		var body map[string]string

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		b.SolicitationHandler([]byte(body[messageBodyKey]))

		return nil
	})

	n.Handle(bmmc.SynchronizationRoute, func(msg maelstrom.Message) error {
		var body map[string]string

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		b.SynchronizationHandler([]byte(body[messageBodyKey]))

		return nil
	})

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		// Unmarshal the message body as a loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		if err := b.AddMessage(body[messageBodyKey], bmmc.NOCALLBACK); err != nil {
			logger.Printf(cannotAddBMMCMessageErrFmt, body[messageBodyKey], err)

		}

		return n.Reply(msg, map[string]any{
			typeBodyKey: "broadcast_ok",
		})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		// Unmarshal the message body as a loosely-typed map.
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body[typeBodyKey] = "read_ok"
		body["messages"] = b.GetMessages()

		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		topology := body["topology"].(map[string]interface{})
		peers := topology[n.ID()]

		for _, peer := range peers.([]interface{}) {
			if err := b.AddPeer(peer.(string)); err != nil {
				return err
			}
		}

		return n.Reply(msg, map[string]any{
			typeBodyKey: "topology_ok",
		})
	})

	if err := n.Run(); err != nil {
		logger.Fatal(err)
	}
}
