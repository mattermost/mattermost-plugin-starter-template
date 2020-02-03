package plan_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattermost/mattermost-plugin-starter-template/build/sync/plan"
)

func TestCopyDirectory(t *testing.T) {
	assert := assert.New(t)

	// Create a temporary directory to copy to.
	dir, err := ioutil.TempDir("", "test")
	assert.Nil(err)
	defer os.RemoveAll(dir)

	wd, err := os.Getwd()
	assert.Nil(err)

	srcDir := path.Join(wd, "testdata")
	err = plan.CopyDirectory(srcDir, dir)
	assert.Nil(err)

	srcContents, err := ioutil.ReadDir(srcDir)
	assert.Nil(err)
	dstContents, err := ioutil.ReadDir(dir)
	assert.Nil(err)
	assert.Len(dstContents, len(srcContents))

	// Check the directory contents are equal.
	for i, srcFInfo := range srcContents {
		dstFInfo := dstContents[i]
		assert.Equal(srcFInfo.Name(), dstFInfo.Name())
		assert.Equal(srcFInfo.Size(), dstFInfo.Size())
		assert.Equal(srcFInfo.Mode(), dstFInfo.Mode())
		assert.Equal(srcFInfo.IsDir(), dstFInfo.IsDir())
	}
}
