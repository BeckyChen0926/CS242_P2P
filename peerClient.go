package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// go through the port number for every peer server that is running
	// go run peerClient.go :8004
	//loop through and set variable port to all the ports that we are choosing / hardcode?
	port := ":8004"
	conn, err := net.Dial("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		// reader := bufio.NewReader(os.Stdin)
		// fmt.Print("Chunk to search: ")
		// text, _ := reader.ReadString('\n')
		// fmt.Fprintf(conn, text+"\n")

		// process the chunk requested for
		requestChunkFromNeighbor();
		// once we get the chunk
		downloadChunk();
		//print a message here
		chunkName := "File 1 Chunk 1/3"
		fmt.Print("Received " + chunkName);
	}
}

//upon start: make persistent TCP connections with all our neighbors

func requestChunkFromNeighbor(){
	//iterate through neighbor ID List / Connection list(assume TCP connections already established)
	for(x : neighborList){
		// conn, err := net.Dial("tcp", os.Args[1])
		// if err != nil {
		// 	log.Fatal(err)
		// }
		//request the chunk by sending message
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Chunk to search: ")
		text, _ := reader.ReadString('\n')
        fmt.Printf("Received: %s", string(buf[:n]))

        data := strings.ToUpper(string(buf[:n]))
        _, err = conn.Write([]byte(data))
        if err != nil {
            log.Println(err)
            return	}

		if(data == "no" ){
			fmt.Print("peer x doesn't have this chunk")
		}else{
			//received string of the data
			downloadChunk(data);
		}
}

//  kitty -- above stuff


// upon receiving 
func downloadChunk(String data){
	// gets chunk from neighbor
	// maybe use a scanner? write to File to take the string to a file chunk

}