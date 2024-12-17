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
	port      string
	IPAddr    string
	neighbors []Peer
}

// Peer list that keeps track of every peer in the network
var peerList []Peer = []Peer{}

/*
* Registers a new peer given the peer's host and port information.
* The tracker assigns neighbors to the peer and adds this peer to its running master list of peers.
 */
func registerPeer(host string, port string) Peer {

	// Create the new peer and add to the list
	newPeer := Peer{
		port:      port,
		IPAddr:    host,
		neighbors: []Peer{},
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
			neighborSet[peer.port] = false
		}

		// populate neighbor list
		for neighborCount < len(peerList) && neighborCount < 10 {
			randomIdx := rand.Intn(len(peerList))
			randomNeighbor := peerList[randomIdx]

			// Avoid duplicate neighbors
			if !neighborSet[randomNeighbor.port] {
				p.neighbors = append(p.neighbors, randomNeighbor)
				neighborSet[randomNeighbor.port] = true
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
			breakdown := strings.Split(args[2], ":")
			peerIP := breakdown[0]
			port := breakdown[1]

			newPeer := registerPeer(peerIP, port)

			// Respond with neighbors
			neighborInfos := []string{}
			for _, neighbor := range newPeer.neighbors {
				neighborInfos = append(neighborInfos, neighbor.IPAddr+":"+neighbor.port)
			}
			response := fmt.Sprintf("REGISTERED NEIGHBORS: %s\n", strings.Join(neighborInfos, ","))
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
			removePeer(conn)
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

/*
* Remove a peer from the peer list
 */
func removePeer(conn net.Conn) {
	newPeerList := []Peer{}
	for _, peer := range peerList {
		if peer.IPAddr != conn.RemoteAddr().String() {
			newPeerList = append(newPeerList, peer)
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
