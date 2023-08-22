package find

import (
	"bufio"
	"io"
	"io/fs"
	"sort"

	"os"
	"path/filepath"
	"regexp"

	"strings"
	"sync"
	"sync/atomic"

	"log/slog"
	"slices"

	"github.com/dhowden/tag"
	csmap "github.com/mhmtszr/concurrent-swiss-map"
)

const ReadWriteAll = 0644

var validFileRegex *regexp.Regexp
var multiDiscAlbumRegex *regexp.Regexp

var albumCount atomic.Int32
var trackCount atomic.Int32
var dupeCount atomic.Int32

func init() {
	validFileRegex = regexp.MustCompile(`.*\.(alac|flac|mp3|m4p|m4a)$`)
	multiDiscAlbumRegex = regexp.MustCompile(`(.*)/CD\d+$`)

	logFile, err := os.OpenFile("dupes.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, ReadWriteAll)
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

func ScanFiles(paths []string) {

	dupes := csmap.Create[string, []string]()
	visitedPaths := csmap.Create[string, bool]()

	files := findMusicFiles(paths)

	for file := range files {
		currentCount := trackCount.Add(1)
		fqnDir := sanitizePath(file)

		freshPath := !visitedPaths.Has(fqnDir)

		if freshPath {
			visitedPaths.Store(fqnDir, true)
		}

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

func processFile(dupes *csmap.CsMap[string, []string], filename string) {
	fqnDir := sanitizePath(filename)

	file, err := os.Open(filename)
	if err != nil {
		slog.Error("Error while opening file", "error", err)
	}
	defer file.Close()

	meta, err := tag.ReadFrom(file)
	if err != nil {
		slog.Error("Error while reading tag from file", "error", err)
	}

	if meta != nil {
		key := generateTagKey(meta)

		items, _ := dupes.Load(key)

		if !slices.Contains(items, fqnDir) {
			items = append(items, fqnDir)
			dupes.Store(key, items)

			if len(items) == 1 {
				albumCount.Add(1)
			} else {
				dupeCount.Add(1)
				displayDupes(dupes)
			}
		}
	} else {
		slog.Warn("Missing metadata", "file", file)
	}
}

func generateTagKey(metadata tag.Metadata) string {
	artistKey := metadata.AlbumArtist()
	if artistKey == "" {
		artistKey = metadata.Artist()
	}

	key := artistKey + ":" + metadata.Album()

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

func displayDupes(dupes *csmap.CsMap[string, []string]) {
	file, err := os.OpenFile("dupes.txt", os.O_CREATE|os.O_WRONLY, ReadWriteAll)
	if err != nil {
		slog.Error("failed creating file", "error", err)
	}

	datawriter := bufio.NewWriter(file)
	_, err = datawriter.WriteString("Duplicate Music Report\n\n")

	if err != nil {
		panic(err)
	}

	keys := make([]string, 0, dupes.Count())
	dupes.Range(func(k string, v []string) (stop bool) {
		keys = append(keys, k)
		return false
	})

	sort.Strings(keys)

	for _, dupe := range keys {
		actual, _ := dupes.Load(dupe)
		if len(actual) > 1 {
			_, err := datawriter.WriteString("Found duplicates for " + dupe + "\n")
			if err != nil {
				panic(err)
			}

			var sortedDupePaths []string
			sortedDupePaths = append(sortedDupePaths, actual...)

			slices.Sort(sortedDupePaths)

			for dupePath := range sortedDupePaths {
				_, err := datawriter.WriteString("  - " + sortedDupePaths[dupePath] + "\n")
				if err != nil {
					panic(err)
				}
			}
		}
	}

	datawriter.Flush()
	file.Close()
}
