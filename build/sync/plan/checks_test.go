package plan_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	git "gopkg.in/src-d/go-git.v4"

	"github.com/mattermost/mattermost-plugin-starter-template/build/sync/plan"
)

// Tests for the RepoIsClean checker.
func TestRepoIsClean(t *testing.T) {
	assert := assert.New(t)

	// Create a git repository in a temporary dir.
	dir, err := ioutil.TempDir("", "test")
	assert.Nil(err)
	defer os.RemoveAll(dir)
	repo, err := git.PlainInit(dir, false)
	assert.Nil(err)

	// Repo should be clean.
	checker := plan.RepoIsCleanChecker{}
	checker.Params.Repo = "plugin"

	ctx := plan.Setup{
		Plugin: plan.RepoSetup{
			Path: dir,
			Git:  repo,
		},
	}
	assert.Nil(checker.Check("", ctx))

	// Create a file in the repository.
	err = ioutil.WriteFile(path.Join(dir, "data.txt"), []byte("lorem ipsum"), 0666)
	assert.Nil(err)
	assert.EqualError(checker.Check("", ctx), "\"plugin\" repository is not clean")
}
