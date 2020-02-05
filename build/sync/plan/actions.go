package plan

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// OverwriteDirectoryAction is used to completely overwrite directories.
// If the target directory exists, it will be removed first.
type OverwriteDirectoryAction struct {
	Params struct {
		// Create determines whether the target (file or directory)
		// will be created if it does not exist.
		Create bool
	}
}

// Run implements plan.Action.Run.
func (a OverwriteDirectoryAction) Run(path string, setup Setup) error {
	src := setup.PathInRepo(TemplateRepo, path)
	dst := setup.PathInRepo(PluginRepo, path)

	dstInfo, err := os.Stat(dst)
	if os.IsNotExist(err) {
		if !a.Params.Create {
			return fmt.Errorf("path %q does not exist, not creating", dst)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check path %q: %w", dst, err)
	} else {
		if !dstInfo.IsDir() {
			return fmt.Errorf("path %q is not a directory", dst)
		}
	}

	srcInfo, err := os.Stat(src)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory %q does not exist", src)
	} else if err != nil {
		return fmt.Errorf("failed to check path %q: %w", src, err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("path %q is not a directory", src)
	}

	err = CopyDirectory(src, dst)
	if err != nil {
		return fmt.Errorf("failed to copy path %q: %w", path, err)
	}
	return nil
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
