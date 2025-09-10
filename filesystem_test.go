package main

import (
	"fmt"
	"os"
	"testing"
)

func TestCleanFilePath(t *testing.T) {
	cwd, err := os.Getwd()

	if err != nil {
		t.Fatal("Could not get cwd")
	}

	var tests = []struct {
		path, rootDir string
		want          string
		err           error
	}{
		{"/foo/bar", "/foo/bar", "/foo/bar", nil},
		{"/foo/bar/", "/foo/bar", "/foo/bar", nil},
		{"/foo/bar", "/foo", "/foo/bar", nil},
		{".", cwd, cwd, nil},
		{".", "..", cwd, nil},
		{"/", "/foo", "", RootDirEscape},
		{"/foo/bar/..", "/foo/bar/", "", RootDirEscape},
		{"..", ".", "", RootDirEscape},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("(path=\"%s\",rootDir=\"%s\")", tt.path, tt.rootDir)
		t.Run(testname, func(t *testing.T) {

			hole, err := newFileSystemGopherHole(tt.rootDir, "", 0)

			if err != nil {
				t.Fatal(err)
			}

			path, err := hole.cleanFilePath(tt.path)

			if path != tt.want || err != tt.err {
				t.Errorf("got (\"%v\",\"%v\"), want (\"%v\",\"%v\")", path, err, tt.want, tt.err)
			}
		})
	}

}
