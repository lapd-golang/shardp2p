package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
)

// BootstrapManager handles all logic for bootstrap nodes
type BootstrapManager struct {
	// map of shard_id's to libp2p nodes
	bootstrapNodes map[int][]host.Host
	seed           int64
	shardCount     int
}

// MakeBootstrapManager creates a new struct instance
func MakeBootstrapManager(randseed int64, shardCount int) *BootstrapManager {
	return &BootstrapManager{
		bootstrapNodes: make(map[int][]host.Host),
		seed:           randseed,
		shardCount:     shardCount,
	}
}

// InitializeNodes sets up 3 boostrap nodes for each shard
func (b *BootstrapManager) InitializeNodes() map[int][]host.Host {
	for i := 0; i < b.shardCount; i++ {
		var hosts []host.Host
		for len(hosts) < 3 {
			h, err := makeRandomHost(b.seed)
			if err != nil {
				panic(err)
			}

			// THIS is what determines what a bootstrap does upon receiving an incoming
			// connection from a sharding client. In this case, we just echo the message
			// sent to us, but in a real system, we can do peer discovery
			h.SetStreamHandler("/shardprotocol/1", func(s net.Stream) {
				log.Println("Bootstrap Node Got a New Incoming Connection!")
				if err := echo(s); err != nil {
					log.Println(err)
					s.Reset()
				} else {
					s.Close()
				}
			})
			hosts = append(hosts, h)
		}
		b.bootstrapNodes[i] = hosts
		log.Printf("Initialized Bootstrap Nodes for Shard %d\n", i)
	}
	return b.bootstrapNodes
}

func makeRandomHost(randseed int64) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	// Generate a key pair for this host. We will use it at least
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	// random port in range
	port := random(3000, 3500)

	// By default uses secure I/O
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)),
		libp2p.Identity(priv),
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	return basicHost, nil
}

func random(min, max int) int {
	mrand.Seed(time.Now().Unix())
	return mrand.Intn(max-min) + min
}
