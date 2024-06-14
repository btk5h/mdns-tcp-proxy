package main

import (
	"io"
	"log"
	"net"
	"sync"
)

func Proxy(from, to *net.TCPConn) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, err := io.Copy(from, to)
		if err != nil {
			log.Println(err)
		}

		// Signal peer that no more data is coming.
		_ = from.CloseWrite()
	}()
	go func() {
		defer wg.Done()
		_, err := io.Copy(to, from)
		if err != nil {
			log.Println(err)
		}

		// Signal peer that no more data is coming.
		_ = to.CloseWrite()
	}()

	wg.Wait()
}
