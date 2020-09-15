package git_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	git "gopkg.in/src-d/go-git.v4"

	gitutil "github.com/mattermost/mattermost-plugin-starter-template/build/sync/plan/git"
)

func TestFileHistory(t *testing.T) {
	assert := assert.New(t)

	repo, err := git.PlainOpenWithOptions("./", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	assert.Nil(err)
	sums, err := gitutil.FileHistory("build/sync/plan/git/testdata/testfile.txt", repo)
	assert.Nil(err)
	assert.Contains(sums, "ba7192052d7cf77c55d3b7bf40b350b8431b208b")

	// Calling with a non-existent file returns error.
	sums, err = gitutil.FileHistory("build/sync/plan/git/testdata/nosuch_testfile.txt", repo)
	assert.Equal(gitutil.ErrNotFound, err)
	assert.Nil(sums)
}
