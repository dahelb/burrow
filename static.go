package main

import (
	"context"
	"io"
	"net"
)

type StaticGopherHole struct {
	files map[string]string
}

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
