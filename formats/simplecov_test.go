package formats

import (
	"github.com/cocov-ci/coverage-reporter/meta"
	"github.com/cocov-ci/coverage-reporter/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func makeSimpleCovParser(t *testing.T) *SimpleCov {
	root := test_helpers.GitRoot(t)
	base := filepath.Join(root, "formats", "samples", "simplecov")

	return &SimpleCov{
		meta: &meta.Metadata{Pwd: base},
	}
}

func TestSimpleCov_Wants(t *testing.T) {
	cov := makeSimpleCovParser(t)
	t.Run("filter wanted", func(t *testing.T) {
		wanted := cov.Wants([]string{
			"foo/bar.json",
			"data/bla.rb",
			"coverage/coverage.json",
		})

		assert.NotNil(t, wanted)
		assert.Equal(t, "coverage/coverage.json", *wanted)
	})

	t.Run("handles read errors, json parsing failure, and non-simplecov files", func(t *testing.T) {
		wanted := cov.Wants([]string{
			"foo/coverage.json",
			"other_samples/not_json/coverage.json",
			"other_samples/not_simplecov/coverage.json",
		})
		assert.Nil(t, wanted)
	})
}

func TestSimpleCov_Name(t *testing.T) {
	cov := SimpleCov{}
	assert.Equal(t, "SimpleCov", cov.Name())
}

func TestSimpleCov_Parse(t *testing.T) {
	cov := makeSimpleCovParser(t)

	t.Run("read errors", func(t *testing.T) {
		v, err := cov.Parse("foo/bar.json")
		assert.Nil(t, v)
		assert.ErrorContains(t, err, "no such file or directory")
	})

	t.Run("json errors", func(t *testing.T) {
		v, err := cov.Parse("other_samples/not_json/coverage.json")
		assert.Nil(t, v)
		assert.ErrorContains(t, err, "invalid character")
	})

	t.Run("parse result", func(t *testing.T) {
		v, err := cov.Parse("coverage/coverage.json")
		assert.NotNil(t, v)
		require.NoError(t, err)
		assert.Len(t, v, 1)
		assert.Equal(t, "AAAxHjEeMQAxHjIAADEeMzQAADEAODgAADg4ADg4AAAAMQAb", v["app/controllers/application_controller.rb"])
	})
}
