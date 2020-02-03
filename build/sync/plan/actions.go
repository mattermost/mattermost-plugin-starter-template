package plan

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// OverwriteDirectoryAction is used to completely
// overwrite files or directories.
// If the target directory exists, it will be removed
// first.
type CleanOverwriteAction struct {
	Params struct {
		// Create determines whether the target (file or directory)
		// will be created if it does not exist.
		Create bool
	}
}

func CopyDirectory(src, dst string) error {
	copier := dirCopier{dst: dst, src: src}
	return filepath.Walk(src, copier.Copy)
}

type dirCopier struct {
	dst string
	src string
}

// Convert a path in the source directory to a path in the destination
// directory.
func (d dirCopier) srcToDst(path string) (string, error) {
	suff := strings.TrimPrefix(path, d.src)
	if suff == path {
		return "", fmt.Errorf("path %q is not in %q", path, d.src)
	}
	return filepath.Join(d.dst, suff), nil
}

// Copy is an implementation of filepatch.WalkFunc that copies the
// source directory to target with all subdirectories.
func (d dirCopier) Copy(srcPath string, info os.FileInfo, err error) error {
	if err != nil {
		return fmt.Errorf("failed to copy directory: %w", err)
	}
	trgPath, err := d.srcToDst(srcPath)
	if err != nil {
		return err
	}
	if info.IsDir() {
		err = os.MkdirAll(trgPath, info.Mode())
		if err != nil {
			return fmt.Errorf("failed to create directory %q: %w", trgPath, err)
		}
		err = os.Chtimes(trgPath, info.ModTime(), info.ModTime())
		if err != nil {
			return fmt.Errorf("failed to create directory %q: %w", trgPath, err)
		}
		return nil
	}
	err = copyFile(srcPath, trgPath, info)
	if err != nil {
		return fmt.Errorf("failed to copy file %q: %w", srcPath, err)
	}
	return nil
}

func copyFile(src, dst string, info os.FileInfo) error {
	srcF, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %q: %w", src, err)
	}
	defer srcF.Close()
	dstF, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, info.Mode())
	if err != nil {
		return fmt.Errorf("failed to open destination file %q: %w", dst, err)
	}
	_, err = io.Copy(dstF, srcF)
	if err != nil {
		dstF.Close()
		return fmt.Errorf("failed to copy file %q: %w", src, err)
	}
	if err := dstF.Close(); err != nil {
		return fmt.Errorf("failed to close file %q: %w", dst, err)
	}
	err = os.Chtimes(dst, info.ModTime(), info.ModTime())
	if err != nil {
		return fmt.Errorf("failed to adjust file modification time for %q: %w", dst, err)
	}
	return nil
}
