package formats

import (
	"github.com/cocov-ci/coverage-reporter/meta"
	"github.com/cocov-ci/coverage-reporter/test_helpers"
	"github.com/cocov-ci/coverage-reporter/tracking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func prepareLcov(t *testing.T) *Lcov {
	root := test_helpers.GitRoot(t)
	base, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	cmd := exec.Command("tar", "-xzf", filepath.Join(root, "formats", "samples", "lcov", "lcov-sample.tar.gz"))
	cmd.Dir = base
	_, err = cmd.CombinedOutput()
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.RemoveAll(base)
	})

	l, err := tracking.FilesOn(base)
	require.NoError(t, err)

	return &Lcov{
		meta: &meta.Metadata{
			Files: l,
			Pwd:   base,
		},
	}
}

func TestLcov_Wants(t *testing.T) {
	cov := prepareLcov(t)
	assert.Nil(t, cov.Wants(map[string]string{
		"lcov.out": "",
	}))

	assert.Equal(t, "coverage.lcov", *cov.Wants(map[string]string{
		"coverage.lcov": "",
	}))
}

func TestLcov_Parse(t *testing.T) {
	cov := prepareLcov(t)
	p, err := cov.Parse("coverage.lcov")
	require.NoError(t, err)
	assert.Len(t, p, 1)
	assert.Equal(t, "AAAxHjEeMQ==", p["main.go"])
}
