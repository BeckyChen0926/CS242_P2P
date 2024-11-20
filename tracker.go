package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"reflect"
)

type Peer struct {
	conn      net.Conn
	IPAddr    net.IP
	neighbors []Peer
	chunks    []string //chunks that the peer now has. TODO: find write type
}

// initiate new list to store all peers in the network
var peerList []Peer

func registerPeer(peerCon net.Conn) {

	// of type IP
	peerIP := net.ParseIP(peerCon.RemoteAddr().String())

	var neighbors []Peer
	var chunks []string

	// append the IPaddr of the peer into peerList
	peerList = append(peerList, Peer{peerCon, peerIP, neighbors, chunks})

}

// assign up to 10 neighbors to p
func giveNewNeighbors(p Peer) {
	if len(peerList) > 10 {

		// TODO: how to get 10 unique neighbors w/o duplication

		var randNeighbors []Peer
		for i := 0; i < 10; i++ {
			// Generate a random integer between 0 and n
			randomInt := rand.Intn(len(peerList))
			randNeighbor := peerList[randomInt]

			randNeighbors = append(randNeighbors, randNeighbor)

			for j := 0; j < len(p.neighbors); j++ {
				if reflect.DeepEqual(randNeighbor, p.neighbors[j]) {
					randomInt = rand.Intn(len(peerList))
					randNeighbor = peerList[randomInt]
				}
			}

			// things := []string{"foo", "bar", "baz"}
			// slices.Contains(p.neighbors., randNeighbor) // true

			// for strings.Contains(p.neighbors, randNeighbor) {
			// 	randomInt := rand.Intn(len(peerList))
			// 	randNeighbor := peerList[randomInt]

			// }

			// add the 10 new neighbors to p's neighborList
			for i := 0; i < 10; i++ {
				p.neighbors = append(p.neighbors, randNeighbors[i])

			}
		}

	} else {
		p.neighbors = peerList
	}

}

// remove peer p who is leaving
func removePeer(p Peer) {
	//Removes desired peer from the peerlist
	var newPeerList []Peer
	for i := 0; i < len(peerList); i++ {
		if !reflect.DeepEqual(p, peerList[i]) {
			newPeerList = append(newPeerList, peerList[i])
		}
	}
	peerList = newPeerList
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	fmt.Println("Listening on port 8000")
	for {
		peer, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go registerPeer(peer)
	}
}
