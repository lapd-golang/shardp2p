package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"

	golog "github.com/ipfs/go-log"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	gologging "github.com/whyrusleeping/go-logging"
)

func main() {
	// LibP2P code uses golog to log messages. They log with different
	// string IDs (i.e. "swarm"). We can control the verbosity level for
	// all loggers with:
	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info

	// Parse options from the command line
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	shardID := flag.Int("shardID", 1, "shard id to connect to")
	flag.Parse()

	bootstrapNodes := make(chan map[int][]host.Host)

	// Initializes bootstrap nodes and makes them begin listening, outputs their
	// full addresses in pretty print to logs
	go func() {
		log.Println("Starting...")
		shardCount := 5
		b := MakeBootstrapManager(*seed, shardCount)
		bootstrapNodes <- b.InitializeNodes()
		log.Println("Initialized and Running Bootstrap Nodes")
	}()

	// Adds the bootstrap nodes to this node's peers so libp2p knows how to find them
	go func() {
		clientnode, err := makeRandomHost(*seed)
		if err != nil {
			panic(err)
		}
		initializedNodes := <-bootstrapNodes
		log.Printf("Started Sharding Client With Node Addr: %s", clientnode.Addrs()[0])
		log.Println("Opening Connection to Bootstrap Nodes...")

		for _, node := range initializedNodes[*shardID] {
			clientnode.Peerstore().AddAddr(node.ID(), node.Addrs()[0], pstore.PermanentAddrTTL)
			// This will be handled by bootstrap nodes that speak the same subprotocol, in this case
			// /shardprotocol/shardID
			s, err := clientnode.NewStream(context.Background(), node.ID(), "/shardprotocol/1")
			if err != nil {
				log.Fatalln(err)
			}

			_, err = s.Write([]byte("Hello, Shard!"))
			if err != nil {
				log.Fatalln(err)
			}

			out, err := ioutil.ReadAll(s)
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("Reply from Bootstrap Node: %q", out)
		}
	}()

	select {} // hang forever...
}

func echo(s net.Stream) error {
	log.Println("Received Message from Sharding Client: Hello, Shard!")
	s.Write([]byte("Hello, Shard!"))
	return nil
}
