package main

/*
import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

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

/*
client registers and initiate search
they receive request
server needs to act as a client sometimes - okay to allow server to act like a client in that situation
do this in separate thread so doesnt mess up other parts
create a thread , act like a client
connect to server of p2 and p3, save info at servers so 2 and 3 know 1 is its neighbor

recommended:
take away separation - have them be one process, go clientcode, go servercode = 2 threads, each do own processes

workaround:
save neighbor information in a file, look for the info from the file




func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
*/
