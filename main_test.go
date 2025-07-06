package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		rootDir  string
		cfg      config
		expected string
	}{
		{name: "NoFilter", rootDir: "testdata", cfg: config{ext: "", size: 0, list: true}, expected: "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{name: "FilterExtensionMatch", rootDir: "testdata", cfg: config{ext: ".log", size: 0, list: true}, expected: "testdata/dir.log\n"},
		{name: "FilterExtensionsSizeMatch", rootDir: "testdata", cfg: config{ext: ".log", size: 10, list: true}, expected: "testdata/dir.log\n"},
		{name: "FilterExtensionsSizeNoMatch", rootDir: "testdata", cfg: config{ext: ".log", size: 20, list: true}, expected: ""},
		{name: "FilterExtensionNoMatch", rootDir: "testdata", cfg: config{ext: ".gz", size: 0, list: true}, expected: ""},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := run(tc.rootDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()
			if tc.expected != res {
				t.Errorf("expected: %q, got: %q instead4", tc.expected, res)
			}
		})
	}
}

func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	for k, n := range files {
		for j := 1; j <= n; j++ {
			fname := fmt.Sprintf("file%d%s", j, k)
			fpath := filepath.Join(tempDir, fname)
			if err := os.WriteFile(fpath, []byte("dummy"), 0666); err != nil {
				t.Fatal(err)
			}
		}
	}
	
	return tempDir, func() { _ = os.RemoveAll(tempDir) }
}
