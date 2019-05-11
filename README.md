# Bimodal Multicast Protocol

[![Build Status](https://semaphoreci.com/api/v1/rstefan1-11/bimodal-multicast/branches/add-example-of-use/shields_badge.svg)](https://semaphoreci.com/rstefan1-11/bimodal-multicast)

This is an implementation of the Bimodal Multicast Protocol written in GO.

Currently can only sync string messages. This will be improved in the following
versions.

## Overview

The Bimodal Multicast Protocol runs in a series of rounds.
At the beginning of each round, every node randomly chooses another node and
sends it a digest of its message histories. The message is called gossip
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
    "github.com/rstefan1/bimodal-multicast/pkg/callback"
)
```

#### Register callbacks
``` golang
    cb := callback.NewRegistry()
    err := cb.Register(
        "awesome-callback",
        func (msg string) (bool, error) {
            fmt.Println("The message is:", msg)
            return true, nil
        })
```
Note! The buffer will be updated only if the callback fucntion call returns true.

#### Configure the protocol

```golang
    host := "localhost"
    port := "14999"

    cfg := bmmc.Config{
        Addr:   host,
        Port:   port,
        Peers: []peer.Peer{
            {
                Addr: host,
                Port: port,
            },
        },
        Callbacks: cb,
    }
```

#### Create an instance for protocol

```golang
    p, err := bmmc.New(cfg)
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
    p.AddMessage("awesome message", "awesome-callback")
```

#### Get all messages from the buffer

```golang
    messages := p.GetMessages()
```

## Performance

When beta is 0.5 and loss is 30% (after 20 retries):

![Alt text](docs/metrics/b_50_l_30.png)


## Roadmap to v0.2.x
 - [x] create an instance of Bimodal Multicast Protocol, start it,
 stop it, add message and retrieve all messages
 - [x] metrics
 - [x] register callbacks for each messages
 - [ ] add and remove peers via protocol
 - [ ] circular message buffer
 - [ ] more details about protocol in readme
 
## License

This project is licensed under Apache 2.0 license. Read the [LICENSE](LICENSE) file
in the top distribution directory for the full license text.
