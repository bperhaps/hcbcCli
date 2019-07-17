package network

import (
	"fmt"
	"net"
)

// NewBroadcaster creates a new UDP multicast connection on which to broadcast
func NewBroadcaster(port string) (*net.UDPConn, error) {
	address := fmt.Sprintf("255.255.255.255:%s", port)
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil

}
