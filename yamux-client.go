package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/hashicorp/yamux"
)

func stream(session *yamux.Session, name string) {

	stream, err := session.Open()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		n, err := stream.Write([]byte("hello " + name))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s %d bytes written\n", name, n)
		time.Sleep(time.Second)
	}
}

func main() {

	conn, err := net.Dial("tcp4", "localhost:3000")
	if err != nil {
		log.Fatalf("TCP dial: %s", err)
	}

	// Setup client side of yamux
	session, err := yamux.Client(conn, nil)
	if err != nil {
		log.Fatal(err)
	}

	go stream(session, "foo")
	go stream(session, "bar")
	stream(session, "zip")
}
