package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
)

// allFiles keeps a string of all the files circulating this network
var allFiles []string = []string{"F1", "F2", "F3", "F4", "F5"}

/*
* Our peer struct (for a peer) which has the following information:
* port - port number
* IPAddr - IP address
* neighbors - list of the peer's direct neighbors' host:port
* 			NOTE: a peer doesn't need to know detailed information about what chunks + neighbors its neighbor has,
*				 hence a list of type string and not peer. It retrieves information it needs through making that TCP connection.
* chunks - list of chunk names that the peer owns
 */
type Peer struct {
	port      string
	IPAddr    string
	neighbors []string
	chunks    []string
}

var host string = GetOutboundIP().String() // the host number of this peer
var port string                            // the port number of this peer

// self initialization
var self Peer = Peer{
	port:      port,
	IPAddr:    host,
	neighbors: []string{},
	chunks:    []string{},
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	// Every peer that joins the network has a client and server thread running
	go peerServerThread()
	go peerClientThread()
	wg.Wait()
}

/*
* Create a new peer's directory to store:
* uploadDirectory: the files it entered with to upload to the network
* downloadDirectory: the files it downloads from the network
 */
func createFolders() []string {
	// Create the main folder
	folderName := "Peer_" + host + "_" + strings.TrimPrefix(port, ":")
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		err = os.MkdirAll(folderName, 0755)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	// Create the upload directory
	// NOT IMPLEMENTED to be populated FOR NOW, so folder will be empty to save downloading time
	uploadDirectory := folderName + "/upload"
	if _, err := os.Stat(uploadDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDirectory, 0755)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	//Create the downpload directory
	downloadDirectory := folderName + "/download"
	if _, err := os.Stat(downloadDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(downloadDirectory, 0755)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	directories := []string{uploadDirectory, downloadDirectory}
	return directories
}

/*
* the peerServerThread is the goroutine that takes care of the peer's server behaviors.
* the peer establishes its host:port combination and listens to incoming connections.
 */
func peerServerThread() {
	if len(os.Args) == 1 {
		fmt.Println("Please provide :port")
		os.Exit(1)
	}
	host := GetOutboundIP().String()
	port = os.Args[1]
	self.port = strings.TrimPrefix(os.Args[1], ":")
	hostPort := host + port
	addr, err := net.ResolveTCPAddr("tcp", hostPort)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Println("Listening at host ", addr.IP, " and port ", addr.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}

}

/*
* the peerClientThread is the goroutine that takes care of the peer's client behaviors.
* the peer first register's itself with the tracker, creates its peer-specific directory,
* and searches for files to download from the network.
 */
func peerClientThread() {

	// Register with tracker & assign files
	trackerAddr := host + ":8000"
	if port == "" {
		registerWithTracker(trackerAddr)
		assignFiles()
	}

	// create folder directory, if not already created
	args := createFolders()
	downloadDirectory := args[1]

	reader := bufio.NewReader(os.Stdin)

	for {
		// Prompt user to enter the file name they want to search for
		fmt.Println("~~~~~~~~")
		fmt.Print("Enter the file name you want to search for (e.g., F1 or F2): ")
		fileName, _ := reader.ReadString('\n')
		fileName = strings.TrimSpace(fileName)
		if fileName == "" {
			fmt.Println("File name cannot be empty. Please try again.")
			continue
		}

		isFileDownloaded := fullTransactionFromNeighbors(self.neighbors, fileName, downloadDirectory, 5)
		if !isFileDownloaded {
			fmt.Printf("No chunks found for file %s.\n", fileName)
		}
	}
}

/*
* Recursive function that either successfully downloads the file wanted or doesnt (return bool)
* We set a TTL (time to live) that tells us how many layer of BFS we would like to search
* If a neighbor we are looking at doesn't have the file we need, we save that neighbors' neighbors in a running list
* If no neighbor in the current layer of BFS has the file we need, we look in the next layer,
* the current peer directly connects to the new layers' neighbors
* if found, we download the chunk and exit the recursion returning true.
* if the TTL runs out, we don't keep searching and assume file cannot be found.
 */
func fullTransactionFromNeighbors(neighbors []string, fileName string, downloadDirectory string, TTL int) bool {
	if TTL > 0 {
		var returnedList []string
		var indirectNeighbors []string

		// populate the original neighbors are part of the first layer, so add them as already checked
		addedNeighborsSet := make(map[string]bool)
		for _, neighbor := range self.neighbors {
			addedNeighborsSet[neighbor] = true
		}

		for _, neighborHostPort := range neighbors {
			// reach out - create TCP connection with the current neighbor them to request a file
			conn, err := net.Dial("tcp", neighborHostPort)
			if err != nil {
				log.Fatal(err)
			}

			returnedList = requestFileFromNeighbor(conn, fileName)

			if len(returnedList) > 0 {
				chars := strings.Split(returnedList[0], "")
				var firstChar string = chars[0]
				if firstChar == "F" {
					// returned a list of chunks
					fmt.Printf("Found file %s.\n", fileName)
					// Download the chunks
					downloadChunks(conn, returnedList, downloadDirectory)
					return true
				} else {
					// returned list of indirect neighbors, add these as the next layer of neighbors to explore
					for _, currIN := range returnedList {
						// Avoid duplicate neighbors that have already been added
						if !addedNeighborsSet[currIN] {
							indirectNeighbors = append(indirectNeighbors, currIN)
							addedNeighborsSet[currIN] = true
						}
					}
				}
			}
		}
		if len(indirectNeighbors) > 0 {
			fmt.Println("Searching another layer deeper into the network...")
			return fullTransactionFromNeighbors(indirectNeighbors, fileName, downloadDirectory, TTL-1)
		}
	}
	return false
}

/*
* PEER SERVER CODE
* Handle a peer's client thread requesting to connect to this peer's server thread
* the peerClient requests, and either send it over or let it know that I don't have it
 */
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
		fmt.Printf("Message: %s\n", request)

		if strings.HasPrefix(request, "SEARCH FILE") {
			fileName := strings.TrimPrefix(request, "SEARCH FILE ")
			availableChunks := searchFile(fileName)
			if len(availableChunks) > 0 {
				fmt.Printf("File %s chunks found: %v\n", fileName, availableChunks)
				sendChunks(conn, availableChunks)
			} else {
				fmt.Printf("File %s not found!\n", fileName)
				conn.Write([]byte(fmt.Sprintf("FILE NOT FOUND, SENDING NLIST: %s\n", strings.Join(self.neighbors, ","))))
			}
			// return a client view for this peer so they may request more files
			go peerClientThread()
		} else if strings.HasPrefix(request, "ADD FRIEND") {
			friendHostPort := strings.TrimPrefix(request, "ADD FRIEND ")
			if !peerExistsInList(self.neighbors, friendHostPort) {
				self.neighbors = append(self.neighbors, friendHostPort)
			}
		} else {
			fmt.Printf("Invalid request: %s\n", request)
			conn.Write([]byte("INVALID REQUEST\n"))
		}
	}
}

/*
* PEER SERVER CODE
* Searches the peer's chunklist and returns any chunks belonging to the file requested, fileName.
 */
func searchFile(fileName string) []string {
	var fileChunks []string
	for _, chunk := range self.chunks {
		if strings.HasPrefix(chunk, fileName) {
			fileChunks = append(fileChunks, chunk)
		}
	}
	return fileChunks
}

/*
* PEER SERVER CODE
* Sends available chunks requested to the requesting peer
 */
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

/*
* PEER SERVER CODE
* Helper method: Uses a UDP dialing to figure out the localhost host number
 */
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

/*
* PEER CLIENT CODE
* Register's the peer with the network's tracker
 */
func registerWithTracker(trackerAddr string) {
	conn, err := net.Dial("tcp", trackerAddr)

	if err != nil {
		log.Fatal("Error connecting to tracker:", err)
	}
	defer conn.Close()

	// Send registration message with chunks
	fmt.Fprintln(conn, "REGISTER Peer", host+port)

	// Receive response from tracker, which includes the neighbors the tracker has assigned to this peer
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err == io.EOF {
		log.Fatal("Error reading from tracker:", err)
	}
	message = strings.TrimSpace(message)
	fmt.Println("Tracker response:", message)

	// Parse tracker's assigned neighbors
	if message == "REGISTERED NEIGHBORS:" {
		fmt.Printf("Assigned neighbors: None. (You are the first in the network!)\n")
	} else {
		if strings.HasPrefix(message, "REGISTERED NEIGHBORS:") {
			neighborIPs := strings.Split(strings.TrimPrefix(message, "REGISTERED NEIGHBORS: "), ",")
			for _, info := range neighborIPs {
				self.neighbors = append(self.neighbors, info)
			}
			fmt.Printf("Assigned neighbors: %v\n", self.neighbors)
		}
	}
}

/*
* PEER CLIENT CODE
* Assigns a subset of the files to the peer to SIMULATE an entering peer who already comes with its files to upload
 */
func assignFiles() {
	//current implementation makes new peers enter with one random full file(all 3 chunks)
	randomIdx := rand.Intn(len(allFiles))
	randomFileName := allFiles[randomIdx]

	// Add chunks to chunk list
	self.chunks = append(self.chunks, randomFileName+"C1")
	self.chunks = append(self.chunks, randomFileName+"C2")
	self.chunks = append(self.chunks, randomFileName+"C3")

	// download chunks to upload folder for this peer ("make them appear" to simulate the peer already owning these files)

	// Print chunklist
	fmt.Println("My chunks: " + strings.Join(self.chunks, ", "))
}

/*
* PEER CLIENT CODE
* Request a file from the peer's direct neighbor & retrieve and return their neighborlist if file not found
 */
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
	} else if strings.HasPrefix(message, "FILE NOT FOUND, SENDING NLIST:") {
		fmt.Printf("File %s not found on peer. Searching further...\n", fileName)
		fixed := strings.TrimSpace(strings.TrimPrefix(message, "FILE NOT FOUND, SENDING NLIST:"))
		if fixed == "" {
			return nil
		}
		theirNeighbors := strings.Split(fixed, ",")
		return theirNeighbors
	} else {
		fmt.Printf("Unexpected response: %s\n", message)
	}
	return nil
}

/*
* PEER CLIENT CODE
* Downloads the chunks from the server into the downloadDirectory of the peer.
* Precondition: The server owns the chunks being downloaded.
 */
func downloadChunks(conn net.Conn, chunks []string, downloadDirectory string) {
	// If we are now exchanging chunks, let the providing peer add us as a direct neighbor if they havent already.
	addFriend(conn)
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
				break
			}
			fileContent.Write(buffer[:n])
			break
		}

		// Save the chunk to a file
		filePath := fmt.Sprintf("./"+downloadDirectory+"/%s", chunkName)
		err := os.WriteFile(filePath, []byte(fileContent.String()), 0644)
		if err != nil {
			log.Printf("Error saving chunk %s: %v\n", chunkName, err)
			return
		}

		fmt.Printf("Chunk '%s' downloaded and saved to %s\n", chunkName, filePath)
		fmt.Println("My chunks: ", self.chunks)
		self.chunks = append(self.chunks, chunkName)
	}
}

/*
* PEER CLIENT CODE (helper func)
* function to initiate to the server that they can add this peer as a neighbor
 */
func addFriend(conn net.Conn) {
	hostPort := self.IPAddr + ":" + self.port
	fmt.Fprintf(conn, "ADD FRIEND %s\n", hostPort)
	// we dont really care about the server replying, we just let them know they can add us into their neighborlist, if they want(more like if they can).
	// best effort, if they add us successfully, YAY! Otherwise oh well.
}

/*
* PEER SERVER CODE (helper func)
* function to see if a given peer is already in their neighbor list
 */
func peerExistsInList(list []string, peer string) bool {
	set := make(map[string]bool)
	for _, p := range list {
		set[p] = true
	}
	if !set[peer] {
		return false
	} else {
		return true
	}
}
