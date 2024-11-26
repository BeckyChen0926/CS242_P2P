// becky version for testing peerServer

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// Peer structure for maintaining peer information
type Peer struct {
	conn      net.Conn //tbh i don't think we need to initiate conn... it's the connection to this Peer
	IPAddr    net.IP
	neighbors []Peer
	chunks    []string //chunks that the peer now has. TODO: find write type
	// MAKE PEER ID's instead
	// peer only keeping the name of chunk file now i suppose... where is it keeping the actual file??
}

var self Peer

// Register the peer with the tracker
func registerWithTracker(trackerAddr string) {
	conn, err := net.Dial("tcp", trackerAddr)
	if err != nil {
		log.Fatal("Error connecting to tracker:", err)
	}
	defer conn.Close()

	// Send registration message with chunks
	// TODO: maybe change to send ipaddr + port or something else
	chunks := strings.Join(self.chunks, ",")
	fmt.Fprintf(conn, "REGISTER %s\n", chunks)

	// Receive response from tracker
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal("Error reading from tracker:", err)
	}
	message = strings.TrimSpace(message)
	fmt.Println("Tracker response:", message)

	// Parse assigned neighbors
	if strings.HasPrefix(message, "REGISTERED NEIGHBORS:") {
		neighborIPs := strings.Split(strings.TrimPrefix(message, "REGISTERED NEIGHBORS: "), ",")
		for _, ip := range neighborIPs {
			self.neighbors = append(self.neighbors, Peer{IPAddr: net.ParseIP(ip)})
		}
		fmt.Printf("Assigned neighbors: %v\n", self.neighbors)
	}
}

// Request a chunk from the neighbor
func requestChunkFromNeighbor(conn net.Conn) {
	// Get the chunk name to search for
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter chunk name to search: ")
	chunkName, _ := reader.ReadString('\n')
	chunkName = strings.TrimSpace(chunkName)

	// Send the chunk name to the neighbor
	_, err := fmt.Fprintf(conn, chunkName+"\n")
	if err != nil {
		log.Println("Error sending chunk request:", err)
	}
}

// Download a chunk from the peer server
func downloadChunk(conn net.Conn, chunkName string) {
	// Read the chunk content from the server
	buffer := make([]byte, 1024)
	fileContent := strings.Builder{}

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Println("Error reading chunk:", err)
			break
		}
		fileContent.Write(buffer[:n])
		break
	}

	// Save the chunk to a file
	filePath := fmt.Sprintf("./%s", chunkName)
	err := os.WriteFile(filePath, []byte(fileContent.String()), 0644)
	if err != nil {
		log.Println("Error saving chunk:", err)
		return
	}

	fmt.Printf("Chunk '%s' downloaded and saved to %s\n", chunkName, filePath)
}

func main() {

	mockNeighbor4 := Peer{
		IPAddr:    net.ParseIP("127.0.0.4"), // mock IP
		neighbors: []Peer{},
		chunks:    []string{"3"}, // mock chunks this neighbor has
	}

	mockNeighbor3 := Peer{
		// conn:      nil,
		IPAddr:    net.ParseIP("127.0.0.3"),
		neighbors: []Peer{},
		chunks:    []string{"3", "2"},
	}

	// Initialize self
	self = Peer{
		IPAddr:    net.ParseIP("127.0.0.2"),
		chunks:    []string{"1"},
		neighbors: []Peer{mockNeighbor4, mockNeighbor3},
	}

	// Register with tracker
	trackerAddr := "127.0.0.1:8000"
	registerWithTracker(trackerAddr)

	port := ":8004"
	// this should be where it established TCP conn to all it's assigned neighbors
	// TODO: neighbor is hard coded for now
	conn, err := net.Dial("tcp", "127.0.0.4"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println("Connected to peer server at", port)

	for {
		// Request a chunk from the peer server
		requestChunkFromNeighbor(conn)

		// Process response
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println("Error reading from server:", err)
			break
		}
		message = strings.TrimSpace(message)

		// Process the response
		if strings.HasPrefix(message, "FOUND") {
			chunkName := strings.TrimPrefix(message, "FOUND ")
			fmt.Printf("Chunk found: %s.\nDownloading...\n", chunkName)
			downloadChunk(conn, chunkName)
		} else if strings.HasPrefix(message, "NOT FOUND") {
			chunkName := strings.TrimPrefix(message, "NOT FOUND ")
			fmt.Printf("Chunk %s not found on peer\n", chunkName)
		} else {
			fmt.Printf("Unknown response: %s\n", message)
		}
	}
}
