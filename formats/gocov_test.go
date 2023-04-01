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

func prepareGoCov(t *testing.T) *GoCov {
	root := test_helpers.GitRoot(t)
	base, err := os.MkdirTemp("", "")
	require.NoError(t, err)

	cmd := exec.Command("tar", "-xzf", filepath.Join(root, "formats", "samples", "gocov", "gocov-sample.tar.gz"))
	cmd.Dir = base
	_, err = cmd.CombinedOutput()
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.RemoveAll(base)
	})

	l, err := tracking.FilesOn(base)
	require.NoError(t, err)

	return &GoCov{
		Meta: &meta.Metadata{
			Files: l,
			Pwd:   base,
		},
		mods: nil,
	}
}

func TestGoCov_Wants(t *testing.T) {
	cov := prepareGoCov(t)
	assert.Nil(t, cov.Wants([]string{
		"coverage.out",
	}))

	assert.Equal(t, "c.out", *cov.Wants([]string{
		"c.out",
	}))
}

func TestGoCov_Parse(t *testing.T) {
	cov := prepareGoCov(t)
	p, err := cov.Parse("c.out")
	require.NoError(t, err)
	assert.Len(t, p, 1)
	assert.Equal(t, "AAAxHjEeMQ==", p["main.go"])
}
