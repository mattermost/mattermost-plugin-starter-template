package plan

import (
	"fmt"
	"os"
	"sort"

	"github.com/mattermost/mattermost-plugin-starter-template/build/sync/plan/git"
)

// CheckFail is a custom error type used to indicate a
// check that did not pass (but did not fail due to external
// causes.
// Use `IsCheckFail` to check if an error is a check failure.
type CheckFail string

func (e CheckFail) Error() string {
	return string(e)
}

// CheckFailf creates an error with the specified message string.
// The error will pass the IsCheckFail filter.
func CheckFailf(msg string, args ...interface{}) CheckFail {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return CheckFail(msg)
}

// IsCheckFail determines if an error is a check fail error.
func IsCheckFail(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(CheckFail)
	return ok
}

// RepoIsCleanChecker checks whether the git repository is clean.
type RepoIsCleanChecker struct {
	Params struct {
		Repo RepoID
	}
}

// Check implements the Checker interface.
// The path parameter is ignored because this checker checks the state of a repository.
func (r RepoIsCleanChecker) Check(_ string, ctx Setup) error {
	ctx.Logf("checking if repository %q is clean", r.Params.Repo)
	rc := ctx.GetRepo(r.Params.Repo)
	repo := rc.Git
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %v", err)
	}
	status, err := worktree.Status()
	if err != nil {
		return fmt.Errorf("failed to get worktree status: %v", err)
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
		Repo RepoID
	}
}

// Check implements the Checker interface.
func (r PathExistsChecker) Check(path string, ctx Setup) error {
	repo := r.Params.Repo
	if repo == "" {
		repo = TargetRepo
	}
	ctx.Logf("checking if path %q exists in repo %q", path, repo)
	absPath := ctx.PathInRepo(repo, path)
	_, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return CheckFailf("path %q does not exist", path)
	} else if err != nil {
		return fmt.Errorf("failed to stat path %q: %v", absPath, err)
	}
	return nil
}

// FileUnalteredChecker checks whether the file in Repo is
// an unaltered version of that same file in ReferenceRepo.
//
// Its purpose is to check that a file has not been changed after forking a repository.
// It could be an old unaltered version, so the git history of the file is traversed
// until a matching version is found.
//
// If the repositories in the parameters are not specified,
// reference will default to the source repository and repo - to the target.
type FileUnalteredChecker struct {
	Params struct {
		ReferenceRepo RepoID `json:"compared-to"`
		Repo          RepoID `json:"in"`
	}
}

// Check implements the Checker interface.
func (f FileUnalteredChecker) Check(path string, setup Setup) error {
	setup.Logf("checking if file %q has not been altered", path)
	repo := f.Params.Repo
	if repo == "" {
		repo = TargetRepo
	}
	reference := f.Params.ReferenceRepo
	if reference == "" {
		reference = SourceRepo
	}
	absPath := setup.PathInRepo(repo, path)

	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return CheckFailf("file %q has been deleted", absPath)
	}
	if err != nil {
		return fmt.Errorf("failed to get stat for %q: %v", absPath, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%q is a directory", absPath)
	}

	fileHashes, err := git.FileHistory(path, setup.GetRepo(reference).Git)
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
	return CheckFailf("file %q has been altered", absPath)
}
