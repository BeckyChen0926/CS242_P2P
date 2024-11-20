package main

import (
	"net"
)

type Peer struct {
	conn      net.Conn
	IPAddr    net.IP
	neighbors []Peer
}
