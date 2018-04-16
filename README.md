# ShardP2P Proof of Concept

Playing around with libp2p for a sharded network. We can instantiate a set of bootnodes for each shard and establish an initial handshake from a sharding client. Then, we create a simple echo implementation using libp2p.

This is a playground repo, allowing us to see how we can move into go-libp2p for sharding with geth.

See: [Prysmatic Labs P2P Specs Overview](https://docs.google.com/document/d/1K9NVV2SBjxdgejWnip3l-ZYyknOdWu6i5Ot_X_y6t1k/edit#)

## Running

```
$ git clone https://github.com/rauljordan/shardp2p
$ go build
$ ./shardp2p -shardID 1
```
