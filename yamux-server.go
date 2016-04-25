package main

import (
	"io"
	"log"
	"net"

	"github.com/hashicorp/yamux"
)

func stream(sconn net.Conn) {
	buff := make([]byte, 0xff)
	for {
		n, err := sconn.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Stream read error: %s", err)
			break
		}
		log.Printf("stream sent %d bytes: %s", n, buff[:n])
	}
}

func handle(conn net.Conn) {

	log.Printf("TCP accepted")

	// Setup server side of yamux
	session, err := yamux.Server(conn, nil)
	if err != nil {
		log.Fatalf("Yamux server: %s", err)
	}

	for {
		sconn, err := session.Accept()
		if err != nil {
			if session.IsClosed() {
				log.Printf("TCP closed")
				break
			}
			log.Printf("Yamux accept: %s", err)
			continue
		}
		go stream(sconn)
	}
}

func main() {
	l, err := net.Listen("tcp4", "0.0.0.0:3000")
	if err != nil {
		log.Fatalf("TCP server: %s", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("TCP accept: %s", err)
		}
		go handle(conn)
	}

}
