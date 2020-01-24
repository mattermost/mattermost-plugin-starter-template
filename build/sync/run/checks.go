package run

import (
	"fmt"

	gogit "gopkg.in/src-d/go-git.v4"
)

// Construct a check to determine if a git repository at the provided path
// is clean.
func RepoIsClean(path string) Check {
	return func() error {
		r, err := gogit.PlainOpenWithOptions(path, &gogit.PlainOpenOptions{
			DetectDotGit: true})
		if err != nil {
			return fmt.Errorf("failed to open git repository %q: %w", path, err)
		}
		worktree, err := r.Worktree()
		if err != nil {
			return fmt.Errorf("failed to get worktree at %q: %w", path, err)
		}
		status, err := worktree.Status()
		if err != nil {
			return fmt.Errorf("failed to get worktree status at %q: %w", path, err)
		}
		if !status.IsClean() {
			return fmt.Errorf("git repository %q is not clean", path)
		}
		return nil
	}
}
