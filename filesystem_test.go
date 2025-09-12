package main

import (
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestDirListing(t *testing.T) {
	testFs := fstest.MapFS{
		"bar.txt": {},
		"foo":     {Mode: fs.ModeDir},
	}

	hole := FileSystemGopherHole{rootFs: testFs, host: "localhost", port: 7070}

	listing, err := hole.getDirListing(".")

	if err != nil {
		t.Error(err)
	}

	want := "0bar.txt\tbar.txt\tlocalhost\t7070\r\n1foo\tfoo\tlocalhost\t7070\r\n.\r\n"

	if listing != want {
		t.Errorf("Got '%s', want '%s'", listing, want)
	}
}
