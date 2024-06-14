package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/pion/mdns/v2"
	"golang.org/x/net/context"
	"log"
	"net"
	"os"
	"strconv"
)

var opts struct {
	Target       string `short:"t" long:"target" description:"The server to proxy traffic to" required:"true"`
	TargetPort   int    `short:"P" long:"target-port" description:"The port to serve traffic to" required:"true"`
	ListenerPort int    `short:"p" long:"port" description:"The port to listen on" required:"true"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	mdnsServer, err := NewMDNSServer()
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("tcp4", ":"+strconv.Itoa(opts.ListenerPort))
	if err != nil {
		log.Fatal(err)
	}
	defer func(l net.Listener) {
		_ = l.Close()
	}(l)

	ctx := context.Background()

	fmt.Printf("Proxying requests on port %d -> %s:%d\n", opts.ListenerPort, opts.Target, opts.TargetPort)

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(ctx, c.(*net.TCPConn), mdnsServer)
	}
}

func handleConnection(ctx context.Context, c *net.TCPConn, mdnsServer *mdns.Conn) {
	fmt.Println("Handling connection from ", c.RemoteAddr())

	_, target, err := mdnsServer.QueryAddr(ctx, opts.Target)
	if err != nil {
		log.Println("Error while looking resolving target: ", err)
	}

	fmt.Printf("Proxying from %s to %s\n", c.RemoteAddr(), target)

	targetConn, err := net.Dial("tcp", target.String()+":"+strconv.Itoa(opts.TargetPort))
	if err != nil {
		log.Println("Error while handling request:", err)
		return
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(targetConn)

	Proxy(c, targetConn.(*net.TCPConn))
}
