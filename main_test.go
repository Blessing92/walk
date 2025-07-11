package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		rootDir  string
		cfg      config
		expected string
	}{
		{name: "NoFilter", rootDir: "testdata",
			cfg: config{exts: []string{}, size: 0, list: true}, expected: "testdata/dir.log\ntestdata/dir2/script.sh\ntestdata/file.txt\n"},
		{name: "FilterExtensionMatch", rootDir: "testdata",
			cfg: config{exts: []string{".log"}, size: 0, list: true}, expected: "testdata/dir.log\n"},
		{name: "FilterExtensionMultiMatch", rootDir: "testdata",
			cfg: config{exts: []string{".log", ".txt"}, size: 0, list: true}, expected: "testdata/dir.log\ntestdata/file.txt\n"},
		{name: "FilterExtensionsSizeMatch", rootDir: "testdata",
			cfg: config{exts: []string{".log"}, size: 10, list: true}, expected: "testdata/dir.log\n"},
		{name: "FilterExtensionsSizeNoMatch", rootDir: "testdata",
			cfg: config{exts: []string{".log"}, size: 20, list: true}, expected: ""},
		{name: "FilterExtensionNoMatch", rootDir: "testdata",
			cfg: config{exts: []string{".gz"}, size: 0, list: true}, expected: ""},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := run(tc.rootDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()
			if tc.expected != res {
				t.Errorf("expected: %q, got: %q instead\n", tc.expected, res)
			}
		})
	}
}

func TestRunDelExtension(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         config
		extNoDelete string
		nDelete     int
		nNoDelete   int
		expected    string
	}{
		{name: "DeleteExtensionNoMatch",
			cfg: config{exts: []string{".log"}, del: true}, extNoDelete: ".gz", nDelete: 0, nNoDelete: 10, expected: ""},
		{name: "DeleteExtensionMatch",
			cfg: config{exts: []string{".log"}, del: true}, extNoDelete: "", nDelete: 10, nNoDelete: 0, expected: ""},
		{name: "DeleteExtensionMixed",
			cfg: config{exts: []string{".log"}, del: true}, extNoDelete: ".gz", nDelete: 5, nNoDelete: 5, expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				buffer    bytes.Buffer
				logBuffer bytes.Buffer
			)

			temDir, cleanup := createTempDir(t, map[string]int{
				strings.Join(tc.cfg.exts, ""): tc.nDelete,
				tc.extNoDelete:                tc.nNoDelete,
			})
			tc.cfg.wLog = &logBuffer

			defer cleanup()

			if err := run(temDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()
			if tc.expected != res {
				t.Errorf("expected: %q, got: %q instead\n", tc.expected, res)
			}

			filesLeft, err := os.ReadDir(temDir)
			if err != nil {
				t.Fatal(err)
			}

			if len(filesLeft) != tc.nNoDelete {
				t.Errorf("Expected %d files left, got %d instead\n", tc.nNoDelete, len(filesLeft))
			}

			expLogLines := tc.nDelete + 1
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
			if len(lines) != expLogLines {
				t.Errorf("Expected %d log lines, got %d instead\n", expLogLines, len(lines))
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

func TestRunArchive(t *testing.T) {
	testCases := []struct {
		name         string
		cfg          config
		extNoArchive string
		nArchive     int
		nNoArchive   int
	}{
		{name: "ArchiveExtensionNoMatch",
			cfg: config{exts: []string{".log"}}, extNoArchive: ".gz", nArchive: 0, nNoArchive: 10},
		{name: "ArchiveExtensionMatch",
			cfg: config{exts: []string{".log"}}, extNoArchive: "", nArchive: 10, nNoArchive: 0},
		{name: "ArchiveExtensionMixed",
			cfg: config{exts: []string{".log"}}, extNoArchive: ".gz", nArchive: 5, nNoArchive: 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			// Create temp dirs for RunArchive test
			tempDir, cleanup := createTempDir(t, map[string]int{
				strings.Join(tc.cfg.exts, ""): tc.nArchive,
				tc.extNoArchive:               tc.nNoArchive,
			})
			defer cleanup()

			archiveDir, cleanupArchive := createTempDir(t, nil)
			defer cleanupArchive()

			tc.cfg.archive = archiveDir

			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			var expFiles []string
			for _, ext := range tc.cfg.exts {
				pattern := filepath.Join(tempDir, fmt.Sprintf("*%s", ext))
				files, err := filepath.Glob(pattern)
				if err != nil {
					t.Fatal(err)
				}
				expFiles = append(expFiles, files...)
			}

			expOut := strings.Join(expFiles, "\n")
			res := strings.TrimSpace(buffer.String())

			if expOut != res {
				t.Errorf("expected: %q, got: %q instead\n", expOut, res)
			}

			filesArchived, err := os.ReadDir(archiveDir)
			if err != nil {
				t.Fatal(err)
			}

			if len(filesArchived) != tc.nArchive {
				t.Errorf("Expected %d files archived, got %d instead\n", tc.nArchive, len(filesArchived))
			}
		})
	}
}
