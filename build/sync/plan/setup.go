package plan

import (
	"fmt"
	"os"
	"path/filepath"

	git "gopkg.in/src-d/go-git.v4"
)

// RepoID identifies a repository - either plugin or template.
type RepoID string

const (
	// TemplateRepo is the id of the template repository (source).
	TemplateRepo RepoID = "template"
	// PluginRepo is the id of the plugin repository (target).
	PluginRepo RepoID = "plugin"
)

// LogLevel sets the level of the log mesage.
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String implements the Stringer interface.
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "<unknown log level>"
	}
}

// Setup contains information about both parties
// in the sync: the plugin repository being updated
// and the source of the update - the template repo.
type Setup struct {
	Template RepoSetup
	Plugin   RepoSetup
}

// Logf logs the provided message with the set level.
func (c Setup) Logf(lvl LogLevel, tpl string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, lvl.String()+": "+tpl+"\n", args...)
}

// GetRepo is a helper to get the required repo setup.
// If the target parameter is not one of "plugin" or "template",
// the function panics.
func (c Setup) GetRepo(r RepoID) RepoSetup {
	switch r {
	case PluginRepo:
		return c.Plugin
	case TemplateRepo:
		return c.Template
	default:
		panic(fmt.Sprintf("cannot get repository setup %q", r))
	}
}

// PathInRepo returns the full path of a file in the specified repository.
func (c Setup) PathInRepo(repo RepoID, path string) string {
	r := c.GetRepo(repo)
	return filepath.Join(r.Path, path)
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
