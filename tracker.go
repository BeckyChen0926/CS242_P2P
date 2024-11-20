package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
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

	newPeer := Peer{peerCon, peerIP, neighbors, chunks}

	// append the IPaddr of the peer into peerList
	peerList = append(peerList, newPeer)
	// assign up to 10 new neighbors
	giveNewNeighbors(newPeer)

}

// assign up to 10 neighbors to p
func giveNewNeighbors(p Peer) {
	if len(peerList) > 10 {

		// TODO: how to get 10 unique neighbors w/o duplication
		// i think this will do?

		var neighborCount int

		// loop until have 10 new neighbors
		for neighborCount < 10 {
			// Generate a random integer between 0 and n
			randomInt := rand.Intn(len(peerList))
			randNeighbor := peerList[randomInt]

			for j := 0; j < len(p.neighbors); j++ {

				// if duplicate, skip and do the next loop
				if randNeighbor.IPAddr.Equal(p.neighbors[j].IPAddr) || randNeighbor.IPAddr.Equal((p.IPAddr)) {
					continue
				}
			}

			p.neighbors = append(p.neighbors, randNeighbor)
			neighborCount++

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
		// i think we can also just check if their ipaddr are equal?
		if !peerList[i].IPAddr.Equal(p.IPAddr) {
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
