package main

import (
	"bufio"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/dhowden/tag"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

var validFileRegex *regexp.Regexp
var multiDiscAlbumRegex *regexp.Regexp

var albumCount atomic.Int32
var trackCount atomic.Int32
var dupeCount atomic.Int32

var dupeMutex = &sync.RWMutex{}
var pathMutex = &sync.RWMutex{}

func init() {
	validFileRegex = regexp.MustCompile(`.*\.(alac|flac|mp3|m4p|m4a)$`)
	multiDiscAlbumRegex = regexp.MustCompile(`(.*)/CD\d+$`)

	logFile, err := os.OpenFile("dupes.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	mw := io.MultiWriter(os.Stderr, logFile)
	logger := slog.New(slog.NewTextHandler(mw, nil))
	slog.SetDefault(logger)
}

func findMusicFiles(paths []string) <-chan string {
	wg := &sync.WaitGroup{}
	ch := make(chan string)
	for pathIdx := range paths {
		path := paths[pathIdx]
		wg.Add(1)
		go func() {
			slog.Info("Scanning for duplicates", "path", path)

			fsys := os.DirFS(path)

			_ = fs.WalkDir(fsys, ".", func(file string, d fs.DirEntry, err error) error {
				if isValidFile(file) {
					fqnFile, _ := filepath.Abs(path + "/" + file)
					ch <- fqnFile
				}
				return nil
			})
			slog.Info("Finnished scanning for duplicates", "path", path)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func main() {
	paths := os.Args[1:]
	dupes := make(map[string][]string)
	visitedPaths := make(map[string]bool)
	files := findMusicFiles(paths)

	for file := range files {

		currentCount := trackCount.Add(1)
		fqnDir := sanitizePath(file)

		pathMutex.Lock()
		freshPath := !visitedPaths[fqnDir]
		if freshPath {
			visitedPaths[fqnDir] = true
		}
		pathMutex.Unlock()

		if freshPath {
			slog.Info("Checking for duplicate music", "albums",
				albumCount.Load(), "tracks", currentCount, "dupes", dupeCount.Load(), "file", file)

			go func(fileToProcess string) {
				processFile(dupes, fileToProcess)
			}(file)
		}

	}

	slog.Info("Finnished.")

}

func processFile(dupes map[string][]string, file string) {

	fqnDir := sanitizePath(file)

	f, err := os.Open(file)
	if err != nil {
		slog.Error("Error while opening file", "error", err)
	}
	defer f.Close()

	meta, err := tag.ReadFrom(f)
	if err != nil {
		slog.Error("Error while reading tag from file", "error", err)
	}

	if meta != nil {

		key := generateTagKey(meta)

		dupeMutex.Lock()
		items := dupes[key]

		if !slices.Contains(items, fqnDir) {
			items = append(items, fqnDir)
			dupes[key] = items

			if len(items) == 1 {
				albumCount.Add(1)
			} else {
				dupeCount.Add(1)
				displayDupes(dupes)
			}
		}
		dupeMutex.Unlock()
	} else {
		slog.Warn("Missing metadata", "file", file)
	}
}

func generateTagKey(m tag.Metadata) string {
	artistKey := m.AlbumArtist()
	if artistKey == "" {
		artistKey = m.Artist()
	}

	key := artistKey + ":" + m.Album()
	return key
}

func sanitizePath(file string) string {
	sanitized := filepath.Dir(file)
	sanitized = multiDiscAlbumRegex.ReplaceAllString(sanitized, "$1")
	return sanitized
}

func isValidFile(file string) bool {
	lower := strings.ToLower(file)
	return validFileRegex.MatchString(lower)
}

func displayDupes(dupes map[string][]string) {

	file, err := os.OpenFile("dupes.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("failed creating file", "error", err)
	}
	datawriter := bufio.NewWriter(file)

	datawriter.WriteString("Duplicate Music Report\n\n")

	keys := make([]string, 0, len(dupes))
	for k := range dupes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, dupe := range keys {

		if len(dupes[dupe]) > 1 {
			datawriter.WriteString("Found duplicates for " + dupe + "\n")
			dupePaths := dupes[dupe]
			for dupePath := range dupePaths {
				datawriter.WriteString("  - " + dupePaths[dupePath] + "\n")
			}
		}
	}
	datawriter.Flush()
	file.Close()

}
