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

// Peer list and mutex for concurrency safety
var peerList []Peer

// var peerListMutex sync.Mutex

// Register a new peer
func registerPeer(ipAddr net.IP, chunks []string) Peer {
	// peerListMutex.Lock()
	// defer peerListMutex.Unlock()

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

// Handle incoming requests from peers
func handlePeerRequest(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("Error reading from connection:", err)
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
			removePeer(conn)
			fmt.Printf("Peer left: %s\n", conn.RemoteAddr())
		} else {
			log.Println("Unknown message:", message)
		}
	}
}

// Remove a peer from the peer list
func removePeer(conn net.Conn) {
	// peerListMutex.Lock()
	// defer peerListMutex.Unlock()

	newPeerList := []Peer{}
	for _, peer := range peerList {
		if peer.IPAddr.String() != conn.RemoteAddr().String() {
			newPeerList = append(newPeerList, peer)
		}
	}
	peerList = newPeerList
}

// Main function
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
		go handlePeerRequest(conn)
	}
}
