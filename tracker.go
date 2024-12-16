package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
)

// Peer structure for maintaining peer information
type Peer struct {
	IPAddr    net.IP
	neighbors []Peer
	chunks    []string
}

// Peer list that keeps every peer in the network
var peerList []Peer

// Register a new peer
func registerPeer(ipAddr net.IP, chunks []string) Peer {

	// Create the new peer and add to the list
	newPeer := Peer{
		IPAddr:    ipAddr,
		neighbors: []Peer{},
		chunks:    chunks,
	}
	peerList = append(peerList, newPeer)
	fmt.Printf("Peerlist: %v\n", peerList)

	// Assign neighbors
	assignNeighbors(&newPeer)

	return newPeer
}

/*
   NOTE:
   We have not implemented the functionality to assign new neigbors for every peer in intervals of 30 seconds.
   Our plan is to give each peer a new neighborhood to allow a peer to access more chunks.
   This will be implemented by creating a timer and calling assignNeighbors() below for each peer per interval.
*/

// Assign up to 10 unique neighbors to the peer
func assignNeighbors(p *Peer) {
	if len(peerList) <= 1 {
		return // No neighbors to assign if only one peer exists
	}

	neighborCount := 0
	neighborSet := make(map[string]bool)

	for neighborCount < 10 && neighborCount < len(peerList)-1 {
		randomIdx := rand.Intn(len(peerList))
		randomNeighbor := peerList[randomIdx]

		// Avoid self and duplicate neighbors
		if randomNeighbor.IPAddr.Equal(p.IPAddr) || neighborSet[randomNeighbor.IPAddr.String()] {
			continue
		}

		p.neighbors = append(p.neighbors, randomNeighbor)
		neighborSet[randomNeighbor.IPAddr.String()] = true
		neighborCount++
	}
}

// Handle incoming requests from peers requesting to connect
func handlePeerRequest(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			// log.Println("Error reading from connection:", err)
			return
		}

		message := strings.TrimSpace(string(buf[:n]))
		fmt.Printf("Received: %s\n", message)

		if strings.HasPrefix(message, "REGISTER") {
			// Register the peer
			args := strings.Split(message, " ")
			if len(args) < 2 {
				log.Println("Invalid REGISTER message")
				return
			}
			peerIP := net.ParseIP(conn.RemoteAddr().String())
			chunks := strings.Split(args[1], ",")
			newPeer := registerPeer(peerIP, chunks)

			// Respond with neighbors
			neighborIPs := []string{}
			for _, neighbor := range newPeer.neighbors {
				neighborIPs = append(neighborIPs, neighbor.IPAddr.String())
			}
			response := fmt.Sprintf("REGISTERED NEIGHBORS: %s\n", strings.Join(neighborIPs, ","))
			conn.Write([]byte(response))
		} else if message == "LEAVE" { //TODO: need some code in client and server upon message in terminal
			/*
				NOTE:

				we don't have the implementation for request to leave yet. Theoritically a peer would send over the
				connection for LEAVE message and just leave. Tracker will handle remove peer from the peer list to
				avoid mistaking it as alive.
			*/
			removePeer(conn)
			fmt.Printf("Peer left: %s\n", conn.RemoteAddr())
		} else {
			log.Println("message:", message)
		}
	}
}

// Remove a peer from the peer list
func removePeer(conn net.Conn) {
	newPeerList := []Peer{}
	for _, peer := range peerList {
		if peer.IPAddr.String() != conn.RemoteAddr().String() {
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
