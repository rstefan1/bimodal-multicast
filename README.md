# Bimodal Multicast Protocol

[![Build Status](https://semaphoreci.com/api/v1/rstefan1-11/bimodal-multicast/branches/master/shields_badge.svg)](https://semaphoreci.com/rstefan1-11/bimodal-multicast)

This is an implementation of the Bimodal Multicast Protocol written in GO.

You can synchronize all types of messages: bool, string, int, 
complex structs, etc.

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

* Imports

```golang
import (
    "github.com/rstefan1/bimodal-multicast/pkg/bmmc"
)
```

* Configure the protocol

```golang
    cfg := bmmc.Config{
        Addr:      "localhost",
        Port:      "14999",
        Callbacks: map[string]func (interface{}, *log.Logger) (error) {
            "awesome-callback":
            func (msg interface{}, logger *log.Logger) (error) {
                fmt.Println("The message is:", msg)
                return nil
            },
        },
    }
```

* Create an instance for protocol

```golang
    p, err := bmmc.New(cfg)
```

* Start the protocol

```golang
    p.Start()
```

* Stop the protocol

```golang
    p.Stop()
```

* Add a new message in buffer

```golang
    err := p.AddMessage("awesome message", "awesome-callback")
    
    err := p.AddMessage(12345, "awesome-callback")
```

For messages without callback, you can use `bmmc.NOCALLBACK` as callback type.

* Get all messages from the buffer

```golang
    messages := p.GetMessages()
```

* Add a new peer in peers buffer

```golang
    err := p.AddPeer("localhost", "18999")
```

* Remove a peer from peers buffer

```golang
    err := p.RemovePeer("localhost", "18999")
```


## Performance

When beta is 0.5 and loss is 30% (after 20 retries):

![Alt text](docs/metrics/b_50_l_30.png)


## Roadmap to v0.2
 - [x] create an instance of Bimodal Multicast Protocol, start it,
 stop it, add message and retrieve all messages
 - [x] metrics
 - [x] register callbacks for each messages
 - [x] add and remove peers via protocol
 - [x] synchronize all types of messages (bool, string, int, complex structs, etc.)
 
## Contributing

We welcome all contributions in the form of new issues for feature requests, bugs
or even pull requests.

## License

This project is licensed under Apache 2.0 license. Read the [LICENSE](LICENSE) file
in the top distribution directory for the full license text.
