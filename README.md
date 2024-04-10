# Bimodal Multicast Protocol

![Build](https://github.com/rstefan1/bimodal-multicast/workflows/Build/badge.svg?branch=master)
[![GoDoc](https://godoc.org/github.com/rstefan1/bimodal-multicast?status.svg)](https://godoc.org/github.com/rstefan1/bimodal-multicast)

This is an implementation of the Bimodal Multicast Protocol written in GO.

You can synchronize all types of messages: bool, string, int, 
complex structs, etc.

---

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

---

## Usage

- ### Step 1: Imports

```go
import "github.com/rstefan1/bimodal-multicast/pkg/bmmc"
```

- ### Step 2: Configure the host

The host must implement [Peer interface](https://github.com/rstefan1/bimodal-multicast/blob/f98c69dbc8ac22decdb438a1d6b5abc4b5db2db0/pkg/internal/peer/peer.go#L20):

```go
type Peer interface {
	String() string
	Send(msg []byte, route string, peerToSend string) error
}
```

- ### Step 3: Configure the bimodal-multicast protocol

```go
cfg := bmmc.Config{
    Host:           host,
    Callbacks:      map[string]func (interface{}, *log.Logger) error {
        "custom-callback":
        func (msg interface{}, logger *log.Logger) error {
            fmt.Println("The message is:", msg)

            return nil
        },
    },
    Beta:           float64,
    Logger:         logger,
    RoundDuration:  time.Second * 5,
    BufferSize:     2048,
}
```

| Config        | Required | Description                                                                                                                                                                                                                 |
|---------------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Host          | Yes      | Host of Bimodal Multicast server. <br/>Must implement [Peer interface](https://github.com/rstefan1/bimodal-multicast/blob/f98c69dbc8ac22decdb438a1d6b5abc4b5db2db0/pkg/internal/peer/peer.go#L20). Check the previous step. |
| Callback      | No       | You can define a list of callbacks.<br/>A callback is a function that is called every time a message on the server is synchronized.                                                                                         |
| Beta          | No       | The beta factor is used to control the ratio of unicast to multicast traffic that the protocol allows.                                                                                                                      |
| Logger        | No       | You can define a [structured logger](https://pkg.go.dev/log/slog).                                                                                                                                                          | 
| RoundDuration | No       | The duration of a gossip round.                                                                                                                                                                                             | 
| BufferSize    | Yes      | The size of messages buffer.<br/>The buffer will also include internal messages (e.g. synchronization of the peer list).<br/>***When the buffer is full, the oldest message will be removed.***                             |


- ### Step 4. Create a bimodal multicast server

```go
bmmcServer, err := bmmc.New(cfg)
```

- ### Step 5. Create the host server (e.g. a HTTP server)

The server must handle a list of predefined requests.
Each of these handlers must read the message body and call a predefined function.

| Handler route     | Function to be called                     |
|-------------------|-------------------------------------------|
| `bmmc.GossipRoute` | `bmmcServer.GossipHandler(body)`          |
| `bmmc.SolicitationRoute` | `bmmcServer.SolicitationHandler(body)`    |
| `bmmc.SynchronizationRoute` | `bmmcServer.SynchronizationHandler(body)` |

For more details, check the [exemples](#examples).

- ### Step 6. Start the host server and the bimodal multicast server

```go
# Start the host server
hostServer.Start()

# Start the bimodal multicast server
bmmcServer.Start()
```

<a name="custom_anchor_name"></a>
- ### Step 7. Add a message to broadcast

```go
bmmcServer.AddMessage("new-message", "my-callback")
bmmcServer.AddMessage(12345, "another-callback")
bmmcServer.AddMessage(true, bmmc.NOCALLBACK)
```

- ### Step 8. Retrieve all messages from buffer

```go
bmcServer.GetMessages()
```

- ### Step 9. Add/Remove peers

```go
bmmcServer.AddPeer(peerToAdd)
bmmcServer.RemovePeer(peerToRemove)
```

- ### Step 10. Stop the bimodal multicast server

```go
bmmcServer.Stop()
```

---

## Examples

<a name="examples"></a>

1. using a [http server](_examples/http)
2. using a [maelstrom server](_examples/maelstrom)

---

## Contributing

I welcome all contributions in the form of new issues for feature requests, bugs
or even pull requests.

---

## License

This project is licensed under Apache 2.0 license. Read the [LICENSE](LICENSE) file
in the top distribution directory for the full license text.
