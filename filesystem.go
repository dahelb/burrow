package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type FileSystemGopherHole struct {
	rootDir string
	rootFs  fs.FS
	host    string
	port    int
}

var RootDirEscape = errors.New("selector escapes root directory")

func newFileSystemGopherHole(rootDir string, host string, port int) (*FileSystemGopherHole, error) {
	absPath, err := filepath.Abs(rootDir)

	if err != nil {
		return nil, fmt.Errorf("invalid root path: %w", err)
	}

	if stat, err := os.Stat(absPath); err != nil {
		return nil, fmt.Errorf("root path does not exist: %w", err)
	} else if !stat.IsDir() {
		return nil, fmt.Errorf("root path is not a directory: %s", absPath)
	}

	return &FileSystemGopherHole{rootDir: absPath}, nil
}

func (h *FileSystemGopherHole) cleanFilePath(path string) (string, error) {
	path, err := filepath.Abs(path)

	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	_, found := strings.CutPrefix(path, h.rootDir)

	if !found {
		return "", RootDirEscape
	}

	return path, nil
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
	dirPath, err := h.cleanFilePath(dirPath)

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
		selector, err := filepath.Rel(h.rootDir, filepath.Join(dirPath, e.Name()))

		if err != nil {
			return "", err
		}

		if e.IsDir() {
			listing += fmt.Sprintf("1%v\t%v\t%v\t%v\n", e.Name(), selector, h.host, h.port)
			continue
		}

		t := mime.TypeByExtension(e.Name())

		switch t {
		case "text/plain":
			listing += fmt.Sprintf("0%v\t%v\t%v\t%v\n", e.Name(), selector, h.host, h.port)
		}
	}

	return listing, nil
}
