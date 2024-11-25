// for testing

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Peer struct {
	conn      net.Conn
	IPAddr    net.IP
	neighbors []Peer
	chunks    []string //chunks that the peer now has. TODO: find write type
	// MAKE PEER ID's instead
	// peer only keeping the name of chunk file now i suppose... where is it keeping the actual file??
}

var self Peer

// a peerclient requested to connect to this peerServer
// but the TCP connection should have been established??
// or maybe 2 cases 1) this is new neighbor and TCP conn is established here
//  2. this is old neighbor w/ established TCP so this is duplicate conn (does this matter?)
func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("Error reading from connection:", err)
			return
		}

		chunkInNeed := strings.TrimSpace(string(buf[:n]))
		fmt.Printf("Peer requested chunk: %s\n", chunkInNeed)

		if searchChunk(chunkInNeed) {
			fmt.Printf("Chunk %s found!\nSending to peer...\n", chunkInNeed)
			fmt.Fprintf(conn, "FOUND %s\n", chunkInNeed)
			sendChunk(conn, chunkInNeed)
		} else {
			fmt.Printf("Chunk %s not available!\n", chunkInNeed)
			fmt.Fprintf(conn, "NOT FOUND %s\n", chunkInNeed)
		}
	}
}

func searchChunk(chunkInNeed string) bool {
	// fmt.Printf("chunk in need:%s\n", chunkInNeed)

	for _, chunk := range self.chunks {
		if chunk == chunkInNeed {
			return true
		}
	}
	return false
}

func sendChunk(conn net.Conn, chunkName string) {

	filePath := fmt.Sprintf("../../files/%s", chunkName)

	// Read the chunk from dir
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(conn, "ERROR finding chunk: %s\n", err.Error())
		return
	}

	// Send the chunk content to the requesting peer
	_, err = conn.Write(fileContent)
	if err != nil {
		log.Println("Error sending chunk:", err)
	} else {
		fmt.Printf("Sent chunk %s to peer\n", chunkName)
	}
}

func main() {

	mockNeighbor2 := Peer{
		// conn:      nil,
		IPAddr:    net.ParseIP("127.0.0.2"), // Example IP
		neighbors: []Peer{},
		chunks:    []string{"1"}, // Example chunks this neighbor has
	}

	mockNeighbor4 := Peer{
		// conn:      nil,

		IPAddr: net.ParseIP("127.0.0.4"), // Use the server's IP address
		// chunks:  []string{"chunk1.txt", "chunk2.txt"}, // Replace with actual file names
		chunks:    []string{"3"}, // Replace with actual file names
		neighbors: []Peer{},      // This will be populated dynamically
	}

	// initialize myself!! example
	self = Peer{
		IPAddr:    net.ParseIP("127.0.0.3"),
		neighbors: []Peer{mockNeighbor4, mockNeighbor2},
		chunks:    []string{"3", "2"},
	}

	// TODO: mayeb should loop for each port? rn can only listen on one port/ listen to one peerclient
	// at a time - bad

	// should be unique for each peerServer
	port := ":8005" // temporary
	addr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	fmt.Println("Listening on port", port)
	fmt.Printf("Connected neighbors: %v\n", self.neighbors)
	fmt.Printf("Available chunks: %v\n", self.chunks)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}

}
