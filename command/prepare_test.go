package command

import (
	"github.com/cocov-ci/coverage-reporter/meta"
	"github.com/cocov-ci/coverage-reporter/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestPrepare(t *testing.T) {
	repo := test_helpers.GitRoot(t)
	wd := filepath.Join(repo, "script")

	t.Run("with commitish flag", func(t *testing.T) {
		ctx := test_helpers.MakeContext(t, map[string]string{
			"commitish": "foo",
		})

		test_helpers.InDir(t, wd, func(t *testing.T) {
			err := Prepare(ctx)
			require.NoError(t, err)
			data, err := meta.ReadMetadata("token")
			require.NoError(t, err)
			assert.Equal(t, "foo", data.Sha)
		})
	})

	t.Run("with github event data", func(t *testing.T) {
		ctx := test_helpers.MakeContext(t, nil)
		t.Setenv("GITHUB_EVENT_PATH", filepath.Join(repo, "models", "github_event", "fixtures", "event.json"))
		test_helpers.InDir(t, wd, func(t *testing.T) {
			err := Prepare(ctx)
			require.NoError(t, err)
			data, err := meta.ReadMetadata("token")
			require.NoError(t, err)
			assert.Equal(t, "78ac78527048240a5d4ed0bbf4d61f62e3b7017f", data.Sha)
		})
	})

	t.Run("from git data", func(t *testing.T) {
		ctx := test_helpers.MakeContext(t, nil)
		err := Prepare(ctx)
		require.NoError(t, err)
		data, err := meta.ReadMetadata("token")
		require.NoError(t, err)
		head := test_helpers.GitHead(t)
		assert.Equal(t, head, data.Sha)
	})
}
