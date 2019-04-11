# Bimodal Multicast Protocol

[![Build Status](https://semaphoreci.com/api/v1/projects/42333e66-e66b-4bdf-bbd6-29e8deae4ebf/2519090/badge.svg)](https://semaphoreci.com/rstefan1-11/bimodal-multicast)

This is an implementation of the Bimodal Multicast Protocol written in GO.

Currently can only sync string messages. This will be improved in the following
versions.

## Overview

The protocol runs in a series of rounds.
At the beginning of each round, every node randomly chooses another node and
sends it a digest of its message histories. The message is the gossip
message.
The node that receive the gossip message compares the given digest with the
messages in its own message buffer.
If the digest differs from its message histories, then it send a message
back to the original sender to request the missing messages. This message is
called solicitation.

## Usage

#### Imports

```golang
import (
    "github.com/rstefan1/bimodal-multicast/pkg/peer"
    "github.com/rstefan1/bimodal-multicast/pkg/bmmc"
)
```

#### Configure the protocol

```golang
    host := "localhost"
    port := "14999"

    cfg := bmmc.Config{
        HTTPAddr:   host,
        HTTPPort:   port,
        GossipAddr: "localhost",
        GossipPort: "15999",
        Peers: []peer.Peer{
            {
                Addr: host,
                Port: port,
            },
        },
    }
```

#### Create an instance for protocol

```golang
    p := bmmc.New(cfg)
```

#### Start the protocol

```golang
    p.Start()
```

#### Stop the protocol

```golang
    p.Stop()
```

#### Add a new string message in buffer

```golang
    p.AddMessage("awesome message")
```

#### Get all messages from the buffer

```golang
    messages := p.GetMessages()
```

