package plan

import (
	"fmt"
	"os"
)

// checkFail is a custom error type used to indicate a
// check that did not pass (but did not fail due to external
// causes.
// Use `IsCheckFail` to check if an error is a check failure.
type checkFail string

func (e checkFail) Error() string {
	return string(e)
}

func checkFailf(msg string, args ...interface{}) checkFail {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return checkFail(msg)
}

// IsCheckFail determines if an error is a check fail error.
func IsCheckFail(err error) bool {
	_, ok := err.(checkFail)
	return ok
}

// Check whether the git repository is clean.
type RepoIsCleanChecker struct {
	Params struct {
		Repo RepoId
	}
}

// Check implements the Checker interface.
func (r RepoIsCleanChecker) Check(_ string, ctx Setup) error {
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
		return checkFailf("%q repository is not clean", r.Params.Repo)
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
	absPath := ctx.PathInRepo(r.Params.Repo, path)
	_, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return checkFailf("path %q does not exist", path)
	} else if err != nil {
		return fmt.Errorf("failed to stat path %q: %w", absPath, err)
	}
	return nil
}
