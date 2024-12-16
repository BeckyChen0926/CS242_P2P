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
	conn      net.Conn //TODO: tbh i don't think we need to initiate conn... it's the connection to this Peer
	IPAddr    net.IP
	neighbors []Peer
	chunks    []string
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

// Request a file from the neighbor
func requestFileFromNeighbor(conn net.Conn, fileName string) []string {
	// Send the file search request
	fmt.Fprintf(conn, "SEARCH FILE %s\n", fileName)

	// Read response from server
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println("Error reading response from server:", err)
		return nil
	}
	message = strings.TrimSpace(message)

	if strings.HasPrefix(message, "FOUND FILE:") {
		// Extract the list of chunks
		chunks := strings.Split(strings.TrimPrefix(message, "FOUND FILE: "), ",")
		fmt.Printf("Chunks found for file %s: %v\n", fileName, chunks)
		return chunks
	} else if strings.HasPrefix(message, "FILE NOT FOUND") {
		fmt.Printf("File %s not found on peer\n", fileName)
	} else {
		fmt.Printf("Unexpected response: %s\n", message)
	}
	return nil
}

// Download chunks from the server
func downloadChunks(conn net.Conn, chunks []string) {
	// skip download for duplicated chunk
	for _, chunkName := range chunks {
		// Check if the chunk already exists in self.chunks
		isDuplicate := false
		for _, selfChunk := range self.chunks {
			if selfChunk == chunkName {
				fmt.Printf("Already have Chunk '%s'. Skipping download.\n", chunkName)
				isDuplicate = true
				break
			}
		}
		if isDuplicate {
			continue // Skip the rest of the loop for this chunk
		}

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
				return
			}
			fileContent.Write(buffer[:n])
			break
		}

		// Save the chunk to a file
		filePath := fmt.Sprintf("./%s", chunkName)
		err := os.WriteFile(filePath, []byte(fileContent.String()), 0644)
		if err != nil {
			log.Printf("Error saving chunk %s: %v\n", chunkName, err)
			return
		}

		fmt.Printf("Chunk '%s' downloaded and saved to %s\n", chunkName, filePath)
		self.chunks = append(self.chunks, chunkName)
	}
}

// Check if all required chunks for a file are downloaded
func checkCompletion(fileName string) bool {
	fmt.Println("My chunks: ", self.chunks)
	expectedChunks := []string{fileName + "C1", fileName + "C2", fileName + "C3"}
	for _, expectedChunk := range expectedChunks {
		found := false
		for _, selfChunk := range self.chunks {
			if selfChunk == expectedChunk {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("Missing chunk: %s\n", expectedChunk)
			return false
		}
	}
	return false
}

/*
	NOTE:


	Our peer is currently only able to connect to and get chunks from one neighbor at a time because of our port configuration.
	We plan to implement something similar to the below code, which will create TCP connections with multiple peers.
	This functionality was not able to be tested quite yet, but we plan to incorporate it soon.


	numNeighbors := 5
	ports:= [numNeighbors]string{":8000",":8001", ":8002", ":8003", ":8004"}
	for i := 0 ; i < numNeighbors ; i++{
		conn, err := net.Dial("tcp", port)
		if err != nil {
			log.Fatal(err)
		}
	}


*/

func main() {

	mockNeighbor4 := Peer{
		IPAddr:    net.ParseIP("127.0.0.4"), // mock IP
		neighbors: []Peer{},
		chunks:    []string{"F1C3"}, // mock chunks this neighbor has
	}

	mockNeighbor3 := Peer{
		// conn:      nil,
		IPAddr:    net.ParseIP("127.0.0.3"),
		neighbors: []Peer{},
		chunks:    []string{"F1C3", "F1C2"},
	}

	// Initialize self
	self = Peer{
		IPAddr:    net.ParseIP("127.0.0.2"),
		chunks:    []string{"F1C1"},
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

	reader := bufio.NewReader(os.Stdin)

	for {
		// Prompt user to enter the file name they want to search for
		fmt.Print("Enter the file name you want to search for (e.g., F1 or F2): ")
		fileName, _ := reader.ReadString('\n')
		fileName = strings.TrimSpace(fileName)

		if fileName == "" {
			fmt.Println("File name cannot be empty. Please try again.")
			continue
		}

		// Request the file from the server
		chunks := requestFileFromNeighbor(conn, fileName)
		if len(chunks) == 0 {
			fmt.Printf("No chunks found for file %s. Retrying...\n", fileName)
			continue
		}

		// Download the chunks
		downloadChunks(conn, chunks)

		checkCompletion(fileName)
	}
}
