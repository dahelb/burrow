package main

import (
	"errors"
	"fmt"
	"log"
	"net"
)

const PORT = 7070

func main() {
	log.Printf("Listening on :%d", PORT)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))

	if err != nil {
		log.Fatal(err)
	}

	hole, err := CreateFsGopherHole("./gopherhole/", "localhost", 7070)

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
