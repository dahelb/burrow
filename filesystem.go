package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"os"
	"path/filepath"
)

type FileSystemGopherHole struct {
	rootFs fs.FS
	host   string
	port   int
}

func CreateFsGopherHole(rootDirPath string, host string, port int) (*FileSystemGopherHole, error) {

	rootDirPath, err := filepath.Abs(rootDirPath)

	if err != nil {
		return nil, fmt.Errorf("error getting absolute path: %w", err)
	}

	root, err := os.OpenRoot(rootDirPath)
	if err != nil {
		return nil, err
	}

	log.Printf("create root file system at %s", rootDirPath)

	return &FileSystemGopherHole{
		rootFs: root.FS(),
		host:   host,
		port:   port,
	}, nil

}

func (h *FileSystemGopherHole) Serve(ctx context.Context, selector string, conn net.Conn) error {

	filePath := filepath.Clean(selector)

	if selector == "/" || selector == "" {
		filePath = "."
	}

	fileInfo, err := fs.Stat(h.rootFs, filePath)

	if err != nil {
		log.Printf("error getting stats: %s", err)
		return err
	}

	if fileInfo.IsDir() {
		log.Printf("getting dir listing for %s", filePath)
		listing, err := h.getDirListing(filePath)

		if err != nil {
			log.Printf("Error getting dir listing: %s", err)
			return err
		}

		_, err = io.WriteString(conn, listing)

		return err
	} else {
		file, err := h.rootFs.Open(selector)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(conn, file)
	}

	return err
}

func (h *FileSystemGopherHole) getDirListing(dirPath string) (string, error) {
	if stat, err := fs.Stat(h.rootFs, dirPath); err != nil {
		return "", fmt.Errorf("error getting file stat: %w", err)
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
			listing += fmt.Sprintf("9%v\t%v\t%v\t%v\r\n", e.Name(), selector, h.host, h.port)
		}
	}

	listing += ".\r\n"

	return listing, nil
}
