# Bimodal Multicast Protocol

[![Build Status](https://semaphoreci.com/api/v1/projects/42333e66-e66b-4bdf-bbd6-29e8deae4ebf/2555382/shields_badge.svg)](https://semaphoreci.com/rstefan1-11/bimodal-multicast)

At the beginning of each round, every node `randomly` chooses another node and
sends it a `digest` of its `message histories`. The messsage is the `gossip 
message`.

The node that receive the `gossip message` compares the given `digest` with the
messages in its own `message buffer`.

If the `digest` differs from its `message histories`, then it send a message
back to the original sender to request the missing messages. The message is 
called `solicitation`.

### Sequence diagram

![alt text](docs/diagrams/sequence-diagram.png)

### Round goroutine vs. HTTP Server goroutine

![alt text](docs/diagrams/threads-diagram.png)

### Receiving a gossip message

![alt text](docs/diagrams/threads-diagram-2.png)

