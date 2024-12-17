package main

/*
import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
)

// Register the peer with the tracker
func registerWithTracker(trackerAddr string) {
	conn, err := net.Dial("tcp", trackerAddr)
	if err != nil {
		log.Fatal("Error connecting to tracker:", err)
	}
	defer conn.Close()

	// Send registration message with chunks
	// TODO: maybe change to send ipaddr + port or something else
	//chunks := strings.Join(self.chunks, ",")

	fmt.Fprintln(conn, "REGISTER Peer_", self.IPAddr, ":", self.port, " with chunks: ", self.chunks)

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
		for _, info := range neighborIPs {
			hostAndPort := strings.Split(info, ":")
			currHost := hostAndPort[0]
			currPort := hostAndPort[1]
			self.neighbors = append(self.neighbors, Peer{IPAddr: currHost, port: currPort})
		}
		fmt.Printf("Assigned neighbors: %v\n", self.neighbors)
	}
}

func assignChunks() {
	if len(allFiles) <= 1 {
		return // No neighbors to assign if only one peer exists
	}

	fileCount := 0
	fileSet := make(map[string]bool)

	for fileCount < 3 && fileCount < len(allFiles)-1 {
		randomIdx := rand.Intn(len(allFiles))
		randomFile := allFiles[randomIdx]

		// Avoid self and duplicate neighbors
		if fileSet[randomFile] {
			continue
		}

		self.chunks = append(self.chunks, randomFile)
		fileSet[randomFile] = true
		fileCount++
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
	fmt.Fprintf(conn, "Message: %s\n", message)
	//
	//call completeOrForwardRequest
	//
	if strings.HasPrefix(message, "FOUND FILE:") {
		// Extract the list of chunks
		chunks := strings.Split(strings.TrimPrefix(message, "FOUND FILE: "), ",")
		fmt.Printf("Chunks found for file %s: %v\n", fileName, chunks)
		return chunks
	} else if strings.HasPrefix(message, "FILE NOT FOUND") {
		fmt.Printf("File %s not found on peer. Searching further...\n", fileName)
		// forward request to neighbors

	} else {
		fmt.Printf("Unexpected response: %s\n", message)
	}
	return nil
}

func completeOrForwardRequest(requesterPeer Peer, neighborConn net.Conn, fileName string, TTL int) []string {
	if TTL >= 1 {
		// Send the file search request
		fmt.Fprintf(neighborConn, "SEARCH FILE %s\n", fileName)

		// Read response from server
		message, err := bufio.NewReader(neighborConn).ReadString('\n')
		if err != nil {
			log.Println("Error reading response from server:", err)
			return nil
		}
		message = strings.TrimSpace(message)
		fmt.Fprintf(neighborConn, "Message: %s\n", message)
		if strings.HasPrefix(message, "FOUND FILE:") {
			// Extract the list of chunks
			chunks := strings.Split(strings.TrimPrefix(message, "FOUND FILE: "), ",")
			fmt.Printf("Chunks found for file %s: %v\n", fileName, chunks)
			return chunks
		} else if strings.HasPrefix(message, "FILE NOT FOUND") {
			fmt.Printf("File %s not found on peer. Searching further...\n", fileName)
			for _, neighbor := range self.neighbors {

				// reach out - create TCP connection with them to ask
				nextNeighborConn, err := net.Dial("tcp", neighbor.IPAddr+neighbor.port)
				if err != nil {
					log.Fatal(err)
				}
				return completeOrForwardRequest(requesterPeer, nextNeighborConn, fileName, TTL-1)
			}
		} else {
			fmt.Printf("Unexpected response: %s\n", message)
		}
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

*/
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
