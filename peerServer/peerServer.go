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

// handle a peerclient requesting to connect to this peerServer, search for chunks of the file
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

		request := strings.TrimSpace(string(buf[:n]))
		fmt.Printf("Peer requested file: %s\n", request)

		if strings.HasPrefix(request, "SEARCH FILE") {
			fileName := strings.TrimPrefix(request, "SEARCH FILE ")
			availableChunks := searchFile(fileName)
			if len(availableChunks) > 0 {
				fmt.Printf("File %s chunks found: %v\n", fileName, availableChunks)
				sendChunks(conn, availableChunks)

			} else {
				fmt.Printf("File %s not found!\n", fileName)
				conn.Write([]byte(fmt.Sprintf("FILE NOT FOUND %s\n", fileName)))
			}
		} else {
			fmt.Printf("Invalid request: %s\n", request)
			conn.Write([]byte("INVALID REQUEST\n"))
		}
	}
}

func searchFile(fileName string) []string {
	var fileChunks []string
	for _, chunk := range self.chunks {
		if strings.HasPrefix(chunk, fileName) {
			fileChunks = append(fileChunks, chunk)
		}
	}
	return fileChunks
}

// Send available chunks of a file to the requesting peer
func sendChunks(conn net.Conn, chunks []string) {
	// Notify the peer that the file is available
	conn.Write([]byte(fmt.Sprintf("FOUND FILE: %s\n", strings.Join(chunks, ","))))

	for _, chunk := range chunks {
		filePath := fmt.Sprintf("../files/%s", chunk)

		// Read the chunk from the directory
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading chunk %s: %v\n", chunk, err)
			conn.Write([]byte(fmt.Sprintf("ERROR READING CHUNK %s\n", chunk)))
			continue
		}

		// Send the chunk content to the requesting peer
		conn.Write(fileContent)
		fmt.Printf("Sent chunk %s to peer\n", chunk)
	}
}

func main() {

	mockNeighbor2 := Peer{
		conn:      nil,
		IPAddr:    net.ParseIP("127.0.0.2"),
		neighbors: []Peer{},
		chunks:    []string{"F1C1"},
	}

	mockNeighbor3 := Peer{
		conn:      nil,
		IPAddr:    net.ParseIP("127.0.0.3"),
		neighbors: []Peer{},
		chunks:    []string{"F1C3", "F1C2"},
	}

	// initialize myself!! example
	self = Peer{
		IPAddr:    net.ParseIP("127.0.0.4"),
		chunks:    []string{"F1C3"},
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
