package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
)

// func handleConnection(conn net.Conn) {
// 	defer conn.Close()
// 	buf := make([]byte, 1024)
// 	for {
// 		n, err := conn.Read(buf)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		fmt.Printf("Received: %s", string(buf[:n]))
// 		data := strings.ToUpper(string(buf[:n]))
// 		_, err = conn.Write([]byte(data))
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 	}
// }

// initiate new list to store all peers in the network
var peerList []Peer

func registerPeer(peerCon net.Conn) {

	// of type IP
	peerIP := net.ParseIP(peerCon.RemoteAddr().String())

	var neighbors []Peer

	// append the IPaddr of the peer into peerList
	peerList = append(peerList, Peer{peerCon, peerIP, neighbors})

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

			// for j := 0; j < len(p.neighbors); j++ {
			// 	if reflect.DeepEqual(randNeighbor, p.neighbors[j]) {
			// 		randomInt = rand.Intn(len(peerList))
			// 		randNeighbor = peerList[randomInt]
			// 	}
			// }

			// things := []string{"foo", "bar", "baz"}
			// slices.Contains(p.neighbors., randNeighbor) // true

			// for strings.Contains(p.neighbors, randNeighbor) {
			// 	randomInt := rand.Intn(len(peerList))
			// 	randNeighbor := peerList[randomInt]

			// }
			// p.neighbors = append(p.neighbors, peerList[randomInt])

		}

	} else {
		p.neighbors = peerList
	}

}

func removePeer(p Peer) {
	//Removes desired peer from the peerlist

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
