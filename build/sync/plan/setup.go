package plan

import (
	"fmt"

	git "gopkg.in/src-d/go-git.v4"
)

// Setup contains information about both parties
// in the sync: the plugin repository being updated
// and the source of the update - the template repo.
type Setup struct {
	Template RepoSetup
	Plugin   RepoSetup
}

// GetRepo is a helper to get the required repo setup.
// If the target parameter is not one of "plugin" or "template",
// the function panics.
func (c Setup) GetRepo(target string) RepoSetup {
	switch target {
	case "plugin":
		return c.Plugin
	case "template":
		return c.Template
	default:
		panic(fmt.Sprintf("cannot get repository setup %q", target))
	}
}

// RepoSetup contains relevant information
// about a single repository (either template or plugin).
type RepoSetup struct {
	Git  *git.Repository
	Path string
}

// GetRepoSetup returns the repository setup for the specified path.
func GetRepoSetup(path string) (RepoSetup, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return RepoSetup{}, fmt.Errorf("failed to access git repository at %q: %w", path, err)
	}
	return RepoSetup{
		Git:  repo,
		Path: path,
	}, nil
}
