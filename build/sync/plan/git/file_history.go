package git

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// ErrNotFound signifies the file was not found.
var ErrNotFound = fmt.Errorf("not found")

// FileHistory will trace all the versions of a file in the git repository
// and return a list of sha1 hashes of that file.
func FileHistory(path string, repo *git.Repository) ([]string, error) {
	logOpts := git.LogOptions{
		FileName: &path,
	}
	commits, err := repo.Log(&logOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits for path %q: %w", path, err)
	}
	defer commits.Close()

	hashHistory := []string{}
	err = commits.ForEach(func(c *object.Commit) error {
		root, err := repo.TreeObject(c.TreeHash)
		if err != nil {
			return fmt.Errorf("failed to get commit tree: %w", err)
		}
		f, err := traverseTree(root, path)
		if err != nil {
			return err
		}
		sum, err := getReaderHash(f)
		f.Close()
		if err != nil {
			return err
		}
		hashHistory = append(hashHistory, sum)
		return nil
	})
	if err != nil && err != io.EOF {
		return nil, err
	}
	if len(hashHistory) == 0 {
		return nil, fmt.Errorf("file %q: %w", path, ErrNotFound)
	}
	return hashHistory, nil

}

func traverseTree(root *object.Tree, path string) (io.ReadCloser, error) {
	dirName, fileName := filepath.Split(path)
	var err error
	t := root
	if dirName != "" {
		t, err = root.Tree(filepath.Clean(dirName))
		if err != nil {
			return nil, fmt.Errorf("failed to traverse tree to %q: %w", dirName, err)
		}
	}
	f, err := t.File(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup file %q: %w", fileName, err)
	}
	reader, err := f.Reader()
	if err != nil {
		return nil, fmt.Errorf("failed to open %q: %w", path, err)
	}
	return reader, nil
}

func getReaderHash(r io.Reader) (string, error) {
	h := sha1.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// GetFileHash calculates the sha1 hash sum of the file.
func GetFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	sum, err := getReaderHash(f)
	if err != nil {
		return "", err
	}
	return sum, nil
}
