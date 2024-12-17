package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
)

/*
* Our peer struct (for a tracker) which has the following information:
* port - port number
* IPAddr - IP address
* neighbors - list of the peer's direct neighbors

* NOTE: we don't store any info about a peer's chunk list for privacy and to maintain true purpose of a tracker (decentralized)
 */
type Peer struct {
	hostPort  string
	neighbors []string
}

// Peer list that keeps track of every peer in the network
var peerList []Peer = []Peer{}

/*
* Registers a new peer given the peer's host:port information.
* The tracker assigns neighbors to the peer and adds this peer to its running master list of peers.
 */
func registerPeer(hostPort string) Peer {

	// Create the new peer and add to the list
	newPeer := Peer{
		hostPort:  hostPort,
		neighbors: []string{},
	}
	// Assign neighbors
	assignNeighbors(&newPeer)
	peerList = append(peerList, newPeer)
	fmt.Printf("AllPeersInNetwork: %v\n", peerList)
	return newPeer
}

/*
* Assigns up to 10 random unique neighbors to the peer
* Or less, if there are less peers in the network
 */
func assignNeighbors(p *Peer) {
	if len(peerList) > 0 {
		neighborCount := 0
		neighborSet := make(map[string]bool)

		// populate neighborSet
		for _, peer := range peerList {
			neighborSet[peer.hostPort] = false
		}

		// populate neighbor list
		maxNumPeers := 2 // Do 10 for large networks!!!!! 2 is just for ease of testing
		for neighborCount < len(peerList) && neighborCount < maxNumPeers {
			randomIdx := rand.Intn(len(peerList))
			randomNeighbor := peerList[randomIdx]

			// Avoid duplicate neighbors
			if !neighborSet[randomNeighbor.hostPort] {
				p.neighbors = append(p.neighbors, randomNeighbor.hostPort)
				neighborSet[randomNeighbor.hostPort] = true
				neighborCount++
			}
		}
	}
}

/*
* Handles incoming requests from peers requesting to connect(for registration)
 */
func handlePeerRequest(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		message := strings.TrimSpace(string(buf[:n]))
		fmt.Println("~~~~~~~~")
		fmt.Printf("Received: %s\n", message)

		if strings.HasPrefix(message, "REGISTER") {
			// Register the peer
			args := strings.Split(message, " ")
			if len(args) < 2 {
				log.Println("Invalid REGISTER message")
				return
			}
			hostPort := args[2]

			newPeer := registerPeer(hostPort)

			// Respond with neighbors
			response := fmt.Sprintf("REGISTERED NEIGHBORS: %s\n", strings.Join(newPeer.neighbors, ","))
			conn.Write([]byte(response))
		} else if message == "LEAVE" {
			/*
				NOTE:

				We have not implemented a leaving peer since we will assume peers will stay in the network
				to continue to contribute files to the network(be helpful, good peers).
				In the future we could implement a leaving peer: Theoritically a peer would send over the
				connection for LEAVE message and just leave. Tracker will handle remove peer from the peer list to
				avoid mistaking it as alive.
			*/
			//removePeer(conn)
			fmt.Printf("Peer left: %s\n", conn.RemoteAddr())
		} else {
			log.Println("message:", message)
		}

		data := strings.ToUpper(string(buf[:n]))
		_, err = conn.Write([]byte(data))
		if err != nil {
			log.Println(err)
			return
		}
	}
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
	fmt.Println("Tracker is listening on port 8000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		// start a new go routine for each incoming peer
		go handlePeerRequest(conn)
	}
}
