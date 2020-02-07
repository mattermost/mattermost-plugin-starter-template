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
func TestRepoIsCleanChecker(t *testing.T) {
	assert := assert.New(t)

	// Create a git repository in a temporary dir.
	dir, err := ioutil.TempDir("", "test")
	assert.Nil(err)
	defer os.RemoveAll(dir)
	repo, err := git.PlainInit(dir, false)
	assert.Nil(err)

	// Repo should be clean.
	checker := plan.RepoIsCleanChecker{}
	checker.Params.Repo = plan.PluginRepo

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
	err = checker.Check("", ctx)
	assert.EqualError(err, "\"plugin\" repository is not clean")
	assert.True(plan.IsCheckFail(err))
}

func TestPathExistsChecker(t *testing.T) {
	assert := assert.New(t)

	wd, err := os.Getwd()
	assert.Nil(err)

	checker := plan.PathExistsChecker{}
	checker.Params.Repo = plan.TemplateRepo

	ctx := plan.Setup{
		Template: plan.RepoSetup{
			Path: wd,
		},
	}

	// Check with existing directory.
	assert.Nil(checker.Check("testdata", ctx))

	// Check with existing file.
	assert.Nil(checker.Check("testdata/a", ctx))

	err = checker.Check("nosuchpath", ctx)
	assert.NotNil(err)
	assert.True(plan.IsCheckFail(err))
}
