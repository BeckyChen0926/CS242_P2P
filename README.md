# CS242_P2P

# Project Description
We coded a decentralized Peer to Peer file sharing network using Golang that allows peers to connect to peers in the network and upload + download files over TCP connections. The entities in our network include the peers and a tracker.

# How to Install and Run/Use the Project
**To get the development environment set up:** please install golang. 

**To run the project:**

0) Prior to every run, please make sure all peer directories from previous runs are deleted to reset the space.
1) Open up one terminal for the tracker and enter in the terminal ```go run tracker.go```, which starts up the tracker in our network at port :8000 on the local host.
2) Open up one terminal for each peer you would like to enter the network.
3) For each peer:

    a. First, enter ```cd peer``` .
   
    b. Then enter in the terminal ```go run peer.go <:port>``` in which the port you enter is a port number after 8000(the port reserved for the tracker). For example, ```go run peer.go :8001``` or ```go run peer.go :8002``` all work.
   
4) Now you are ready to request files as prompted from each peer's terminal.

**To use the project:**

Some variables you can tweak for your own testing purposes:

In method ```tracker.go:assignNeighbors``` line 60, you may want to set ```maxNumPeers``` to a different number to allow an incoming peer to be assigned less or more neighbors. 

In method ```peer.go:peerClientThread``` line 142, you may want to hardcode ```trackerAddr``` to a different host:port combination.

# Example Demo


# Credits
Completed by [Becky](https://github.com/BeckyChen0926), [Kitty](), and [Jessica](https://github.com/jessica-b-dai). 


