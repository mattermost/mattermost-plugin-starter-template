package plan_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	git "github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"

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
	checker.Params.Repo = plan.TargetRepo

	ctx := plan.Setup{
		Target: plan.RepoSetup{
			Path: dir,
			Git:  repo,
		},
	}
	assert.Nil(checker.Check("", ctx))

	// Create a file in the repository.
	err = ioutil.WriteFile(path.Join(dir, "data.txt"), []byte("lorem ipsum"), 0600)
	assert.Nil(err)
	err = checker.Check("", ctx)
	assert.EqualError(err, "\"target\" repository is not clean")
	assert.True(plan.IsCheckFail(err))
}

func TestPathExistsChecker(t *testing.T) {
	assert := assert.New(t)

	wd, err := os.Getwd()
	assert.Nil(err)

	checker := plan.PathExistsChecker{}
	checker.Params.Repo = plan.SourceRepo

	ctx := plan.Setup{
		Source: plan.RepoSetup{
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

func TestUnalteredCheckerSameFile(t *testing.T) {
	assert := assert.New(t)

	// Path to the root of the repo.
	wd, err := filepath.Abs("../../../")
	assert.Nil(err)

	gitRepo, err := git.PlainOpen(wd)
	assert.Nil(err)

	ctx := plan.Setup{
		Source: plan.RepoSetup{
			Path: wd,
			Git:  gitRepo,
		},
		Target: plan.RepoSetup{
			Path: wd,
		},
	}

	checker := plan.FileUnalteredChecker{}
	checker.Params.ReferenceRepo = plan.SourceRepo
	checker.Params.Repo = plan.TargetRepo

	// Check with the same file - check should succeed
	hashPath := "build/sync/plan/testdata/a"
	err = checker.Check(hashPath, ctx)
	assert.Nil(err)
}

func TestUnalteredCheckerDifferentContents(t *testing.T) {
	assert := assert.New(t)

	// Path to the root of the repo.
	wd, err := filepath.Abs("../../../")
	assert.Nil(err)

	gitRepo, err := git.PlainOpen(wd)
	assert.Nil(err)

	ctx := plan.Setup{
		Source: plan.RepoSetup{
			Path: wd,
			Git:  gitRepo,
		},
		Target: plan.RepoSetup{
			Path: wd,
		},
	}

	checker := plan.FileUnalteredChecker{}
	checker.Params.ReferenceRepo = plan.SourceRepo
	checker.Params.Repo = plan.TargetRepo

	// Create a file with the same suffix path, but different contents.
	hashPath := "build/sync/plan/testdata/a"
	tmpDir, err := ioutil.TempDir("", "test")
	assert.Nil(err)
	defer os.RemoveAll(tmpDir)
	fullPath := filepath.Join(tmpDir, "build/sync/plan/testdata")
	err = os.MkdirAll(fullPath, 0777)
	assert.Nil(err)
	file, err := os.OpenFile(filepath.Join(fullPath, "a"), os.O_CREATE|os.O_WRONLY, 0755)
	assert.Nil(err)
	_, err = file.WriteString("this file has different contents")
	assert.Nil(err)
	assert.Nil(file.Close())

	// Set the plugin path to the temporary directory.
	ctx.Target.Path = tmpDir

	err = checker.Check(hashPath, ctx)
	assert.True(plan.IsCheckFail(err))
	assert.EqualError(err, fmt.Sprintf("file %q has been altered", filepath.Join(tmpDir, hashPath)))

}

// TestUnalteredCheckerNonExistant tests running the unaltered file checker
// in the case where the target file does not exist. If the files has no history,
// the checker should pass.
func TestUnalteredCheckerNonExistant(t *testing.T) {
	assert := assert.New(t)
	hashPath := "build/sync/plan/testdata/a"

	// Path to the root of the repo.
	wd, err := filepath.Abs("../../../")
	assert.Nil(err)
	gitRepo, err := git.PlainOpen(wd)
	assert.Nil(err)

	// Temporary repo.
	tmpDir, err := ioutil.TempDir("", "test")
	assert.Nil(err)
	defer os.RemoveAll(tmpDir)

	trgRepo, err := git.PlainInit(tmpDir, false)
	assert.Nil(err)

	ctx := plan.Setup{
		Source: plan.RepoSetup{
			Path: wd,
			Git:  gitRepo,
		},
		Target: plan.RepoSetup{
			Path: tmpDir,
			Git:  trgRepo,
		},
	}

	checker := plan.FileUnalteredChecker{}
	checker.Params.ReferenceRepo = plan.SourceRepo
	checker.Params.Repo = plan.TargetRepo

	err = checker.Check(hashPath, ctx)
	assert.Nil(err)
}
