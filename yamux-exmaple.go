package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/hashicorp/yamux"
)

func main() {
	fmt.Println("Starting yamux demo")

	localAddr := "127.0.0.1:4444"
	done := make(chan bool, 0)
	go server(localAddr, done)
	<-done

	if err := client(localAddr); err != nil {
		log.Println(err)
	}

	time.Sleep(time.Second * 5)
}

func client(serverAddr string) error {
	// Get a TCP connection
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return err
	}

	// Setup client side of yamux
	log.Println("creating client session")
	session, err := yamux.Client(conn, nil)
	if err != nil {
		return err
	}

	// Open a new stream
	log.Println("opening stream")
	stream, err := session.Open()
	if err != nil {
		return err
	}

	// Stream implements net.Conn
	_, err = stream.Write([]byte("ping"))
	return err
}

func server(localAddr string, done chan bool) error {
	// Accept a TCP connection
	listener, err := net.Listen("tcp", localAddr)

	close(done)

	conn, err := listener.Accept()
	if err != nil {
		return err
	}

	// Setup server side of yamux
	log.Println("creating server session")
	session, err := yamux.Server(conn, nil)
	if err != nil {
		return err
	}

	// Accept a stream
	log.Println("accepting stream")
	stream, err := session.Accept()
	if err != nil {
		return err
	}

	// Listen for a message
	buf := make([]byte, 4)
	_, err = stream.Read(buf)

	fmt.Printf("buf = %+v\n", string(buf))
	return err
}
