package main

import (
	"context"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
)

type FileSystemGopherHole struct {
	rootFs fs.FS
	host   string
	port   int
}

func (h *FileSystemGopherHole) Serve(ctx context.Context, selector string, conn net.Conn) error {

	fileInfo, err := os.Stat(selector)

	if err != nil {
		return err
	}

	if fileInfo.IsDir() {

	}

	return nil
}

func (h *FileSystemGopherHole) getDirListing(dirPath string) (string, error) {
	dirPath = filepath.Clean(dirPath)

	if stat, err := fs.Stat(h.rootFs, dirPath); err != nil {
		return "", fmt.Errorf("path does not exist: %w", err)
	} else if !stat.IsDir() {
		return "", fmt.Errorf("path is not a directory")
	}

	entries, err := fs.ReadDir(h.rootFs, dirPath)

	var listing string

	if err != nil {
		return "", err
	}

	for _, e := range entries {
		selector := filepath.Join(dirPath, e.Name())

		if e.IsDir() {
			listing += fmt.Sprintf("1%v\t%v\t%v\t%v\r\n", e.Name(), selector, h.host, h.port)
			continue
		}

		ext := filepath.Ext(e.Name())

		switch ext {
		case ".txt":
			listing += fmt.Sprintf("0%v\t%v\t%v\t%v\r\n", e.Name(), selector, h.host, h.port)
		default:
			return listing, fmt.Errorf("Unsupported file extension %s", ext)
		}
	}

	listing += ".\r\n"

	return listing, nil
}
