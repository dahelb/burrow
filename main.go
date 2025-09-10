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

	fileMap := make(map[string]string)
	fileMap[""] = "0About\t/about\tlocalhost\t7070\r\n.\r\n"
	fileMap["/about"] = "Hello here is some information about me!\r\n.\r\n"

	staticHole := StaticGopherHole{fileMap}

	for {
		conn, err := ln.Accept()
		log.Printf("Got new connection from %s", conn.RemoteAddr())
		if err != nil {
			log.Print(err)
		}
		go handleConnection(conn, &staticHole)
	}

}

var NotFound = errors.New("resource not found")
