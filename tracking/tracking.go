package tracking

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var directoriesToIgnore = []string{
	".git",
	".hg",
	".bzr",
	".svn",
	".cargo",
	".idea",
	".nova",
	"node_modules",
}

func shouldIgnore(rel string) bool {
	for _, d := range directoriesToIgnore {
		if strings.HasPrefix(rel, d+"/") {
			return true
		}
	}

	return false
}

func FilesOn(path string) (map[string]string, error) {
	pathCh := make(chan string, 1)
	shaCh := make(chan []string, 1)
	errCh := make(chan error, 1)

	go func() {
		sha := sha256.New()

		for v := range pathCh {
			sha.Reset()
			f, err := os.Open(v)
			if err != nil {
				errCh <- err
				return
			}

			if _, err := io.Copy(sha, f); err != nil {
				errCh <- err
				return
			}

			digest := hex.EncodeToString(sha.Sum(nil))
			v = strings.TrimPrefix(v, path+"/")
			shaCh <- []string{v, digest}
		}
		close(shaCh)
	}()

	go func() {
		err := filepath.Walk(path, func(filePath string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			relativePath := strings.TrimPrefix(filePath, path+"/")
			if info.IsDir() || shouldIgnore(relativePath) {
				return nil
			}
			pathCh <- filePath
			return nil
		})

		if err != nil {
			errCh <- err
		}

		close(pathCh)
	}()

	output := map[string]string{}

loop:
	for {
		select {
		case sh, ok := <-shaCh:
			if !ok {
				break loop
			}
			output[sh[0]] = sh[1]
		case err := <-errCh:
			return nil, err
		}
	}

	return output, nil
}

func DiffFiles(first, second map[string]string) map[string]string {
	diffs := map[string]string{}

	for k1, v1 := range first {
		v2, ok := second[k1]
		if !ok {
			continue
		}

		if v2 != v1 {
			diffs[k1] = v2
		}
	}

	for k1, v1 := range second {
		_, ok := first[k1]
		if ok {
			continue
		}

		diffs[k1] = v1
	}

	return diffs
}
