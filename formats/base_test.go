package formats

import (
	"github.com/cocov-ci/coverage-reporter/meta"
	"github.com/cocov-ci/coverage-reporter/test_helpers"
	"github.com/cocov-ci/coverage-reporter/tracking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestAutoFind_NoData(t *testing.T) {
	root := test_helpers.GitRoot(t)
	target := filepath.Join(root, "meta")
	f, err := tracking.FilesOn(target)
	require.NoError(t, err)

	metadata := meta.Metadata{
		Pwd:   target,
		Files: f,
	}

	r, err := AutoFindAll(f, &metadata)
	assert.Nil(t, r)
	assert.ErrorContains(t, err, "could not auto-detect")
}

func TestAutoFind_OK(t *testing.T) {
	root := test_helpers.GitRoot(t)
	target := filepath.Join(root, "formats", "samples", "simplecov")
	f, err := tracking.FilesOn(target)
	require.NoError(t, err)

	metadata := meta.Metadata{
		Pwd:   target,
		Files: f,
	}

	r, err := AutoFindAll(f, &metadata)
	require.NoError(t, err)
	assert.NotNil(t, r)
}
