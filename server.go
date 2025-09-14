package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type GopherHole interface {
	Serve(ctx context.Context, selector string, conn net.Conn) error
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
