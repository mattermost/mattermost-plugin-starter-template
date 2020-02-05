package plan

import (
	"fmt"
)

// NilCheck is used for testing only.
type NilCheck struct {
	Params struct {
		Echo string `json:"echo"`
	}
}

func (NilCheck) Check() error {
	println("ok")
	return nil
}

// Check whether the git repository is clean.
type RepoIsCleanChecker struct {
	Params struct {
		Repo string
	}
}

// Check implements the Checker interface.
func (r RepoIsCleanChecker) Check(ctx Setup) error {
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
		return fmt.Errorf("%q repository is not clean", r.Params.Repo)
	}
	return nil

}
