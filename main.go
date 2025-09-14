package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	portPtr := flag.Int("port", 7000, "The port to listen on.")
	binAddrPtr := flag.String("bind", "127.0.0.1", "bind to this address")
	rootPtr := flag.String("root", ".", "The root directory to serve")
	hostPtr := flag.String("host", "127.0.0.1", "Hostname for gopher menus")

	flag.Parse()

	addr := fmt.Sprintf("%v:%v", *binAddrPtr, *portPtr)

	log.Printf("Listening on %s", addr)
	ln, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}

	hole, err := CreateFsGopherHole(*rootPtr, *hostPtr, *portPtr)

	if err != nil {
		log.Fatal("error creating gopher hole: %w", err)
	}

	for {
		conn, err := ln.Accept()
		log.Printf("Got new connection from %s", conn.RemoteAddr())
		if err != nil {
			log.Print(err)
		}
		go handleConnection(conn, hole)
	}

}

var NotFound = errors.New("resource not found")
