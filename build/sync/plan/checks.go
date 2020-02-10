package plan

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/mattermost/mattermost-plugin-starter-template/build/sync/plan/git"
)

// checkFail is a custom error type used to indicate a
// check that did not pass (but did not fail due to external
// causes.
// Use `IsCheckFail` to check if an error is a check failure.
type checkFail string

func (e checkFail) Error() string {
	return string(e)
}

func CheckFailf(msg string, args ...interface{}) checkFail {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return checkFail(msg)
}

// IsCheckFail determines if an error is a check fail error.
func IsCheckFail(err error) bool {
	if err == nil {
		return false
	}
	e := errors.Unwrap(err)
	if e == nil {
		e = err
	}
	_, ok := e.(checkFail)
	return ok
}

// Check whether the git repository is clean.
type RepoIsCleanChecker struct {
	Params struct {
		Repo RepoId
	}
}

// Check implements the Checker interface.
// The path parameter is ignored because this checker checks the state of a repository.
func (r RepoIsCleanChecker) Check(_ string, ctx Setup) error {
	ctx.Logf(DEBUG, "checking if repository %q is clean", r.Params.Repo)
	rc := ctx.GetRepo(r.Params.Repo)
	repo := rc.Git
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}
	status, err := worktree.Status()
	if err != nil {
		return fmt.Errorf("failed to get worktree status: %w", err)
	}
	if !status.IsClean() {
		return CheckFailf("%q repository is not clean", r.Params.Repo)
	}
	return nil

}

// PathExistsChecker checks whether the fle or directory with the
// path exists. If it does not, an error is returned.
type PathExistsChecker struct {
	Params struct {
		Repo RepoId
	}
}

// Check implements the Checker interface.
func (r PathExistsChecker) Check(path string, ctx Setup) error {
	ctx.Logf(DEBUG, "checking if path %q exists in repo %q", path, r.Params.Repo)
	absPath := ctx.PathInRepo(r.Params.Repo, path)
	_, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return CheckFailf("path %q does not exist", path)
	} else if err != nil {
		return fmt.Errorf("failed to stat path %q: %w", absPath, err)
	}
	return nil
}

// FileUnalteredChecker checks whether the file in Repo is
// an unaltered version of that same file in ReferenceRepo.
//
// Its purpose is to check that a file has not been changed after forking a repository.
// It could be an old unaltered version, so the git history of the file is traversed
// until a matching version is found.
type FileUnalteredChecker struct {
	Params struct {
		ReferenceRepo RepoId `json:"reference-repo"`
		Repo          RepoId `json:"repo"`
	}
}

// Check implements the Checker interface.
func (f FileUnalteredChecker) Check(path string, setup Setup) error {
	setup.Logf(DEBUG, "checking if file %q has not been altered", path)
	absPath := setup.PathInRepo(f.Params.Repo, path)

	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("failed to get stat for %q: %w", absPath, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%q is a directory", absPath)
	}

	fileHashes, err := git.FileHistory(path, setup.GetRepo(f.Params.ReferenceRepo).Git)
	if err != nil {
		return err
	}

	currentHash, err := git.GetFileHash(absPath)
	if err != nil {
		return err
	}

	sort.Strings(fileHashes)
	idx := sort.SearchStrings(fileHashes, currentHash)
	if idx < len(fileHashes) && fileHashes[idx] == currentHash {
		return nil
	}
	return CheckFailf("file %q has been altered", path)
}
