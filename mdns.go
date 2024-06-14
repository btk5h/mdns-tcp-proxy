package main

import (
	"github.com/pion/mdns/v2"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
)

func NewMDNSServer() (*mdns.Conn, error) {
	addr4, err := net.ResolveUDPAddr("udp4", mdns.DefaultAddressIPv4)
	if err != nil {
		return nil, err
	}
	l4, err := net.ListenUDP("udp4", addr4)
	if err != nil {
		return nil, err
	}
	packetConnV4 := ipv4.NewPacketConn(l4)

	addr6, err := net.ResolveUDPAddr("udp6", mdns.DefaultAddressIPv6)
	if err != nil {
		return nil, err
	}
	l6, err := net.ListenUDP("udp6", addr6)
	if err != nil {
		return nil, err
	}
	packetConnV6 := ipv6.NewPacketConn(l6)

	server, err := mdns.Server(packetConnV4, packetConnV6, &mdns.Config{})
	if err != nil {
		return nil, err
	}

	return server, nil
}
