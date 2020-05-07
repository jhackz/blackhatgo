package main

import (
	"io"
	"log"
	"net"
)

func handler(src net.Conn) {
	// Dial connects to the address on the named network
	dst, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalln("[-] Unable to connect to host")
	}
	// Upon handler returning, our connection from net.Dial will close
	defer dst.Close()

	// Run in goroutine to prevent io.Copy from blocking
	go func() {
		// Copy our source's output to the destination
		if _, err := io.Copy(dst, src); err != nil {
			log.Fatalln(err)
		}
	}()
	// Copy our desination's output back to our source
	if _, err := io.Copy(src, dst); err != nil {
		log.Fatalln(err)
	}

}

func main() {
	// Listen on local port 80
	// Listen announces on the local network address.
	lstn, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalln("[-] Unable to bind port")
	}

	for {
		// Accept accepts connections on the listener and serves requests
		conn, err := lstn.Accept()
		if err != nil {
			log.Fatalln("[-] Unable to accept connections")
		}

		go handler(conn)
	}
}
