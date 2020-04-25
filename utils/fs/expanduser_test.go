package fs

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestExpandUserNothingToExpand(t *testing.T) {
	expected := "does/not/expand"
	path, err := ExpandUser(expected)

	assert.NilError(t, err)
	assert.Equal(t, expected, path)
}

func TestExpandUserJustTilde(t *testing.T) {
	path, err := ExpandUser("~")

	assert.NilError(t, err)
	assert.Equal(t, os.Getenv("HOME"), path)
}

func TestExpandUserCurrentUser(t *testing.T) {
	path, err := ExpandUser("~/rest/of/path")
	expected := os.Getenv("HOME") + "/rest/of/path"

	assert.NilError(t, err)
	assert.Equal(t, expected, path)
}

func TestExpandUserOtherUser(t *testing.T) {
	_, err := ExpandUser("~otheruser/rest/of/path")

	assert.Check(t, is.Error(err, "expanding ~user/ paths are not supported yet"))
}
