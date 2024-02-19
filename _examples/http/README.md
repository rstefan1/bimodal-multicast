# Demo

Start a Docker container with BMMC instance:

```bash
make docker-net
make docker-build
make docker-start ADDR="172.24.0.4" PORT="19000"
```

# Scenario

## Start 7 nodes

```bash
make docker-start ADDR="172.24.101.1" PORT="19001"
	
make docker-start ADDR="172.24.102.2" PORT="19002"
	
make docker-start ADDR="172.24.103.3" PORT="19003"
	
make docker-start ADDR="172.24.104.4" PORT="19004"
	
make docker-start ADDR="172.24.105.5" PORT="19005"
	
make docker-start ADDR="172.24.106.6" PORT="19006"
	
make docker-start ADDR="172.24.107.7" PORT="19007"
```

## Add peers in buffer

First node:

```bash
add-peer 172.24.102.2 19002
add-peer 172.24.103.3 19003
add-peer 172.24.104.4 19004
add-peer 172.24.105.5 19005
add-peer 172.24.106.6 19006
add-peer 172.24.107.7 19007
```

Second node:

```bash
add-peer 172.24.101.1 19001
```

## Send messages

Add a message with `first-callback` type for callback to *third node*:

```bash
add-message aaaaa first-callback
```

#### Waiting to sync...

Add another message with `second-callback` type for callback to *fifth node*

```bash
add-message bbbbb second-callback
```

#### Waiting to sync...


## Start another two nodes

```bash
make docker-start ADDR="172.24.108.8" PORT="19008"
	
make docker-start ADDR="172.24.109.9" PORT="19009"
```

## Send messages between last two nodes added

Add peers to these nodes

Eighth node:

```bash
add-peer 172.24.109.9 19009
```

Ninth node:

```bash
add-peer 172.24.108.8 19008
```

Send a message

```bash
add-message ZZZZZ second-callback
```

#### Waiting to sync...

## Add the last two nodes turned on at the first nodes turned on

Third node:

```bash
add-peer 172.24.109.9 19009
```

#### Waiting to sync...

## Add another message

Seventh node:

```bash
add-message IIIII first-callback
```

#### Waiting to sync...
