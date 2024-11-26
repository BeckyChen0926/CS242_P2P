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
	chunks    []string
}

var self Peer

// handle a peerclient requesting to connect to this peerServer, search for the chunk
// the peerClient requests, and either send it over or let it know that I don't have it
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

	filePath := fmt.Sprintf("../files/%s", chunkName)

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
		conn:      nil,
		IPAddr:    net.ParseIP("127.0.0.2"),
		neighbors: []Peer{},
		chunks:    []string{"1"},
	}

	mockNeighbor3 := Peer{
		conn:      nil,
		IPAddr:    net.ParseIP("127.0.0.3"),
		neighbors: []Peer{},
		chunks:    []string{"3", "2"},
	}

	// initialize myself!! example
	self = Peer{
		IPAddr:    net.ParseIP("127.0.0.4"),
		chunks:    []string{"3"},
		neighbors: []Peer{mockNeighbor2, mockNeighbor3},
	}

	// should be unique for each peerServer
	port := ":8004" // temporary
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
