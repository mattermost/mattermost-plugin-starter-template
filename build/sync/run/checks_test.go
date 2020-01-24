package run_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	git "gopkg.in/src-d/go-git.v4"

	"github.com/mattermost/mattermost-plugin-starter-template/build/sync/run"
)

// Tests for the RepoIsClean checker.
func TestRepoIsClean(t *testing.T) {
	assert := assert.New(t)

	dir, err := ioutil.TempDir("", "test")
	assert.Nil(err)
	defer os.RemoveAll(dir)

	_, err = git.PlainInit(dir, false)
	assert.Nil(err)

	// Repo should be clean.
	checker := run.RepoIsClean(dir)
	assert.Nil(checker())

	// Create a file in the repository.
	err = ioutil.WriteFile(path.Join(dir, "data.txt"), []byte("lorem ipsum"), 0666)
	assert.Nil(err)
	assert.EqualError(checker(), fmt.Sprintf("git repository %q is not clean", dir))
}
