package plan

import (
	"fmt"

	git "gopkg.in/src-d/go-git.v4"
)

// Context contains information about both parties
// in the sync: the plugin repository being updated
// and the source of the update - the template repo.
type Context struct {
	Template RepoContext
	Plugin   RepoContext
}

// GetRepo is a helper to get the required repo context.
// If the target parameter is not one of "plugin" or "template",
// the function panics.
func (c Context) GetRepo(target string) RepoContext {
	switch target {
	case "plugin":
		return c.Plugin
	case "template":
		return c.Template
	default:
		panic(fmt.Sprintf("cannot get repository context %q", target))
	}
}

// RepoContext contains relevant information
// about a single repository (either template or plugin).
type RepoContext struct {
	Git  *git.Repository
	Path string
}
