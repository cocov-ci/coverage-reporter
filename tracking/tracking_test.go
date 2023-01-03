package tracking

import (
	"github.com/cocov-ci/coverage-reporter/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestFilesOn(t *testing.T) {
	root := test_helpers.GitRoot(t)
	list, err := FilesOn(filepath.Join(root, "tracking", "fixtures", "a"))
	require.NoError(t, err)
	files := map[string]string{
		"README.md": "a0ef3e2d6478e1484b7fc6d7a153572288712150683cfd54353f0b379c477c13",
		"foo.clj":   "9f2ea84d9d4e6d8570ce8d98477164d57707eb56601d0a98bbc38a5f92e39d52",
		"other":     "06f961b802bc46ee168555f066d28f4f0e9afdf3f88174c1ee6f9de004fc30a0",
	}

	require.Equal(t, list, files)
}

func TestDiffFiles(t *testing.T) {
	root := test_helpers.GitRoot(t)
	listA, err := FilesOn(filepath.Join(root, "tracking", "fixtures", "a"))
	require.NoError(t, err)
	listB, err := FilesOn(filepath.Join(root, "tracking", "fixtures", "b"))
	require.NoError(t, err)

	diff := DiffFiles(listA, listB)
	assert.NotZero(t, diff["c.out"])
	assert.Equal(t, diff["c.out"], listB["c.out"])
	assert.NotZero(t, diff["other"])
	assert.Equal(t, diff["other"], listB["other"])
	assert.Len(t, diff, 2)
}
