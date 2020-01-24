package git

import (
	"fmt"

	gogit "gopkg.in/src-d/go-git.v4"
)

// Repo provides an API to interact with a git repository.
type Repo struct {
	repo *gogit.Repository
}

// Open git repository at the given path.
func Open(path string) (*Repo, error) {
	r, err := gogit.PlainOpenWithOptions(path, &gogit.PlainOpenOptions{
		DetectDotGit: true})
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}
	return &Repo{repo: r}, nil
}

// Check if the repository worktree is clean.
func (r *Repo) IsClean() (bool, error) {
	worktree, err := r.repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("failed to get worktree: %w", err)
	}
	status, err := worktree.Status()
	if err != nil {
		return false, fmt.Errorf("failed to get worktree status: %w", err)
	}
	fmt.Println(status.String())
	return status.IsClean(), nil
}
