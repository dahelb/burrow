package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
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

type GopherHole interface {
	Serve(ctx context.Context, selector string, conn net.Conn) error
}

type StaticGopherHole struct {
	files map[string]string
}

var NotFound = errors.New("resource not found")

func (h *StaticGopherHole) Serve(ctx context.Context, selector string, conn net.Conn) error {

	file, ok := h.files[selector]

	if !ok {
		return NotFound
	}

	if deadline, ok := ctx.Deadline(); ok {
		conn.SetWriteDeadline(deadline)
	}

	_, err := io.WriteString(conn, file)
	return err
}

func readSelector(ctx context.Context, conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)

	if deadline, ok := ctx.Deadline(); ok {
		conn.SetReadDeadline(deadline)
	}

	line, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	selector := strings.TrimRight(line, "\r\n")

	return selector, nil
}

func handleConnection(conn net.Conn, hole GopherHole) {
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	selector, err := readSelector(ctx, conn)

	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("Got selector '%v'", selector)

	err = hole.Serve(ctx, selector, conn)

	if err != nil {
		switch err {
		case NotFound:
			io.WriteString(conn, "3File Not Found\terror\terror.host\t1\r\n.\r\n")
		default:
			io.WriteString(conn, "3Server error\terror\terror.host\t1\r\n.\r\n")
		}
	}
}
